package sessions

import (
	"context"
	"net/http"
	"os"
	"time"
)

const (
	startupStageKey    = "BUGSNAG_STARTUP_STAGE"
	initialSessionSent = "INITIAL_SESSION_SENT"
)

// SendStartupSession is called by Bugsnag on startup, which will send a
// session to Bugsnag and return a context to represent the session of the main
// thread. This is the session associated with any fatal panics that are caught
// by panicwrap.
func SendStartupSession(config *SessionTrackingConfiguration) context.Context {
	ctx := context.Background()
	if alreadySentStartupSession() {
		return ctx
	}
	session := newSession()
	publisher := &publisher{
		config: config,
		client: &http.Client{Transport: config.Transport},
	}
	publisher.publish([]*Session{session})
	// This blocks the application from continuing (and possibly crashing)
	// before we've sent the session, but don't block for too long, i.e.
	// nothing is synchronous.
	// TODO: make this number configurable for test sanity, if nothing else.
	time.Sleep(100 * time.Millisecond)
	return context.WithValue(ctx, contextSessionKey, session)
}

func alreadySentStartupSession() bool {
	stage := os.Getenv(startupStageKey)
	if stage == "" {
		os.Setenv(startupStageKey, initialSessionSent)
	}
	return stage == initialSessionSent
}
