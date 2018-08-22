package sessions

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Session represents a start time and a unique ID that identifies the session.
type Session struct {
	startedAt time.Time
	id        uuid.UUID
}

func newSession() *Session {
	sessionID, _ := uuid.NewV4()
	return &Session{
		startedAt: time.Now(),
		id:        sessionID,
	}
}
