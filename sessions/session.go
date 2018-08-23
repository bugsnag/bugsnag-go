package sessions

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Session represents a start time and a unique ID that identifies the session.
type Session struct {
	StartedAt time.Time
	ID        uuid.UUID
}

func newSession() *Session {
	sessionID, _ := uuid.NewV4()
	return &Session{
		StartedAt: time.Now(),
		ID:        sessionID,
	}
}
