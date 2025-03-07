package completions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/sourcegraph/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"

	"github.com/sourcegraph/sourcegraph/cmd/cody-gateway/internal/actor"
	"github.com/sourcegraph/sourcegraph/cmd/cody-gateway/internal/events"
	"github.com/sourcegraph/sourcegraph/cmd/cody-gateway/internal/httpapi/featurelimiter"
	"github.com/sourcegraph/sourcegraph/cmd/cody-gateway/internal/httpapi/overhead"
	"github.com/sourcegraph/sourcegraph/cmd/cody-gateway/internal/limiter"
	"github.com/sourcegraph/sourcegraph/cmd/cody-gateway/internal/notify"
	"github.com/sourcegraph/sourcegraph/cmd/cody-gateway/internal/response"
	"github.com/sourcegraph/sourcegraph/internal/codygateway"
	"github.com/sourcegraph/sourcegraph/internal/completions/client/fireworks"
	"github.com/sourcegraph/sourcegraph/internal/httpcli"
	"github.com/sourcegraph/sourcegraph/internal/requestclient"
	sgtrace "github.com/sourcegraph/sourcegraph/internal/trace"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

type usageStats struct {
	// characters is the number of characters in the input or response.
	characters int
	// tokens is the number of tokens consumed in the input or response.
	tokens int
}

// Hop-by-Hop headers that should not be copied when proxying upstream requests
// List from https://cs.opensource.google/go/go/+/master:src/net/http/httputil/reverseproxy.go;l=294;drc=7abeefd2b1a03932891e581f1f90656ffebebce4
var hopHeaders = map[string]struct{}{
	"Connection":          {},
	"Proxy-Connection":    {}, // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive":          {},
	"Proxy-Authenticate":  {},
	"Proxy-Authorization": {},
	"Te":                  {}, // canonicalized version of "TE"
	"Trailer":             {}, // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding":   {},
	"Upgrade":             {},
}

// upstreamHandlerMethods declares a set of methods that are used throughout the
// lifecycle of a request to an upstream API. All methods are required, and called
// in the order they are defined here.
//
// Methods do not need to be concurrency-safe, as they are only called sequentially.
type upstreamHandlerMethods[ReqT UpstreamRequest] interface {
	// getAPIURLByFeature returns the upstream API endpoint to call for the given feature.
	getAPIURLByFeature(codygateway.Feature) string

	// validateRequest can be used to validate the HTTP request before it is sent upstream.
	// This is where we enforce things like character/token limits, etc. Any non-nil errors
	// will block the processing of the request, and serve the error directly to the end
	// user along with an HTTP status code 400 Bad Request.
	//
	// The provided logger already contains actor context.
	validateRequest(context.Context, log.Logger, codygateway.Feature, ReqT) error

	// shouldFlagRequest is called after the request has been validated, and is where we
	// run various heuristics to check if the request is abusive in nature. (e.g. suspiciously
	// long, contains words/phrases from a blocklist, etc.)
	//
	// All implementations of this function should call isFlaggedRequest(...), along with
	// any LLM or provider-specific logic.
	//
	// Any errors returned from shouldFlagRequest will be swallowed, and the request will be
	// considered unflagged. (So implementations should return errors rather than swallowing
	// them directly.)
	shouldFlagRequest(context.Context, log.Logger, ReqT) (*flaggingResult, error)
	// transformBody can be used to modify the request body before it is sent
	// upstream. To manipulate the HTTP request, use transformRequest.
	//
	// If the upstream supports it, the given identifier string should be
	// provided to assist in abuse detection.
	transformBody(_ *ReqT, identifier string)
	// transformRequest can be used to modify the HTTP request before it is sent
	// upstream. To manipulate the body, use transformBody.
	transformRequest(*http.Request)
	// getRequestMetadata should extract details about the request we are sending
	// upstream for validation and tracking purposes. Usage data does not need
	// to be reported here - instead, use parseResponseAndUsage to extract usage,
	// which for some providers we can only know after the fact based on what
	// upstream tells us.
	getRequestMetadata(ReqT) (model string, additionalMetadata map[string]any)
	// parseResponseAndUsage should extract details from the response we get back from
	// upstream as well as overall usage for tracking purposes.
	//
	// If data is unavailable, implementations should set relevant usage fields
	// to -1 as a sentinel value.
	parseResponseAndUsage(log.Logger, ReqT, io.Reader) (promptUsage, completionUsage usageStats)
}

