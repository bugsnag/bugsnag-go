package sessions

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type session struct {
	startedAt time.Time
	id        uuid.UUID
}

func newSession() *session {
	sessionID, _ := uuid.NewV4()
	return &session{
		startedAt: time.Now(),
		id:        sessionID,
	}
}
