package bugsnag

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	//Unique key for accessing and setting Bugsnag session data on a context.Context object
	contextSessionKey ctxKey = 1
)

var defaultSessionTracker = makeDefaultSessionTracker()

type sessionTracker struct {
	interval       time.Duration
	sessionChannel chan (session)
	sessions       []session
	publisher      sessionPublisher
}

type ctxKey int

func (s *sessionTracker) startSession() *session {
	session := session{
		startedAt: time.Now(),
		id:        uuid.NewV4(),
	}
	s.sessionChannel <- session
	return &session
}

func (s *sessionTracker) processSessions() {
	tic := time.Tick(s.interval)
	for {
		select {
		case session := <-s.sessionChannel:
			s.sessions = append(s.sessions, session)
		case <-tic:
			oldSessions := s.sessions
			s.sessions = nil
			s.publisher.publish(oldSessions)
		} //TODO: case for shutdown signal
	}
}

// StartSession creates a clone of the context.Context instance with Bugsnag
// session data attached.
func StartSession(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextSessionKey, defaultSessionTracker.startSession())
}

func makeDefaultSessionTracker() *sessionTracker {
	p := defaultSessionPublisher{config: Config}
	return &sessionTracker{
		interval:       60 * time.Second,
		sessionChannel: make(chan session, 1),
		sessions:       []session{},
		publisher:      &p,
	}
}
