package bugsnag

type (
	// Middleware are functions that let you modify events before they
	// are sent to Bugsnag. They should call next() to send the notification,
	// or avoid calling it to prevent the notification being sent.
	Middleware func(event *Event, config *Configuration, next func())

	// MiddlewareStacks keep middleware in the correct order. They are
	// called in reverse order, so if you add a new middleware it will
	// be called before all existing middleware.
	MiddlewareStack struct {
		middleware []Middleware
	}
)

// AddMiddleware adds a new middleware to the outside of the existing ones,
// when the MiddlewareStack is Run it will be run before all middleware that
// have been added before.
func (stack *MiddlewareStack) AddMiddleware(middleware Middleware) {
	stack.middleware = append(stack.middleware, middleware)
}

// BeforeNotify adds a new middleware that runs before any existing ones,
// it can be used to easily modify the event or abort processing.
func (stack *MiddlewareStack) BeforeNotify(middleware func(*Event, *Configuration) bool) {
	stack.AddMiddleware(func(e *Event, n *Configuration, next func()) {
		if middleware(e, n) {
			next()
		}
	})
}

// Run causes all the middleware to be run. If they all permit it the next callback
// will be called with all the middleware on the stack.
func (stack *MiddlewareStack) Run(event *Event, config *Configuration, next func()) {
	for i, _ := range stack.middleware {
		next = (func(f Middleware, next func()) func() {
			return func() {
				defer catchMiddlewarePanic(event, config, next)
				f(event, config, next)
			}
		})(stack.middleware[len(stack.middleware)-1-i], next)
	}
	next()
}

func NewMiddleware() *MiddlewareStack {
	return &MiddlewareStack{middleware: make([]Middleware, 0)}
}

func DefaultMiddleware() *MiddlewareStack {
	return NewMiddleware()
}

// catchMiddlewarePanic is used to log any panics that happen inside Middleware,
// we wouldn't want to not notify Bugsnag in this case.
func catchMiddlewarePanic(event *Event, config *Configuration, next func()) {
	if err := recover(); err != nil {
		println("TODO: Use a logger!")
		next()
	}
}
