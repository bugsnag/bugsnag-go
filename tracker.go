package bugsnag

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
)

type sessionTracker struct {
	interval int
	sessions []session
}

type ctxKey struct{}

var defaultSessionTracker = newSessionTracker(60)

var contextSessionKey = ctxKey{}

// StartSession creates a clone of the context.Context instance with Bugsnag
// session data attached.
func StartSession(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextSessionKey, defaultSessionTracker.startSession())
}

func (s *sessionTracker) startSession() *session {
	session := session{
		startedAt: time.Now(),
		id:        uuid.NewV4(),
	}
	s.sessions = append(s.sessions, session)
	return &session
}

func newSessionTracker(interval int) *sessionTracker {
	return &sessionTracker{
		interval: interval,
	}
}
