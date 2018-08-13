package bugsnag

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type sessionTracker struct {
	interval int
	sessions []session
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
