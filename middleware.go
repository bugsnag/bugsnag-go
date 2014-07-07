package bugsnag

import (
	"net/http"
	"strings"
)

type (
	middlewareFunc func(event *Event, config *Configuration, next func())

	beforeNotifyFunc func(*Event, *Configuration) bool

	// MiddlewareStacks keep middleware in the correct order. They are
	// called in reverse order, so if you add a new middleware it will
	// be called before all existing middleware.
	middlewareStack struct {
		middleware []middlewareFunc
	}
)

// AddMiddleware adds a new middleware to the outside of the existing ones,
// when the middlewareStack is Run it will be run before all middleware that
// have been added before.
func (stack *middlewareStack) AddMiddleware(middleware middlewareFunc) {
	stack.middleware = append(stack.middleware, middleware)
}

// BeforeNotify adds a new middleware that runs before any existing ones,
// it can be used to easily modify the event or abort processing.
func (stack *middlewareStack) BeforeNotify(middleware beforeNotifyFunc) {
	stack.AddMiddleware(func(e *Event, n *Configuration, next func()) {
		if middleware(e, n) {
			next()
		}
	})
}

// Run causes all the middleware to be run. If they all permit it the next callback
// will be called with all the middleware on the stack.
func (stack *middlewareStack) Run(event *Event, config *Configuration, next func()) {
	for _, middleware := range stack.middleware {
		next = (func(f middlewareFunc, next func()) func() {
			return func() {
				defer catchMiddlewarePanic(event, config, next)
				f(event, config, next)
			}
		})(middleware, next)
	}
	next()
}

// catchMiddlewarePanic is used to log any panics that happen inside Middleware,
// we wouldn't want to not notify Bugsnag in this case.
func catchMiddlewarePanic(event *Event, config *Configuration, next func()) {
	if err := recover(); err != nil {
		config.log("bugsnag/middleware: unexpected panic: %v", err)
		next()
	}
}

// httpRequestMiddleware is added OnBeforeNotify by default. It takes information
// from an http.Request passed in as rawData, and adds it to the Event. You can
// use this as a template for writing your own Middleware.
func httpRequestMiddleware(event *Event, config *Configuration) bool {
	for _, datum := range event.RawData {
		if request, ok := datum.(*http.Request); ok {
			proto := "http://"
			if request.TLS != nil {
				proto = "https://"
			}

			event.MetaData.Update(MetaData{
				"Request": {
					"RemoteAddr": request.RemoteAddr,
					"Method":     request.Method,
					"Url":        proto + request.Host + request.RequestURI,
					"Params":     request.URL.Query(),
				},
			})

			// Add headers as a separate tab.
			event.MetaData.AddStruct("Headers", request.Header)

			// Default context to Path
			if event.Context == "" {
				event.Context = request.URL.Path
			}

			// Default user.id to IP so that users-affected works.
			if event.User == nil {
				ip := request.RemoteAddr
				if idx := strings.LastIndex(ip, ":"); idx != -1 {
					ip = ip[:idx]
				}
				event.User = &User{Id: ip}
			}
		}
	}
	return true
}
