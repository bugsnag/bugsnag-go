package bugsnag

import (
	"context"
	"net/http"
	"strings"

	"github.com/bugsnag/bugsnag-go/v2/errors"
)

// Context is the context of the error in Bugsnag.
// This can be passed to Notify, Recover or AutoNotify as rawData.
type Context struct {
	String string
}

// User represents the searchable user-data on Bugsnag. The Id is also used
// to determine the number of users affected by a bug. This can be
// passed to Notify, Recover or AutoNotify as rawData.
type User struct {
	Id    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// ErrorClass overrides the error class in Bugsnag.
// This struct enables you to group errors as you like.
type ErrorClass struct {
	Name string
}

// Sets the severity of the error on Bugsnag. These values can be
// passed to Notify, Recover or AutoNotify as rawData.
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
type StackFrame struct {
	Method     string `json:"method"`
	File       string `json:"file"`
	LineNumber int    `json:"lineNumber"`
	InProject  bool   `json:"inProject,omitempty"`
}

type SeverityReason string

const (
	SeverityReasonCallbackSpecified        SeverityReason = "userCallbackSetSeverity"
	SeverityReasonHandledError                            = "handledError"
	SeverityReasonHandledPanic                            = "handledPanic"
	SeverityReasonUnhandledError                          = "unhandledError"
	SeverityReasonUnhandledMiddlewareError                = "unhandledErrorMiddleware"
	SeverityReasonUnhandledPanic                          = "unhandledPanic"
	SeverityReasonUserSpecified                           = "userSpecifiedSeverity"
)

type HandledState struct {
	SeverityReason   SeverityReason
	OriginalSeverity severity
	Unhandled        bool
	Framework        string
}

// Event represents a payload of data that gets sent to Bugsnag.
// This is passed to each OnBeforeNotify hook.
type Event struct {

	// The original error that caused this event, not sent to Bugsnag.
	Error *errors.Error

	// The rawData affecting this error, not sent to Bugsnag.
	RawData []interface{}

	// The error class to be sent to Bugsnag. This defaults to the type name of the Error, for
	// example *error.String
	ErrorClass string
	// The error message to be sent to Bugsnag. This defaults to the return value of Error.Error()
	Message string
	// The stacktrrace of the error to be sent to Bugsnag.
	Stacktrace []StackFrame

	// The context to be sent to Bugsnag. This should be set to the part of the app that was running,
	// e.g. for http requests, set it to the path.
	Context string
	// The severity of the error. Can be SeverityError, SeverityWarning or SeverityInfo.
	Severity severity
	// The grouping hash is used to override Bugsnag's grouping. Set this if you'd like all errors with
	// the same grouping hash to group together in the dashboard.
	GroupingHash string

	// User data to send to Bugsnag. This is searchable on the dashboard.
	User *User
	// Other MetaData to send to Bugsnag. Appears as a set of tabbed tables in the dashboard.
	MetaData MetaData
	// Ctx is the context of the session the event occurred in. This allows Bugsnag to associate the event with the session.
	Ctx context.Context
	// Request is the request information that populates the Request tab in the dashboard.
	Request *RequestJSON
	// The reason for the severity and original value
	handledState HandledState
	// True if the event was caused by an automatic event
	Unhandled bool
}

func newEvent(rawData []interface{}, notifier *Notifier) (*Event, *Configuration) {
	config := notifier.Config
	event := &Event{
		RawData:  append(notifier.RawData, rawData...),
		Severity: SeverityWarning,
		MetaData: make(MetaData),
		handledState: HandledState{
			SeverityReason:   SeverityReasonHandledError,
			OriginalSeverity: SeverityWarning,
			Unhandled:        false,
			Framework:        "",
		},
		Unhandled: false,
	}

	var err *errors.Error
	var callbacks []func(*Event)

	for _, datum := range event.RawData {
		switch datum := datum.(type) {

		case error, errors.Error:
			err = errors.New(datum.(error), 1)
			event.Error = err
			// Only assign automatically if not explicitly set through ErrorClass already
			if event.ErrorClass == "" {
				event.ErrorClass = err.TypeName()
			}
			event.Message = err.Error()
			event.Stacktrace = make([]StackFrame, len(err.StackFrames()))

		case bool:
			// clone without merging to avoid double-call to Configuration.update
			config = config.clone()
			config.Synchronous = bool(datum)

		case severity:
			event.Severity = datum
			event.handledState.OriginalSeverity = datum
			event.handledState.SeverityReason = SeverityReasonUserSpecified

		case Context:
			event.Context = datum.String

		case context.Context:
			populateEventWithContext(datum, event)

		case *http.Request:
			populateEventWithRequest(datum, event)

		case Configuration:
			config = config.merge(&datum)

		case MetaData:
			event.MetaData.Update(datum)

		case User:
			event.User = &datum

		case ErrorClass:
			event.ErrorClass = datum.Name

		case HandledState:
			event.handledState = datum
			event.Severity = datum.OriginalSeverity
			event.Unhandled = datum.Unhandled
		case func(*Event):
			callbacks = append(callbacks, datum)
		}
	}

	event.Stacktrace = generateStacktrace(err, config)

	for _, callback := range callbacks {
		callback(event)
		if event.Severity != event.handledState.OriginalSeverity {
			event.handledState.SeverityReason = SeverityReasonCallbackSpecified
		}
	}

	return event, config
}

func generateStacktrace(err *errors.Error, config *Configuration) []StackFrame {
	stack := make([]StackFrame, len(err.StackFrames()))
	for i, frame := range err.StackFrames() {
		file := frame.File
		inProject := config.isProjectPackage(frame.Package)

		// This will trim path before package name for external packages and golang default packages
		// Excluding main package as it's special case
		// This will NOT trim paths for packages in current module because path won't contain the package name
		// Example: path is "/user/name/work/internal/internal.go" and module package name is "example.com/mymodule/internal"
		if idx := strings.Index(file, frame.Package); idx > -1 && frame.Package != "main" {
			file = file[idx:]
		}

		// This should trim path for main and other current module packages with correct config
		// If input path is "/user/name/work/internal/internal.go"
		// SourceRoot is "/user/name/work" and ProjectPackages are []string{"main*", "example.com/mymodule/**"}
		// Then when package name is "example.com/mymodule/internal"
		// The path will be trimmed to "/internal/internal.go"
		if inProject {
			file = config.stripProjectPackages(file)
		}

		stack[i] = StackFrame{
			Method:     frame.Name,
			File:       file,
			LineNumber: frame.LineNumber,
			InProject:  inProject,
		}
	}

	return stack
}

func populateEventWithContext(ctx context.Context, event *Event) {
	event.Ctx = ctx
	reqJSON, req := extractRequestInfo(ctx)
	if event.Request == nil {
		event.Request = reqJSON
	}
	populateEventWithRequest(req, event)

}

func populateEventWithRequest(req *http.Request, event *Event) {
	if req == nil {
		return
	}

	event.Request = extractRequestInfoFromReq(req)

	if event.Context == "" {
		event.Context = req.URL.Path
	}

	// Default user.id to IP so that the count of users affected works.
	if event.User == nil {
		ip := req.RemoteAddr
		if idx := strings.LastIndex(ip, ":"); idx != -1 {
			ip = ip[:idx]
		}
		event.User = &User{Id: ip}
	}
}
