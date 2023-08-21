package redispool

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

var (
	tokenBucketGlobalPrefix                   = "v2:rate_limiters"
	bucketLastReplenishmentTimestampKeySuffix = "last_replenishment_timestamp"
	bucketQuotaConfigKeySuffix                = "config:bucket_quota"
	bucketReplenishmentConfigKeySuffix        = "config:bucket_replenishment_interval_seconds"
	bucketMaxCapacity                         = 10
)

type RateLimiter interface {
	// GetTokensFromBucket gets tokens from the specified rate limit token bucket.
	// bucketName: the name of the bucket where the tokens are, e.g. github.com:api_tokens
	// tokensWanted: the number of tokens you want from the bucket; if the number is greater than the remaining capacity, no tokens will be given.
	GetTokensFromBucket(ctx context.Context, bucketName string, maxTimeToWait int) (allowed bool, remainingTokens int, err error)

	// SetTokenBucketReplenishment sets the configuration for the specified token bucket.
	// bucketName: the name of the bucket where the tokens are, e.g. github.com:api_tokens
	// bucketCapacity: the number of tokens the bucket can hold.
	// bucketReplenishRateSeconds: how often (in seconds) the bucket should be completely replenished.
	SetTokenBucketReplenishment(ctx context.Context, bucketName string, bucketCapacity, bucketReplenishRateSeconds int) error
}

type rateLimiter struct {
	prefix                 string
	pool                   *redis.Pool
	getTokensScript        redis.Script
	setReplenishmentScript redis.Script
}

func NewRateLimiter() (RateLimiter, error) {
	var err error

	pool, ok := Cache.Pool()
	if !ok {
		err = errors.New("unable to get redis connection")
	}

	return &rateLimiter{
		prefix: tokenBucketGlobalPrefix,
		pool:   pool,
		// 3 is the key count, keys are arguments passed to the lua script that will be used to get values from Redis KV.
		getTokensScript:        *redis.NewScript(3, getTokensFromBucketLuaScript),
		setReplenishmentScript: *redis.NewScript(3, setTokenBucketReplenishmentLuaScript),
	}, err
}

func (r *rateLimiter) GetTokenFromBucket(ctx context.Context, bucketName string) (getTokensFromBucketResponse, error) {
	// The test code calls lim.wait with a fake timer generator.
	// This is the real timer generator.
	newTimer := func(d time.Duration) (<-chan time.Time, func() bool, func()) {
		timer := time.NewTimer(d)
		return timer.C, timer.Stop, func() {}
	}
	// Check if ctx is already cancelled
	select {
	case <-ctx.Done():
		return getTokensFromBucketResponse{}, ctx.Err()
	default:
	}
	// Determine wait limit
	waitLimit := time.Duration(math.MaxInt32)
	if deadline, ok := ctx.Deadline(); ok {
		waitLimit = deadline.Sub(t)
	}
	// Reserve
	resp, err := r.getTokenFromBucket(ctx, bucketName, waitLimit)
	if err != nil {
		return getTokensFromBucketResponse{}, err
	}

	// Wait if necessary
	delay := r.DelayFrom(t)
	if delay == 0 {
		return nil
	}
	ch, stop, advance := newTimer(delay)
	defer stop()
	advance() // only has an effect when testing
	select {
	case <-ch:
		// We can proceed.
		return nil
	case <-ctx.Done():
		// Context was canceled before we could proceed.  Cancel the
		// reservation, which may permit other events to proceed sooner.
		r.Cancel()
		return ctx.Err()
	}

	return r.wait(ctx, n, time.Now(), newTimer)
}

func (r *rateLimiter) getTokenFromBucket(ctx context.Context, bucketName string, maxTimeToWait time.Duration) (getTokensFromBucketResponse, error) {
	var response getTokensFromBucketResponse
	keys := r.getRateLimiterKeys(bucketName)
	now := time.Now().Unix()
	result, err := r.getTokensScript.DoContext(ctx, r.pool.Get(), keys.BucketKey, keys.LastReplenishmentTimestampKey, keys.QuotaKey, keys.ReplenishmentIntervalSecondsKey, bucketMaxCapacity, now, maxTimeToWait)
	if err != nil {
		return response, errors.Wrapf(err, "error while getting tokens from bucket %s", keys.BucketKey)
	}

	scriptResponse, ok := result.([]interface{})
	if !ok || len(scriptResponse) != 3 {
		return response, errors.Newf("unexpected response from Redis when getting tokens from bucket: %s, response: %+v", keys.BucketKey, response)
	}

	allowedInt, ok := scriptResponse[0].(int64)
	if !ok {
		return response, errors.Newf("unexpected response for allowed, expected int64 but got %T", allowedInt)
	}

	timeToWaitSeconds, ok := scriptResponse[1].(int64)
	if !ok {
		return response, errors.Newf("unexpected response for timeToWait, expected int64, got %T", timeToWaitSeconds)
	}

	remTokens, ok := scriptResponse[2].(int64)
	if !ok {
		return response, errors.Newf("unexpected response for tokens left, expected int64, got %T", remTokens)
	}

	response = getTokensFromBucketResponse{
		Allowed:           rateLimitScriptGrantResponse(allowedInt),
		TimeToWaitSeconds: time.Duration(timeToWaitSeconds) * time.Second,
		TokensRemaining:   int(remTokens),
	}

	return response, response.Error()
}

