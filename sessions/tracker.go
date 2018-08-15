package sessions

import (
	"context"
	"net/http"
	"time"
)

const (
	//contextSessionKey is a unique key for accessing and setting Bugsnag
	//session data on a context.Context object
	contextSessionKey ctxKey = 1
)

// Customa type alias to ensure uniqueness of context key
type ctxKey int

// SessionTracker exposes a methods for starting sessions that can be used for
// gauging your application's health
type SessionTracker interface {
	StartSession(context.Context) context.Context
}

type sessionTracker struct {
	sessionChannel chan (session)
	sessions       []session
	config         SessionTrackingConfiguration
	publisher      sessionPublisher
}

// NewSessionTracker creates a new SessionTracker based on the provided config,
func NewSessionTracker(config SessionTrackingConfiguration) SessionTracker {
	publisher := defaultPublisher{
		config: config,
		client: &http.Client{Transport: config.Transport},
	}
	return &sessionTracker{
		sessionChannel: make(chan session, 1),
		sessions:       []session{},
		config:         config,
		publisher:      &publisher,
	}
}

func (s *sessionTracker) StartSession(ctx context.Context) context.Context {
	session := newSession()
	s.sessionChannel <- session
	return context.WithValue(ctx, contextSessionKey, session)
}

func (s *sessionTracker) processSessions() {
	tic := time.Tick(s.config.PublishInterval)
	for {
		select {
		case session := <-s.sessionChannel:
			s.sessions = append(s.sessions, session)
		case <-tic:
			oldSessions := s.sessions
			s.sessions = nil
			if len(oldSessions) > 0 {
				s.publisher.publish(oldSessions)
			}
		} //TODO: case for shutdown signal
	}
}