type UpstreamRequest interface {
	GetModel() string
	ShouldStream() bool
	// BuildPrompt returns the aggregated prompt (either full prompt as generated by Client, or all messages concatenated)
	BuildPrompt() string
}

// makeUpstreamHandler a big deal. This method will produce an http.Handler that will handle converting
// the Cody Gateway user's request to the backing LLM ("upstream provider"). This is how we provide a
// consistent way for providing logging, telemetry, rate limiting, etc. across multiple upstream providers.
func makeUpstreamHandler[ReqT UpstreamRequest](
	baseLogger log.Logger,
	eventLogger events.Logger,
	rs limiter.RedisStore,
	rateLimitNotifier notify.RateLimitNotifier,
	httpClient httpcli.Doer,

	// upstreamName is the name of the upstream provider. It MUST match the
	// provider names defined clientside, i.e. "anthropic" or "openai".
	upstreamName string,

	allowedModels []string,

	methods upstreamHandlerMethods[ReqT],
	flaggedPromptRecorder PromptRecorder,

	// defaultRetryAfterSeconds sets the retry-after policy on upstream rate
	// limit events in case a retry-after is not provided by the upstream
	// response.
	defaultRetryAfterSeconds int,
	autoFlushStreamingResponses bool,
) http.Handler {
	baseLogger = baseLogger.Scoped(upstreamName)

	// Convert allowedModels to the Cody Gateway configuration format with the
	// provider as a prefix. This aligns with the models returned when we query
	// for rate limits from actor sources.
	clonedAllowedModels := make([]string, len(allowedModels))
	copy(clonedAllowedModels, allowedModels)
	for i := range clonedAllowedModels {
		clonedAllowedModels[i] = fmt.Sprintf("%s/%s", upstreamName, clonedAllowedModels[i])
	}

	// upstreamHandler is the actual HTTP handle that will perform "all of the things"
	// in order to call the upstream API. e.g. calling the upstreamHandlerMethods in
	// the correct order, enforcing rate limits and anti-abuse mechanisms, etc.
	upstreamHandler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		act := actor.FromContext(ctx)

		// TODO: Investigate using actor propagation handler for extracting
		// this. We had some issues before getting that to work, so for now
		// just stick with what we've seen working so far.
		sgActorID := r.Header.Get("X-Sourcegraph-Actor-UID")
		sgActorAnonymousUID := r.Header.Get("X-Sourcegraph-Actor-Anonymous-UID")

		// Build logger for lifecycle of this request with lots of details.
		logger := act.Logger(sgtrace.Logger(ctx, baseLogger)).With(
			append(
				requestclient.FromContext(ctx).LogFields(),
				// Sourcegraph actor details
				log.String("sg.actorID", sgActorID),
				log.String("sg.anonymousID", sgActorAnonymousUID),
			)...,
		)

		feature := featurelimiter.GetFeature(ctx)
		if feature == "" {
			response.JSONError(logger, w, http.StatusBadRequest, errors.New("no feature provided"))
			return
		}

		// This will never be nil as the rate limiter middleware checks this before.
		// TODO: Should we read the rate limit from context, and store it in the rate
		// limiter to make this less dependent on these two logics to remain the same?
		rateLimit, ok := act.RateLimits[feature]
		if !ok {
			response.JSONError(logger, w, http.StatusInternalServerError, errors.Newf("rate limit for %q not found", string(feature)))
			return
		}

		// TEMPORARY: Add provider prefixes to AllowedModels for back-compat
		// if it doesn't look like there is a prefix yet.
		//
		// This isn't very robust, but should tide us through a brief transition
		// period until everything deploys and our caches refresh.
		for i := range rateLimit.AllowedModels {
			if !strings.Contains(rateLimit.AllowedModels[i], "/") {
				rateLimit.AllowedModels[i] = fmt.Sprintf("%s/%s", upstreamName, rateLimit.AllowedModels[i])
			}
		}

		// Parse the request body.
		var body ReqT
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			response.JSONError(logger, w, http.StatusBadRequest, errors.Wrap(err, "failed to parse request body"))
			return
		}
		// Validate the request. (e.g. hard-caps on maximum token size, a known model, etc.)
		if err := methods.validateRequest(ctx, logger, feature, body); err != nil {
			response.JSONError(logger, w, http.StatusBadRequest, err)
			return
		}

		// Check the request to see if it should be flagged for abuse, or additional inspection.
		flaggingResult, err := methods.shouldFlagRequest(ctx, logger, body)
		if err != nil {
			logger.Error("error checking if request should be flagged, treating as non-flagged", log.Error(err))
		}
		if flaggingResult != nil && flaggingResult.IsFlagged() {
			// Record flagged prompts to aid in combating ongoing abuse waves.
			if actor.FromContext(ctx).IsDotComActor() {
				prompt := body.BuildPrompt()
				// We don't record code completions until we get the false-positive count
				// under control. (It's just noise.)
				if feature != codygateway.FeatureCodeCompletions {
					if err := flaggedPromptRecorder.Record(ctx, prompt); err != nil {
						logger.Warn("failed to record flagged prompt", log.Error(err))
					}
				}
			}

			// Requests that are flagged but not outright blocked, will have some of the
			// metadata from flaggingResult attached to the request event telemetry. That's
			// how the data flows into other backend systems for downstream analysis.
			if !flaggingResult.shouldBlock {
				logger.Info("request was flagged, but not blocked. Proceeding.", log.Strings("reasons", flaggingResult.reasons))
			} else {
				requestMetadata := getFlaggingMetadata(flaggingResult, act)
				err := eventLogger.LogEvent(
					ctx,
					events.Event{
						Name:       codygateway.EventNameRequestBlocked,
						Source:     act.Source.Name(),
						Identifier: act.ID,
						Metadata: events.MergeMaps(requestMetadata, map[string]any{
							codygateway.CompletionsEventFeatureMetadataField: feature,
							"model":    fmt.Sprintf("%s/%s", upstreamName, body.GetModel()),
							"provider": upstreamName,

							// Request metadata
							"prompt_token_count":   flaggingResult.promptTokenCount,
							"max_tokens_to_sample": flaggingResult.maxTokensToSample,

							// Actor details, specific to the actor Source
							"sg_actor_id":            sgActorID,
							"sg_actor_anonymous_uid": sgActorAnonymousUID,
						}),
					},
				)
				if err != nil {
					logger.Error("failed to log event", log.Error(err))
				}
				response.JSONError(logger, w, http.StatusBadRequest, requestBlockedError(ctx))
				return
			}
		}

		// identifier that can be provided to upstream for abuse detection
		// has the format '$ACTOR_ID:$SG_ACTOR_ID'. The latter is anonymized
		// (specific per-instance)
		identifier := fmt.Sprintf("%s:%s", act.ID, sgActorID)
		methods.transformBody(&body, identifier)

		// Re-marshal the payload for upstream to unset metadata and remove any properties
		// not known to us.
		upstreamPayload, err := json.Marshal(body)
		if err != nil {
			response.JSONError(logger, w, http.StatusInternalServerError, errors.Wrap(err, "failed to marshal request body"))
			return
		}

		// Create a new request to send upstream, making sure we retain the same context.
		upstreamURL := methods.getAPIURLByFeature(feature)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, upstreamURL, bytes.NewReader(upstreamPayload))
		if err != nil {
			response.JSONError(logger, w, http.StatusInternalServerError, errors.Wrap(err, "failed to create request"))
			return
		}

		// Run the request transformer.
		methods.transformRequest(req)

		// Retrieve metadata from the initial request.
		model, requestMetadata := methods.getRequestMetadata(body)

		// Match the model against the allowlist of models, which are configured
		// with the Cody Gateway model format "$PROVIDER/$MODEL_NAME". Models
		// are sent as if they were against the upstream API, so they don't have
		// the prefix yet when extracted - we need to add it back here. This
		// full gatewayModel is also used in events tracking.
		gatewayModel := fmt.Sprintf("%s/%s", upstreamName, model)
		if allowed := intersection(clonedAllowedModels, rateLimit.AllowedModels); !isAllowedModel(allowed, gatewayModel) {
			response.JSONError(logger, w, http.StatusBadRequest,
				errors.Newf("model %q is not allowed, allowed: [%s]",
					gatewayModel, strings.Join(allowed, ", ")))
			return
		}

		w.Header().Add("x-cody-resolved-model", gatewayModel)

		var (
			upstreamStarted    = time.Now()
			upstreamLatency    time.Duration
			upstreamStatusCode int = -1
			// resolvedStatusCode is the status code that we returned to the
			// client - in most case it is the same as upstreamStatusCode,
			// but sometimes we write something different.
			resolvedStatusCode int = -1
			// promptUsage and completionUsage are extracted from parseResponseAndUsage.
			promptUsage, completionUsage usageStats
		)
		defer func() {
			if span := oteltrace.SpanFromContext(ctx); span.IsRecording() {
				span.SetAttributes(
					attribute.Int("upstreamStatusCode", upstreamStatusCode),
					attribute.Int("resolvedStatusCode", resolvedStatusCode))
			}
			if flaggingResult.IsFlagged() {
				requestMetadata = events.MergeMaps(requestMetadata, getFlaggingMetadata(flaggingResult, act))
			}
			usageData := map[string]any{
				"prompt_character_count":     promptUsage.characters,
				"prompt_token_count":         promptUsage.tokens,
				"completion_character_count": completionUsage.characters,
				"completion_token_count":     completionUsage.tokens,
			}
			for k, v := range usageData {
				// Drop usage fields that are invalid/unimplemented. All
				// usageData fields are ints - we use map[string]any for
				// convenience with mergeMaps utility.
				if n, _ := v.(int); n < 0 {
					delete(usageData, k)
				}
			}
			o := overhead.FromContext(ctx)
			o.Feature = feature
			o.UpstreamLatency = upstreamLatency
			o.Provider = upstreamName
			o.Stream = body.ShouldStream()

			err := eventLogger.LogEvent(
				ctx,
				events.Event{
					Name:       codygateway.EventNameCompletionsFinished,
					Source:     act.Source.Name(),
					Identifier: act.ID,
					Metadata: events.MergeMaps(requestMetadata, usageData, map[string]any{
						codygateway.CompletionsEventFeatureMetadataField: feature,
						"model":    gatewayModel,
						"provider": upstreamName,

						// Request details
						"upstream_request_duration_ms": upstreamLatency.Milliseconds(),
						"upstream_status_code":         upstreamStatusCode,
						"resolved_status_code":         resolvedStatusCode,

						// Actor details, specific to the actor Source
						"sg_actor_id":            sgActorID,
						"sg_actor_anonymous_uid": sgActorAnonymousUID,
					}),
				},
			)
			if err != nil {
				logger.Error("failed to log event", log.Error(err))
			}
		}()
		resp, err := httpClient.Do(req)
		if err != nil {
			// Ignore reporting errors where client disconnected
			if req.Context().Err() == context.Canceled && errors.Is(err, context.Canceled) {
				oteltrace.SpanFromContext(req.Context()).
					SetStatus(codes.Error, err.Error())
				logger.Info("request canceled", log.Error(err))
				return
			}

			// More user-friendly message for timeouts
			if errors.Is(err, context.DeadlineExceeded) {
				resolvedStatusCode = http.StatusGatewayTimeout
				response.JSONError(logger, w, resolvedStatusCode,
					errors.Newf("request to upstream provider %s timed out", upstreamName))
				return
			}

			resolvedStatusCode = http.StatusInternalServerError
			response.JSONError(logger, w, resolvedStatusCode,
				errors.Wrapf(err, "failed to make request to upstream provider %s", upstreamName))
			return
		}
		defer func() { _ = resp.Body.Close() }()
		// Forward upstream http headers.
		for k, vv := range resp.Header {
			if _, ok := hopHeaders[http.CanonicalHeaderKey(k)]; ok {
				// do not forward
				continue
			}
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}

		// Record upstream's status code and decide what we want to send to
		// the client. By default, we just send upstream's status code.
		upstreamStatusCode = resp.StatusCode
		resolvedStatusCode = upstreamStatusCode
		if upstreamStatusCode == http.StatusTooManyRequests {
			// Rewrite 429 to 503 because we share a quota when talking to upstream,
			// and a 429 from upstream should NOT indicate to the client that they
			// should liberally retry until the rate limit is lifted. To ensure we are
			// notified when this happens, log this as an error and record the headers
			// that are provided to us.
			var headers bytes.Buffer
			_ = resp.Header.Write(&headers)
			logger.Error("upstream returned 429, rewriting to 503",
				log.Error(errors.New(resp.Status)), // real error needed for Sentry reporting
				log.String("resp.headers", headers.String()))
			resolvedStatusCode = http.StatusServiceUnavailable
		}

		// This handles upstream 429 responses as well, since they get
		// resolved to http.StatusServiceUnavailable.
		if resolvedStatusCode == http.StatusServiceUnavailable {
			// Propagate retry-after in case it is handle-able by the client,
			// or write our default. 503 errors can have retry-after as well.
			if upstreamRetryAfter := resp.Header.Get("retry-after"); upstreamRetryAfter != "" {
				w.Header().Set("retry-after", upstreamRetryAfter)
			} else {
				w.Header().Set("retry-after", strconv.Itoa(defaultRetryAfterSeconds))
			}
		}

		// Write the resolved status code.
		w.WriteHeader(resolvedStatusCode)

		// Set up a buffer to capture the response as it's streamed and sent to the client.
		var responseBuf bytes.Buffer
		respBody := io.TeeReader(resp.Body, &responseBuf)
		// if this is a streaming request, we want to flush ourselves instead of leaving that to the http.Server
		// (so events are sent to the client as soon as possible)
		var responseWriter io.Writer = w
		if autoFlushStreamingResponses && body.ShouldStream() {
			if fw, err := response.NewAutoFlushingWriter(w); err == nil {
				responseWriter = fw
			} else {
				// We can't stream the response, but it's better to write it without streaming that fail, so we just log the error
				logger.Error("failed to create auto-flushing writer", log.Error(err))
			}
		}
		_, _ = io.Copy(responseWriter, respBody)
		// record latency of upstream request after we read the whole response, but without recording the time we spend parsing it in parseResponseAndUsage()
		upstreamLatency = time.Since(upstreamStarted)

		if upstreamStatusCode >= 200 && upstreamStatusCode < 300 {
			// Pass reader to response transformer to capture token counts.
			promptUsage, completionUsage = methods.parseResponseAndUsage(logger, body, &responseBuf)
		} else if upstreamStatusCode >= 500 {
			logger.Error("error from upstream",
				log.Int("status_code", upstreamStatusCode))
		}
	}

	return featurelimiter.Handle(
		baseLogger,
		eventLogger,
		rs,
		rateLimitNotifier,
		http.HandlerFunc(upstreamHandler))
}

