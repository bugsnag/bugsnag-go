package bugsnag

import (
	"testing"
)

func TestMiddlewareOrder(t *testing.T) {

	result := make([]int, 0, 7)
	stack := MiddlewareStack{}
	stack.AddMiddleware(func(e *Event, c *Configuration, next func()) {
		result = append(result, 1)
		next()
		result = append(result, 7)
	})
	stack.AddMiddleware(func(e *Event, c *Configuration, next func()) {
		result = append(result, 2)
		next()
		result = append(result, 6)
	})
	stack.AddMiddleware(func(e *Event, c *Configuration, next func()) {
		result = append(result, 3)
		next()
		result = append(result, 5)
	})

	stack.Run(nil, nil, func() {
		result = append(result, 4)
	})

	if !(result[0] == 1 && result[1] == 2 && result[2] == 3 &&
		result[3] == 4 && result[4] == 5 && result[5] == 6 && result[6] == 7) {
		t.Errorf("unexpected middleware order %%", result)
	}
}

func TestBeforeNotifyReturnFalse(t *testing.T) {

	stack := MiddlewareStack{}

	stack.BeforeNotify(func(e *Event, c *Configuration) bool {
		return false
	})

	called := false

	stack.Run(nil, nil, func() {
		called = true
	})

	if called == true {
		t.Errorf("Notify was called when BeforeNotify returned False")
	}
}

func TestBeforeNotifyReturnTrue(t *testing.T) {

	stack := MiddlewareStack{}

	stack.BeforeNotify(func(e *Event, c *Configuration) bool {
		return true
	})

	called := false

	stack.Run(nil, nil, func() {
		called = true
	})

	if called == false {
		t.Errorf("Notify was not called when BeforeNotify returned True")
	}
}

func TestPanicHandling(t *testing.T) {

	stack := MiddlewareStack{}

	stack.BeforeNotify(func(e *Event, c *Configuration) bool {
		panic("oops")
	})

	called := false

	stack.Run(nil, nil, func() {
		called = true
	})

	if called == false {
		t.Errorf("Notify was not called when BeforeNotify panicked")
	}
}
