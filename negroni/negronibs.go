package negronibs

import(
	"github.com/bugsnag/bugsnag-go"
	"net/http"
)

type handler struct {
	notifier *bugsnag.Notifier
}

func AutoNotify(rawData ...interface{}) *handler {
	n := bugsnag.New(rawData)
	handle := &handler{
		notifier: n,
	}
	return handle
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	notifier := bugsnag.New(append(append(h.notifier.RawData, h.notifier.Config), r)...)
	defer notifier.AutoNotify(r)
	next(rw, r)
}