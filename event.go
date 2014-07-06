package bugsnag

import (
	"github.com/bugsnag/bugsnag-go/errors"
	"strings"
)

// Sets the context of the error in Bugsnag. You can pass this in
// where-ever rawData is expected.
type Context struct {
	String string
}

// Sets the searchable user-data on Bugsnag. The Id is also used
// to determine the number of users affected by a bug. You can pass
// this in where-ever rawData is expected.
type User struct {
	Id     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
}

// Tags used to mark the severity of the error in the Bugsnag dashboard.
// You can pass these tags where-ever rawData is expected.
var (
	SeverityError   = severity{"error"}
	SeverityWarning = severity{"warning"}
	SeverityInfo    = severity{"info"}
)
// The severity tag type, private so that people can only use Error,Warning,Info
type severity struct {
	String string
}

// The form of stacktrace that Bugsnag expects
type stackFrame struct {
	Method     string `json:"method"`
	File       string `json:"file"`
	LineNumber int    `json:"lineNumber"`
	InProject  bool   `json:"inProject,omitempty"`
}

// An event to send to Bugsnag. This is passed through the middleware stack.
type Event struct {
	// The original error that caused this event, not sent to Bugsnag
	Error   *errors.Error
	// The rawData affecting this error, not sent to Bugsnag
	RawData []interface{}

	// The error class to be sent to Bugsnag. This defaults to the type name of the Error
	ErrorClass string
	// The error message to be sent to Bugsnag. This defaults to the return value of Error.Error()
	Message    string
	// The stacktrrace of the error to be sent to Bugsnag.
	Stacktrace []stackFrame

	// The context to be sent to Bugsnag. This should be set to the part of the app that was running,
	// e.g. for http requests, set it to the path.
	Context      string
	// The severity of the error. Can be SeverityError, SeverityWarning or SeverityInfo
	Severity     severity
	// The grouping hash is used to override Bugsnag's grouping. Set this if you'd like all errors with
	// the same grouping hash to group together in the dashboard.
	GroupingHash string

	// The searchable user data to send to Bugsnag.
	User     *User
	// Other meta-data to send to Bugsnag. Appears as a set of tabbed tables in the dashboard.
	MetaData MetaData
}

func newEvent(err *errors.Error, rawData []interface{}, notifier *Notifier) (*Event, *Configuration) {

	config := notifier.Config
	event := &Event{
		Error:   err,
		RawData: append(notifier.RawData, rawData...),

		ErrorClass: err.TypeName(),
		Message:    err.Error(),
		Stacktrace: make([]stackFrame, len(err.StackFrames())),

		Severity: SeverityWarning,

		MetaData: make(MetaData),
	}

	for _, datum := range event.RawData {
		switch datum := datum.(type) {
		case severity:
			event.Severity = datum

		case Context:
			event.Context = datum.String

		case Configuration:
			config = config.merge(&datum)

		case MetaData:
			event.MetaData.Update(datum)

		case User:
			event.User = &datum
		}
	}

	for i, frame := range err.StackFrames() {
		file := frame.File
		// make in-project frames really nice.
		file = strings.TrimPrefix(file, config.ProjectRoot)
		// remove $GOROOT and $GOHOME from other frames
		if idx := strings.Index(file, frame.Package); idx > -1 {
			file = file[idx:]
		}

		event.Stacktrace[i] = stackFrame{
			Method:     frame.Name,
			File:       file,
			LineNumber: frame.LineNumber,
			InProject:  config.isProjectPackage(frame.Package),
		}
	}

	return event, config
}
