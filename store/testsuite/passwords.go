package testsuite

import (
	"sync/atomic"
	"testing"

	"golang.org/x/net/context"

	"src.sourcegraph.com/sourcegraph/store"
)

// Passwords_SetPassword_empty tests changing the password to an
// empty password.
func Passwords_SetPassword_empty(ctx context.Context, t *testing.T, s store.Password) {
	uid := nextUID()
	if err := s.SetPassword(ctx, uid, ""); err == nil {
		t.Fatal("err == nil")
	}
}

// Passwords_SetPassword_setToEmpty tests changing the password FROM a
// valid password to an empty password.
func Passwords_SetPassword_setToEmpty(ctx context.Context, t *testing.T, s store.Password) {
	uid := nextUID()
	if err := s.SetPassword(ctx, uid, "p"); err != nil {
		t.Fatal(err)
	}

	// Set to empty
	if err := s.SetPassword(ctx, uid, ""); err == nil {
		t.Fatal("err == nil")
	}

	// Password should remain as "p".
	if err := s.CheckUIDPassword(ctx, uid, "p"); err != nil {
		t.Fatal(err)
	}
	if err := s.CheckUIDPassword(ctx, uid, "p2"); err == nil {
		t.Fatal("err == nil")
	}
}

var testUID int32

// nextUID returns a unique test user UID for this process. This is needed
// since we do sets and compares on passwords for users, and if tests are
// running in parallel the results returned will be racey.
func nextUID() int32 {
	return atomic.AddInt32(&testUID, 1)
}
