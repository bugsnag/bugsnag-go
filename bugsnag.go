package bugsnag

import (
	"github.com/bugsnag/bugsnag-go/errors"
	"github.com/mitchellh/panicwrap"
	"log"
	"net/http"
	"os"
	"sync"
)

var once sync.Once
var middleware middlewareStack

// The configuration for the default bugsnag notifier.
var Config Configuration

var defaultNotifier = Notifier{&Config, nil}

// Configure Bugsnag. The most important setting is the APIKey.
// This must be called before any other function on Bugsnag, and
// should be called as early as possible in your program.
func Configure(config Configuration) {
	defaultNotifier.Config.update(&config)
	once.Do(func() {
		if !config.DisablePanicHandler {
			handleUncaughtPanics()
		}
	})
}

// Notify sends an error to Bugsnag. The rawData can be anything supported by Bugsnag,
// e.g. User, Context, SeverityError, MetaData, Configuration,
// or anything supported by your custom middleware. Unsupported values will be silently ignored.
func Notify(err error, rawData ...interface{}) {
	defaultNotifier.Notify(err, rawData...)
}

// defer AutoNotify notifies Bugsnag about any panic()s. It then re-panics() so that existing
// error handling continues to work. The rawData is used to add information to the notification,
// see Notify for more information.
func AutoNotify(rawData ...interface{}) {
	if err := recover(); err != nil {
		rawData = defaultNotifier.addDefaultSeverity(rawData, SeverityError)
		defaultNotifier.Notify(errors.New(err, 2), rawData...)
		panic(err)
	}
}

// defer Recover notifies Bugsnag about any panic()s, and stops panicking so that your program doesn't
// crash. The rawData is used to add information to the notification, see Notify for more information.
func Recover(rawData ...interface{}) {
	if err := recover(); err != nil {
		rawData = defaultNotifier.addDefaultSeverity(rawData, SeverityWarning)
		defaultNotifier.Notify(errors.New(err, 2), rawData...)
	}
}

// OnBeforeNotify adds a callback to be run before a notification is sent to Bugsnag.
// It can be used to modify the event or the config to be used, or to completely cancel
// the notification by returning false. You should return true to continue processing.
func OnBeforeNotify(callback func(event *Event, config *Configuration) bool) {
	middleware.BeforeNotify(callback)
}

// OnAroundNotify adds a callback to be run before a notification is sent to Bugsnag.
// It can be used to modify the event or the config to be used, or to move the request
// to a different goroutine, etc.
// It should call next() to actually send the notification, or avoid calling next() to cancel it.
// Consider using OnBeforeNotify instead for simple cases.
func OnAroundNotify(callback func(event *Event, config *Configuration, next func())) {
	middleware.AddMiddleware(callback)
}

// Handler wraps the HTTP handler in bugsnag.AutoNotify(). It includes details about the
// HTTP request in all error reports. If you don't pass a handler, the default http handlers
// will be used.
func Handler(h http.Handler, rawData ...interface{}) http.Handler {
	notifier := NewNotifier(rawData...)
	if h == nil {
		h = http.DefaultServeMux
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer notifier.AutoNotify(r)
		h.ServeHTTP(w, r)
	})
}

func init() {
	// Set up builtin middlewarez
	OnBeforeNotify(httpRequestMiddleware)

	// Default configuration
	Config.update(&Configuration{
		APIKey:              "",
		Endpoint:            "https://notify.bugsnag.com/",
		Hostname:            "",
		AppVersion:          "",
		ReleaseStage:        "",
		ParamsFilters:       []string{"password", "secret"},
		ProjectPackages:     []string{"main"},
		NotifyReleaseStages: nil,
		Logger:              log.New(os.Stdout, log.Prefix(), log.Flags()),
	})

	hostname, err := os.Hostname()
	if err == nil {
		Config.Hostname = hostname
	}
}

// NOTE: this function does not return when you call it, instead it
// re-exec()s the current process with panic monitoring.
func handleUncaughtPanics() {
	defer defaultNotifier.dontPanic()

	exitStatus, err := panicwrap.Wrap(&panicwrap.WrapConfig{
		CookieKey:   "bugsnag_wrapped",
		CookieValue: "bugsnag_wrapped",
		Handler: func(output string) {

			toNotify, err := errors.ParsePanic(output)

			if err != nil {
				defaultNotifier.Config.log("bugsnag.handleUncaughtPanic: %v", err)
			}
			Notify(toNotify, SeverityError)
		},
	})

	if err != nil {
		defaultNotifier.Config.log("bugsnag.handleUncaughtPanic: %v", err)
		return
	}

	if exitStatus >= 0 {
		os.Exit(exitStatus)
	} else {
		return
	}
}
