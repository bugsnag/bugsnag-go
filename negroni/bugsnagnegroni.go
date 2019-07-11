package bugsnagnegroni

import (
	"net/http"

	"github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/device"
	"github.com/urfave/negroni"
)

// FrameworkName is the name of the framework this middleware applies to
const FrameworkName string = "Negroni"

type handler struct {
	rawData []interface{}
}

// AutoNotify sends any panics to bugsnag, and then re-raises them.
func AutoNotify(rawData ...interface{}) negroni.Handler {
	updateGlobalConfig(rawData...)
	device.AddVersion(FrameworkName, "unknown") // Negroni exposes no version prop.
	state := bugsnag.HandledState{
		SeverityReason:   bugsnag.SeverityReasonUnhandledMiddlewareError,
		OriginalSeverity: bugsnag.SeverityError,
		Unhandled:        true,
		Framework:        FrameworkName,
	}
	rawData = append(rawData, state)
	return &handler{rawData: rawData}
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Record a session if auto capture sessions is enabled
	ctx := bugsnag.AttachRequestData(r.Context(), r)
	if bugsnag.Config.IsAutoCaptureSessions() {
		ctx = bugsnag.StartSession(ctx)
	}
	request := r.WithContext(ctx)
	notifier := bugsnag.New(h.rawData...)
	notifier.FlushSessionsOnRepanic(false)
	defer notifier.AutoNotify(ctx)
	next(rw, request)

}

func updateGlobalConfig(rawData ...interface{}) {
	for i, datum := range rawData {
		if c, ok := datum.(bugsnag.Configuration); ok {
			bugsnag.Configure(c)
			rawData[i] = nil
		}
	}
}
