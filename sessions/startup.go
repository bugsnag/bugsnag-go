package sessions

import (
	"context"
	"net/http"
	"os"
)

const (
	startupSessionIDKey        = "BUGSNAG_STARTUP_SESSION_ID"
	startupSessionTimestampKey = "BUGSNAG_STARTUP_SESSION_TIMESTAMP"
)

// SendStartupSession is called by Bugsnag on startup, which will send a
// session to Bugsnag and return a context to represent the session of the main
// goroutine. This is the session associated with any fatal panics that are
// caught by panicwrap.
func SendStartupSession(config *SessionTrackingConfiguration) context.Context {
	ctx := context.Background()
	session := newSession()
	if !config.IsAutoCaptureSessions() || isApplicationProcess(session) {
		return ctx
	}
	publisher := &publisher{
		config: config,
		client: &http.Client{Transport: config.Transport},
	}
	go publisher.publish([]*Session{session})
	return context.WithValue(ctx, contextSessionKey, session)
}

// Checks to see if this is the application process, as opposed to the process
// that monitors for panics
func isApplicationProcess(session *Session) bool {
	// Application process is run first, and this will only have been set when
	// the monitoring process runs
	envID := os.Getenv(startupSessionIDKey)
	os.Setenv(startupSessionIDKey, session.ID.String())
	os.Setenv(startupSessionTimestampKey, session.StartedAt.String())
	return envID == ""
}
