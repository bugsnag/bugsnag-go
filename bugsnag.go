package bugsnag

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/bugsnag/bugsnag-go/errors"
	"github.com/bugsnag/bugsnag-go/sessions"

	// Fixes a bug with SHA-384 intermediate certs on some platforms.
	// - https://github.com/bugsnag/bugsnag-go/issues/9
	_ "crypto/sha512"
)

// The current version of bugsnag-go.
const VERSION = "1.3.1"
const configuredMultipleTimes = "WARNING: Bugsnag was configured twice. It is recommended to only call bugsnag.Configure once to ensure consistent session tracking behavior"

var once sync.Once
var middleware middlewareStack

// The configuration for the default bugsnag notifier.
var Config Configuration
var sessionTrackingConfig sessions.SessionTrackingConfiguration

// DefaultSessionPublishInterval defines how often sessions should be sent to
// Bugsnag.
// Deprecated: Exposed for developer sanity in testing. Modify at own risk.
var DefaultSessionPublishInterval = 60 * time.Second
var defaultNotifier = Notifier{&Config, nil}
var sessionTracker sessions.SessionTracker

// Configure Bugsnag. The only required setting is the APIKey, which can be
// obtained by clicking on "Settings" in your Bugsnag dashboard. This function
// is also responsible for installing the global panic handler, so it should be
// called as early as possible in your initialization process.
func Configure(config Configuration) {
	Config.update(&config)
	once.Do(Config.PanicHandler)
	startSessionTracking()
}

// StartSession creates a clone of the context.Context instance with Bugsnag
// session data attached.
func StartSession(ctx context.Context) context.Context {
	return sessionTracker.StartSession(ctx)
}

// Notify sends an error to Bugsnag along with the current stack trace. The
// rawData is used to send extra information along with the error. For example
// you can pass the current http.Request to Bugsnag to see information about it
// in the dashboard, or set the severity of the notification.
func Notify(rawData ...interface{}) error {
	return defaultNotifier.Notify(rawData...)
}

// AutoNotify logs a panic on a goroutine and then repanics.
// It should only be used in places that have existing panic handlers further
// up the stack. The rawData is used to send extra information along with any
// panics that are handled this way.
// Usage:
//  go func() {
//		defer bugsnag.AutoNotify()
//      // (possibly crashy code)
//  }()
// See also: bugsnag.Recover()
func AutoNotify(rawData ...interface{}) {
	if err := recover(); err != nil {
		severity := defaultNotifier.getDefaultSeverity(rawData, SeverityError)
		state := HandledState{SeverityReasonHandledPanic, severity, true, ""}
		rawData = append([]interface{}{state}, rawData...)
		defaultNotifier.NotifySync(append(rawData, errors.New(err, 2), true)...)
		panic(err)
	}
}

// Recover logs a panic on a goroutine and then recovers.
// The rawData is used to send extra information along with
// any panics that are handled this way
// Usage: defer bugsnag.Recover()
func Recover(rawData ...interface{}) {
	if err := recover(); err != nil {
		severity := defaultNotifier.getDefaultSeverity(rawData, SeverityWarning)
		state := HandledState{SeverityReasonHandledPanic, severity, false, ""}
		rawData = append([]interface{}{state}, rawData...)
		defaultNotifier.Notify(append(rawData, errors.New(err, 2))...)
	}
}

// OnBeforeNotify adds a callback to be run before a notification is sent to
// Bugsnag.  It can be used to modify the event or its MetaData. Changes made
// to the configuration are local to notifying about this event. To prevent the
// event from being sent to Bugsnag return an error, this error will be
// returned from bugsnag.Notify() and the event will not be sent.
func OnBeforeNotify(callback func(event *Event, config *Configuration) error) {
	middleware.OnBeforeNotify(callback)
}

// Handler creates an http Handler that notifies Bugsnag any panics that
// happen. It then repanics so that the default http Server panic handler can
// handle the panic too. The rawData is used to send extra information along
// with any panics that are handled this way.
func Handler(h http.Handler, rawData ...interface{}) http.Handler {
	notifier := New(rawData...)
	if h == nil {
		h = http.DefaultServeMux
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer notifier.AutoNotify(StartSession(r.Context()), r)
		h.ServeHTTP(w, r)
	})
}

// HandlerFunc creates an http HandlerFunc that notifies Bugsnag about any
// panics that happen. It then repanics so that the default http Server panic
// handler can handle the panic too. The rawData is used to send extra
// information along with any panics that are handled this way. If you have
// already wrapped your http server using bugsnag.Handler() you don't also need
// to wrap each HandlerFunc.
func HandlerFunc(h http.HandlerFunc, rawData ...interface{}) http.HandlerFunc {
	notifier := New(rawData...)

	return func(w http.ResponseWriter, r *http.Request) {
		defer notifier.AutoNotify(r)
		h(w, r)
	}
}

func init() {
	// Set up builtin middlewarez
	OnBeforeNotify(httpRequestMiddleware)

	// Default configuration
	sourceRoot := ""
	if gopath := os.Getenv("GOPATH"); len(gopath) > 0 {
		sourceRoot = filepath.Join(gopath, "src") + "/"
	} else {
		sourceRoot = filepath.Join(runtime.GOROOT(), "src") + "/"
	}
	Config.update(&Configuration{
		APIKey: "",
		Endpoints: Endpoints{
			Notify:   "https://notify.bugsnag.com",
			Sessions: "https://sessions.bugsnag.com",
		},
		Hostname:            "",
		AppType:             "",
		AppVersion:          "",
		AutoCaptureSessions: true,
		ReleaseStage:        "",
		ParamsFilters:       []string{"password", "secret"},
		SourceRoot:          sourceRoot,
		// * for app-engine
		ProjectPackages:     []string{"main*"},
		NotifyReleaseStages: nil,
		Logger:              log.New(os.Stdout, log.Prefix(), log.Flags()),
		PanicHandler:        defaultPanicHandler,
		Transport:           http.DefaultTransport,
	})

	hostname, err := os.Hostname()
	if err == nil {
		Config.Hostname = hostname
	}
}

func startSessionTracking() {
	sessionTrackingConfig.Update(&sessions.SessionTrackingConfiguration{
		APIKey:          Config.APIKey,
		Endpoint:        Config.Endpoints.Sessions,
		Version:         VERSION,
		PublishInterval: DefaultSessionPublishInterval,
		Transport:       Config.Transport,
		ReleaseStage:    Config.ReleaseStage,
		Hostname:        Config.Hostname,
		AppType:         Config.AppType,
		AppVersion:      Config.AppVersion,
		Logger:          Config.Logger,
	})
	if sessionTracker != nil {
		Config.logf(configuredMultipleTimes)
	} else {
		sessionTracker = sessions.NewSessionTracker(&sessionTrackingConfig)
	}
}
