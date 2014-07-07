package bugsnag

import (
	"github.com/bugsnag/bugsnag-go/errors"
)

type Notifier struct {
	Config  *Configuration
	RawData []interface{}
}

// Creates a new notifier. You can pass an instance of bugsnag.Configuration
// in rawData to change the configuration. Other values of rawData will be
// passed to Notify.
func NewNotifier(rawData ...interface{}) *Notifier {
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
func (notifier *Notifier) Notify(err error, rawData ...interface{}) {
	defer notifier.dontPanic()
	event, config := newEvent(errors.New(err, 2), rawData, notifier)

	// Never block, start throwing away errors if we have too many.
	middleware.Run(event, config, func() {
		defer notifier.dontPanic()
		if config.notifyInReleaseStage() {
			(&payload{event, config}).deliver()
		}
	})
}

// defer AutoNotify() sends any panics that happen to Bugsnag, along with any
// rawData you set here. After the notification is sent, panic() is called again
// with the same error so that the panic() bubbles out.
func (notifier *Notifier) AutoNotify(rawData ...interface{}) {
	if err := recover(); err != nil {
		rawData = append(rawData, SeverityError)
		notifier.Notify(errors.New(err, 2), rawData...)
		panic(err)
	}
}

// defer AutoNotify() sends any panics that happen to Bugsnag, along with any
// rawData you set here. After the notification is sent, the panic() is considered
// to have been recovered() and execution proceeds as normal.
func (notifier *Notifier) Recover(rawData ...interface{}) {
	if err := recover(); err != nil {
		rawData = append(rawData, SeverityError)
		notifier.Notify(errors.New(err, 2), rawData...)
	}
}

func (notifier *Notifier) dontPanic() {
	if err := recover(); err != nil {
		notifier.Config.log("bugsnag/notifier.Notify: panic! %s", err)
	}
}
