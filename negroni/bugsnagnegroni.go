package bugsnagnegroni

import (
	"net/http"

	"github.com/bugsnag/bugsnag-go"
)

// FrameworkName is the name of the framework this middleware applies to
const FrameworkName string = "Negroni"

type handler struct {
	rawData []interface{}
}

// AutoNotify sends any panics to bugsnag, and then re-raises them.
func AutoNotify(rawData ...interface{}) *handler {
	state := bugsnag.HandledState{
		SeverityReason:   bugsnag.SeverityReasonUnhandledMiddlewareError,
		OriginalSeverity: bugsnag.SeverityError,
		Unhandled:        true,
		Framework:        FrameworkName,
	}
	rawData = append(rawData, state)
	return &handler{
		rawData: rawData,
	}
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	request := r
	// Record a session if auto capture sessions is enabled
	if bugsnag.Config.IsAutoCaptureSessions() {
		ctx := bugsnag.StartSession(r.Context())
		request = r.WithContext(ctx)
	}

	notifier := bugsnag.New(append(h.rawData, request)...)
	defer notifier.AutoNotify(request)
	next(rw, request)
}