// Trim detected phrases to this many characters (to avoid storing too much repetitive data in BigQuery)
const phrasePrefixLength = 5

func truncateToPrefix(p string) string {
	pat := p
	if len(p) > phrasePrefixLength {
		pat = p[:phrasePrefixLength]
	}
	return pat
}

func getFlaggingMetadata(flaggingResult *flaggingResult, act *actor.Actor) map[string]any {
	requestMetadata := map[string]any{}

	requestMetadata["flagged"] = true
	flaggingMetadata := map[string]any{
		"reason":       flaggingResult.reasons,
		"should_block": flaggingResult.shouldBlock,
	}
	if flaggingResult.blockedPhrase != nil {
		flaggingMetadata["blocked_phrase"] = truncateToPrefix(*flaggingResult.blockedPhrase)
	}

	if act.IsDotComActor() {
		flaggingMetadata["prompt_prefix"] = flaggingResult.promptPrefix
	}
	requestMetadata["flagging_result"] = flaggingMetadata
	return requestMetadata
}

func isAllowedModel(allowedModels []string, model string) bool {
	for _, m := range allowedModels {
		if strings.EqualFold(m, model) {
			return true
		}

		// Expand virtual model names
		if m == "fireworks/starcoder" && (model == "fireworks/"+fireworks.Starcoder7b ||
			model == "fireworks/"+fireworks.Starcoder16b ||
			model == "fireworks/"+fireworks.Starcoder7b8bit ||
			model == "fireworks/"+fireworks.Starcoder16b8bit ||
			model == "fireworks/"+fireworks.Starcoder16bSingleTenant) {
			return true
		}
	}
	return false
}

func intersection(a, b []string) (c []string) {
	for _, val := range a {
		if slices.Contains(b, val) {
			c = append(c, val)
		}
	}
	return c
}
