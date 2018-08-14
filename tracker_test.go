package bugsnag

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

func (p *testPublisher) publish(sessions []session) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.sessionsReceived = append(p.sessionsReceived, sessions)
	return nil
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

func TestShouldOnlyWriteWhenReceivingSessions(t *testing.T) {
	c := make(chan session, 1)
	defer close(c)
	p := testPublisher{
		mutex:            sync.Mutex{},
		sessionsReceived: [][]session{},
	}
	st := &sessionTracker{
		interval:       time.Millisecond * 10, //Publish very fast
		sessionChannel: c,
		sessions:       []session{},
		publisher:      &p,
	}

	//Would publish many times in this time period if there were sessions
	time.Sleep(10 * st.interval)

	if got := len(p.sessionsReceived); got != 0 {
		t.Errorf("Publisher was invoked unexpectedly %d times with arguments: %v", got, p.sessionsReceived)
	}

	go st.processSessions()
	for i := 0; i < 5; i++ {
		st.startSession()
		time.Sleep(st.interval)
	}

	p.mutex.Lock()
	if got := len(p.sessionsReceived); got == 0 {
		t.Errorf("Publisher was not invoked")
	}
	p.mutex.Unlock()
	var sessions []session
	for _, s := range p.sessionsReceived {
		for _, session := range s {
			verifyValidSession(t, &session)
			sessions = append(sessions, session)
		}
	}

	if exp, got := 5, len(sessions); exp != got {
		t.Errorf("Expected %d sessions but got %d", exp, got)
	}

}

func verifyValidSession(t *testing.T, s *session) {
	if (s.startedAt == time.Time{}) {
		t.Errorf("Expected start time to be set but was nil")
	}
	if len(s.id) != 16 {
		t.Errorf("Expected UUID to be a valid V4 UUID but was %s", s.id)
	}
}
