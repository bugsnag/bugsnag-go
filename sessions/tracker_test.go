package sessions

import (
	"context"
	"sync"
	"testing"
	"time"
)

type testPublisher struct {
	mutex            sync.Mutex
	sessionsReceived [][]session
}

var pub = testPublisher{
	mutex:            sync.Mutex{},
	sessionsReceived: [][]session{},
}

func (pub *testPublisher) publish(sessions []session) error {
	pub.mutex.Lock()
	defer pub.mutex.Unlock()
	pub.sessionsReceived = append(pub.sessionsReceived, sessions)
	return nil
}

func TestStartSessionModifiesContext(t *testing.T) {
	type ctxKey string
	var k ctxKey
	k, v := "key", "val"
	st, c := makeSessionTracker()
	defer close(c)

	ctx := st.StartSession(context.WithValue(context.Background(), k, v))
	if got, exp := ctx.Value(k), v; got != exp {
		t.Errorf("Changed pre-existing key '%s' with value '%s' into %s", k, v, got)
	}
	if got := ctx.Value(contextSessionKey); got == nil {
		t.Fatalf("No session information applied to context %v", ctx)
	}

	var s session
	got := ctx.Value(contextSessionKey)
	switch got.(type) {
	case session:
		s = got.(session)
	default:
		t.Fatalf("Expected a session to be set on the context but was of wrong type: %T", got)
	}

	verifyValidSession(t, s)
}

func TestShouldOnlyWriteWhenReceivingSessions(t *testing.T) {
	st, c := makeSessionTracker()
	defer close(c)
	go st.processSessions()
	time.Sleep(10 * st.config.PublishInterval) // Would publish many times in this time period if there were sessions

	if got := pub.sessionsReceived; len(got) != 0 {
		t.Errorf("pub was invoked unexpectedly %d times with arguments: %v", len(got), got)
	}

	for i := 0; i < 50000; i++ {
		st.StartSession(context.Background())
	}
	time.Sleep(st.config.PublishInterval * 2)

	var sessions []session
	pub.mutex.Lock()
	defer pub.mutex.Unlock()
	for _, s := range pub.sessionsReceived {
		for _, session := range s {
			verifyValidSession(t, session)
			sessions = append(sessions, session)
		}
	}
	if exp, got := 50000, len(sessions); exp != got {
		t.Errorf("Expected %d sessions but got %d", exp, got)
	}

}

func makeSessionTracker() (*sessionTracker, chan session) {
	c := make(chan session, 1)
	return &sessionTracker{
		config: &SessionTrackingConfiguration{
			PublishInterval: time.Millisecond * 10, //Publish very fast
		},
		sessionChannel: c,
		sessions:       []session{},
		publisher:      &pub,
	}, c
}

func verifyValidSession(t *testing.T, s session) {
	if (s.startedAt == time.Time{}) {
		t.Errorf("Expected start time to be set but was nil")
	}
	if len(s.id) != 16 {
		t.Errorf("Expected UUID to be a valid V4 UUID but was %s", s.id)
	}
}
