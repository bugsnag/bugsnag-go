package sessions

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type session struct {
	startedAt time.Time
	id        uuid.UUID
}

func newSession() session {
	return session{
		startedAt: time.Now(),
		id:        uuid.NewV4(),
	}
}
