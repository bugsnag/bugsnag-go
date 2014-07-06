package bugsnag

import (
	"github.com/bugsnag/bugsnag-go/errors"
)

var backgroundQueue = make(chan func(), 10)

func init() {
	go func() {
		for job := range backgroundQueue {
			job()
		}
	}()
}

type Notifier struct {
	Config  *Configuration
	RawData []interface{}
}

// Creates a new notifier with different configuration or data
func NewNotifier(config Configuration, rawData ...interface{}) *Notifier {
	return &Notifier{
		Config:  Config.merge(&config),
		RawData: rawData,
	}
}

// Notify sends an error to Bugsnag. Any extraData you pass here will be sent to
// Bugsnag after being converted to JSON.
func (notifier *Notifier) Notify(err error, rawData ...interface{}) {
	defer notifier.dontPanic()
	event, config := newEvent(errors.New(err, 1), rawData, notifier)

	// Never block, start throwing away errors if we have too many.
	if len(backgroundQueue) < cap(backgroundQueue) {
		backgroundQueue <- func() {
			middleware.Run(event, config, func() {
				defer notifier.dontPanic()
				(&payload{event, config}).deliver()
			})
		}
	} else {
		notifier.Config.Logger.Println("bugsnag/notifier.Notify: discarding error due to long queue")
	}
}

// defer AutoNotify() sends any panics that happen to Bugsnag, along with any
// metaData you set here. After the notification is sent, panic() is called again
// with the same error so that the panic() bubbles out.
func (notifier *Notifier) AutoNotify(rawData ...interface{}) {
	if err := recover(); err != nil {
		rawData = append(rawData, SeverityError)
		defaultNotifier.Notify(errors.New(err, 2), rawData...)
		panic(err)
	}
}

// defer AutoNotify() sends any panics that happen to Bugsnag, along with any
// metaData you set here. After the notification is sent, the panic() is considered
// to have been recovered() and execution proceeds as normal.
func (notifier *Notifier) Recover(rawData ...interface{}) {
	if err := recover(); err != nil {
		rawData = append(rawData, SeverityError)
		defaultNotifier.Notify(errors.New(err, 2), rawData...)
	}
}

func (notifier *Notifier) dontPanic() {
	if err := recover(); err != nil {
		notifier.Config.Logger.Println("bugsnag/notifier.Notify: panic! %s", err)
	}
}
