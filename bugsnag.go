package bugsnag

import (
	"github.com/bugsnag/bugsnag-go/errors"
	"log"
	"net/http"
	"os"
	"sync"
)

// The current version of the notifier
const VERSION = "0.1"

var once sync.Once
var middleware middlewareStack

// The configuration for the default bugsnag notifier.
var Config Configuration

var defaultNotifier = Notifier{&Config, nil}

// Configure Bugsnag. The most important setting is the APIKey.
// This must be called before any other function on Bugsnag, and
// should be called as early as possible in your program.
func Configure(config Configuration) {
	Config.update(&config)
	once.Do(Config.PanicHandler)
}

// Notify sends an error to Bugsnag. The rawData can be anything supported by Bugsnag,
// e.g. User, Context, SeverityError, MetaData, Configuration,
// or anything supported by your custom middleware. Unsupported values will be silently ignored.
func Notify(err error, rawData ...interface{}) error {
	return defaultNotifier.Notify(err, rawData...)
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
// It can be used to modify the event or the config to be used.
// If you want to prevent the error from being sent to bugsnag, return an error that
// explains why the notification was cancelled.
func OnBeforeNotify(callback func(event *Event, config *Configuration) error) {
	middleware.OnBeforeNotify(callback)
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

// HandlerFunc wraps a HTTP handler in bugsnag.AutoNotify(). It includes details about the
// HTTP request in all error reports. If you've wrapped your server in an http.Handler,
// you don't also need to wrap each function.
func HandlerFunc(h http.HandlerFunc, rawData ...interface{}) http.HandlerFunc {
	notifier := NewNotifier(rawData...)

	return func(w http.ResponseWriter, r *http.Request) {
		defer notifier.AutoNotify(r)
		h(w, r)
	}
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
		PanicHandler:        defaultPanicHandler,
		Transport:           http.DefaultTransport,
	})

	hostname, err := os.Hostname()
	if err == nil {
		Config.Hostname = hostname
	}
}
