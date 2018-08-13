package bugsnag

import (
	"testing"
	"time"
)

func TestTrackerStartsSession(t *testing.T) {
	tracker := newSessionTracker(60)
	session := tracker.startSession()
	if (session.startedAt == time.Time{}) {
		t.Errorf("Expected start time to be set but was nil")
	}
	if got := session.id; len(got) != 16 {
		t.Errorf("Expected UUID to be a valid V4 UUID but was %s", got)
	}
	if exp, got := 1, len(tracker.sessions); exp != got {
		t.Errorf("Expected '%d' created sessions but was %d", exp, got)
	}
}
