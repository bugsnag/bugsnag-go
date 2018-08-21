package bugsnag

import (
	"fmt"

	"github.com/bugsnag/bugsnag-go/errors"
)

// Notifier sends errors to Bugsnag.
type Notifier struct {
	Config  *Configuration
	RawData []interface{}
}

// New creates a new notifier.
// You can pass an instance of bugsnag.Configuration in rawData to change the configuration.
// Other values of rawData will be passed to Notify.
func New(rawData ...interface{}) *Notifier {
	config := Config.clone()
	for i, datum := range rawData {
		if c, ok := datum.(Configuration); ok {
			config.update(&c)
			rawData[i] = nil
		}
	}

	return &Notifier{
		Config:  config,
		RawData: rawData,
	}
}

// Notify sends an error to Bugsnag. Any rawData you pass here will be sent to
// Bugsnag after being converted to JSON. e.g. bugsnag.SeverityError, bugsnag.Context,
// or bugsnag.MetaData.
func (notifier *Notifier) Notify(err error, rawData ...interface{}) (e error) {
	config := notifier.Config
	return notifier.NotifySync(err, config.Synchronous, rawData...)
}

// NotifySync sends an error to Bugsnag. The synchronous parameter specifies
// whether to send the report in the current context. Any rawData you pass here
// will be sent to Bugsnag after being converted to JSON. e.g.
// bugsnag.SeverityError,  bugsnag.Context, or bugsnag.MetaData.
func (notifier *Notifier) NotifySync(err error, synchronous bool, rawData ...interface{}) (e error) {
	event, config := newEvent(errors.New(err, 1), rawData, notifier)

	// Never block, start throwing away errors if we have too many.
	e = middleware.Run(event, config, func() error {
		config.logf("notifying bugsnag: %s", event.Message)
		if config.notifyInReleaseStage() {
			if synchronous {
				return config.Shipper.Deliver(&payload{event, config})
			}
			// Ensure that any errors are logged if they occur in a goroutine.
			go func(event *Event, config *Configuration) {
				err := config.Shipper.Deliver(&payload{event, config})
				if err != nil {
					config.logf("bugsnag.Notify: %v", err)
				}
			}(event, config)

			return nil
		}
		return fmt.Errorf("not notifying in %s", config.ReleaseStage)
	})

	if e != nil {
		config.logf("bugsnag.Notify: %v", e)
	}
	return e
}

// AutoNotify notifies Bugsnag of any panics, then repanics.
// It sends along any rawData that gets passed in.
// Usage:
//  go func() {
//		defer AutoNotify()
//      // (possibly crashy code)
//  }()
func (notifier *Notifier) AutoNotify(rawData ...interface{}) {
	if err := recover(); err != nil {
		severity := notifier.getDefaultSeverity(rawData, SeverityError)
		state := HandledState{SeverityReasonHandledPanic, severity, true, ""}
		notifier.appendStateIfNeeded(rawData, state)
		notifier.Notify(errors.New(err, 2), rawData...)
		panic(err)
	}
}

// Recover logs any panics, then recovers.
// It sends along any rawData that gets passed in.
// Usage: defer Recover()
func (notifier *Notifier) Recover(rawData ...interface{}) {
	if err := recover(); err != nil {
		severity := notifier.getDefaultSeverity(rawData, SeverityWarning)
		state := HandledState{SeverityReasonHandledPanic, severity, false, ""}
		notifier.appendStateIfNeeded(rawData, state)
		notifier.Notify(errors.New(err, 2), rawData...)
	}
}

func (notifier *Notifier) dontPanic() {
	if err := recover(); err != nil {
		notifier.Config.logf("bugsnag/notifier.Notify: panic! %s", err)
	}
}

// Get defined severity from raw data or a fallback value
func (notifier *Notifier) getDefaultSeverity(rawData []interface{}, s severity) severity {
	allData := append(notifier.RawData, rawData...)
	for _, datum := range allData {
		if _, ok := datum.(severity); ok {
			return datum.(severity)
		}
	}

	for _, datum := range allData {
		if _, ok := datum.(HandledState); ok {
			return datum.(HandledState).OriginalSeverity
		}
	}

	return s
}

func (notifier *Notifier) appendStateIfNeeded(rawData []interface{}, h HandledState) []interface{} {

	for _, datum := range append(notifier.RawData, rawData...) {
		if _, ok := datum.(HandledState); ok {
			return rawData
		}
	}

	return append(rawData, h)
}
