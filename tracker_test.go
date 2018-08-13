package bugsnag

import (
	"context"
	"testing"
	"time"
)

func TestTrackerStartsSession(t *testing.T) {
	tracker := newSessionTracker(60)
	session := tracker.startSession()
	verifyValidSession(t, session)
	if exp, got := 1, len(tracker.sessions); exp != got {
		t.Errorf("Expected '%d' created sessions but was %d", exp, got)
	}
}

func TestStartSessionModifiesContext(t *testing.T) {
	type ctxKey string
	var k ctxKey
	k, v := "key", "val"

	ctx := StartSession(context.WithValue(context.Background(), k, v))
	if got, exp := ctx.Value(k), v; got != exp {
		t.Errorf("Changed pre-existing key '%s' with value '%s' into %s", k, v, got)
	}
	if got := ctx.Value(contextSessionKey); got == nil {
		t.Fatalf("No session information applied to context %v", ctx)
	}

	var s *session
	got := ctx.Value(contextSessionKey)
	switch got.(type) {
	case *session:
		s = got.(*session)
	default:
		t.Fatalf("Expected a session to be set on the context but was of wrong type: %T", got)
	}

	verifyValidSession(t, s)
}

func verifyValidSession(t *testing.T, s *session) {
	if (s.startedAt == time.Time{}) {
		t.Errorf("Expected start time to be set but was nil")
	}
	if len(s.id) != 16 {
		t.Errorf("Expected UUID to be a valid V4 UUID but was %s", s.id)
	}
}