func (r *rateLimiter) SetTokenBucketReplenishment(ctx context.Context, bucketName string, bucketQuota, bucketReplenishRateSeconds int) error {
	keys := r.getRateLimiterKeys(bucketName)
	_, err := r.setReplenishmentScript.DoContext(ctx, r.pool.Get(), keys.BucketKey, keys.QuotaKey, keys.ReplenishmentIntervalSecondsKey, bucketQuota, bucketReplenishRateSeconds)
	return errors.Wrapf(err, "error while setting token bucket replenishment for bucket %s", bucketName)
}

func (r *rateLimiter) getRateLimiterKeys(bucketName string) rateLimitBucketConfigKeys {
	var keys rateLimitBucketConfigKeys
	// e.g. v2:rate_limiters:github.com:api_tokens
	keys.BucketKey = fmt.Sprintf("%s:%s", r.prefix, bucketName)
	// e.g. v2:rate_limiters:github.com:api_tokens:config:bucket_quota
	keys.QuotaKey = fmt.Sprintf("%s:%s", keys.BucketKey, bucketQuotaConfigKeySuffix)
	// e.g.. v2:rate_limiters:github.com:api_tokens:config:bucket_replenishment_interval_seconds
	keys.ReplenishmentIntervalSecondsKey = fmt.Sprintf("%s:%s", keys.BucketKey, bucketReplenishmentConfigKeySuffix)
	// e.g.. v2:rate_limiters:github.com:api_tokens:last_replenishment_timestamp
	keys.LastReplenishmentTimestampKey = fmt.Sprintf("%s:%s", keys.BucketKey, bucketLastReplenishmentTimestampKeySuffix)

	return keys
}

// getTokensFromBucketLuaScript gets a single token from the specified bucket.
// bucket_key: the key in Redis that stores the bucket, under which all the bucket's tokens and rate limit configs are found, e.g. v2:rate_limiters:github.com:api_tokens.
// last_replenishment_timestamp_key: the key in Redis that stores the timestamp (seconds since epoch) of the last bucket replenishment, e.g. v2:rate_limiters:github.com:api_tokens:last_replenishment_timestamp.
// bucket_quota_key: the key in Redis that stores how many tokens the bucket should refill in a `bucket_replenishment_interval` period of time, e.g. v2:rate_limiters:github.com:api_tokens:config:bucket_quota.
// bucket_replenishment_interval_key: the key in Redis that stores how often (in seconds), the bucket should be replenished bucket_quota tokens, e.g. v2:rate_limiters:github.com:api_tokens:config:bucket_replenishment_interval_seconds.
// bucket_capacity: the amount of tokens the bucket can hold, always bucketMaxCapacity right now.
// current_time: current time (seconds since epoch).
// max_time_to_wait_for_token: the maximum amount of time (in seconds) the requester is willing to wait before acquiring/using a token.
const getTokensFromBucketLuaScript = `local bucket_key = KEYS[1]
local last_replenishment_timestamp_key = KEYS[2]
local bucket_quota_key = KEYS[3]
local bucket_replenishment_interval_key = KEYS[4]
local bucket_capacity tonumber(ARGV[1])
local current_time = tonumber(ARGV[2])
local max_time_to_wait_for_token = tonumber(ARGV[3])

-- Check if the bucket exists.
local bucket_exists = redis.call('EXISTS', bucket_key)

-- If the bucket does not exist, create the bucket, and set the last replenishment time.
if bucket_exists == 0 then
    redis.call('SET', bucket_key, capacity)
    redis.call('SET', last_replenishment_timestamp_key, current_time)
end

-- Calculate the time difference in seconds since last replenishment
local last_replenishment_timestamp = tonumber(redis.call('GET', last_replenishment_timestamp_key))
local time_difference = current_time - last_replenishment_timestamp
-- Shouldn't happen, but check just in case.
if time_difference < 0 then
	return {-1, 0, 0}
end

-- Get the rate (tokens/second) that the bucket should replenish
local bucket_quota = tonumber(redis.call('GET', bucket_quota_key))
local bucket_replenishment_interval =  tonumber(redis.call('GET', bucket_replenishment_interval_key))
local replenishment_rate = bucket_quota / bucket_replenishment_interval

-- Calculate the amount of tokens to replenish, round down for the number of 'full' tokens.
local num_tokens_to_replenish = math.floor(replenishment_rate * time_difference)

-- Get the current token count in the bucket.
local current_tokens = tonumber(redis.call('GET', bucket_key))

-- Replenish the bucket if there are tokens to replenish
if num_tokens_to_replenish > 0 then
    local available_capacity = capacity - current_tokens
    if available_capacity > 0 then
		-- The number of tokens we add is either the number of tokens we have replenished over
		-- the last time_difference, or enough tokens to refill the bucket completely, whichever
		-- is lower.
		current_tokens = math.min(capacity, current_tokens + num_tokens_to_replenish)
    	redis.call('SET', bucket_key, math.min(capacity, current_tokens))
    	redis.call('SET', last_replenishment_timestamp_key, current_time)
    end
end

local time_to_wait_for_token = 0
-- This is for calculations with us removing a token.
local current_tokens_after_consumption = current_tokens - 1

-- If the bucket will be negative, calculate the needed to 'wait' before using the token.
-- i.e. if we are going to be at -15 tokens after this consumption, and the token replenishment
-- rate is 0.33/s, then we need to wait 45.45 (46 because we round up) seconds before making the request.
if current_tokens_after_consumption < 0 then
    time_to_wait_for_token = math.ceil((current_tokens_after_consumption * -1) / replenishment_rate)
end

if time_to_wait_for_token >= max_time_to_wait_for_token then
    return {0, time_to_wait_for_token, current_tokens}
end

-- Decrement the token bucket by 1, we are granted a token
redis.call('DECRBY', bucket_key, 1)

return {1, time_to_wait_for_token, current_tokens_after_consumption}`

const setTokenBucketReplenishmentLuaScript = `local bucket_key = KEYS[1]
local bucket_capacity_key = KEYS[2]
local replenish_interval_seconds_key = KEYS[3]
local bucket_capacity = tonumber(ARGV[1])
local bucket_replenish_rate_seconds = tonumber(ARGV[2])

-- Get the current number of tokens in the bucket.
local current_tokens = tonumber(redis.call('GET', bucket_key) or 0)

-- Set bucket capacity if the current # of tokens is > the new capacity.
if current_tokens > bucket_capacity then
	redis.call('SET', bucket_key, bucket_capacity)
end

redis.call('SET', bucket_capacity_key, bucket_capacity)
redis.call('SET', replenish_interval_seconds_key, bucket_replenish_rate_seconds)`

type getTokensFromBucketResponse struct {
	Allowed           rateLimitScriptGrantResponse
	TimeToWaitSeconds time.Duration
	TokensRemaining   int
}

type rateLimitScriptGrantResponse int64

var (
	tokenGranted                 rateLimitScriptGrantResponse = 1
	tokenGrantExceedsMaxWaitTime rateLimitScriptGrantResponse = 0
	negativeTimeDifference       rateLimitScriptGrantResponse = -1
)

func (r getTokensFromBucketResponse) Error() error {
	switch r.Allowed {
	case tokenGranted:
		return nil
	case tokenGrantExceedsMaxWaitTime:
		return errors.Errorf("acquiring token would require a wait of %+v seconds which exceeds the limit", r.TimeToWaitSeconds)
	case negativeTimeDifference:
		return errors.New("time difference between now and the last replenishment is negative")
	default:
		return errors.New("unexpected return code from rate limit script")
	}
}

type rateLimitBucketConfigKeys struct {
	BucketKey                       string
	QuotaKey                        string
	ReplenishmentIntervalSecondsKey string
	LastReplenishmentTimestampKey   string
}

type RateLimiterConfigNotCreatedError struct {
	tokenBucketKey string
}

func (r *RateLimiterConfigNotCreatedError) Error() string {
	return fmt.Sprintf("config for rate limiter not found: %s", r.tokenBucketKey)
}
