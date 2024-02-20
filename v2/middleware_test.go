package bugsnag

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/bugsnag/bugsnag-go/v2/errors"
)

func TestMiddlewareOrder(t *testing.T) {

	err := fmt.Errorf("test")
	data := []interface{}{errors.New(err, 1)}
	event, config := newEvent(data, &defaultNotifier)

	result := make([]int, 0, 7)
	stack := middlewareStack{}
	stack.OnBeforeNotify(func(e *Event, c *Configuration) error {
		result = append(result, 2)
		return nil
	})
	stack.OnBeforeNotify(func(e *Event, c *Configuration) error {
		result = append(result, 1)
		return nil
	})
	stack.OnBeforeNotify(func(e *Event, c *Configuration) error {
		result = append(result, 0)
		return nil
	})

	stack.Run(event, config, func() error {
		result = append(result, 3)
		return nil
	})

	if !reflect.DeepEqual(result, []int{0, 1, 2, 3}) {
		t.Errorf("unexpected middleware order %v", result)
	}
}

func TestBeforeNotifyReturnErr(t *testing.T) {

	stack := middlewareStack{}
	err := fmt.Errorf("test")
	data := []interface{}{errors.New(err, 1)}
	event, config := newEvent(data, &defaultNotifier)

	stack.OnBeforeNotify(func(e *Event, c *Configuration) error {
		return err
	})

	called := false

	e := stack.Run(event, config, func() error {
		called = true
		return nil
	})

	if e != err {
		t.Errorf("Middleware didn't return the error")
	}

	if called == true {
		t.Errorf("Notify was called when BeforeNotify returned False")
	}
}

func TestBeforeNotifyPanic(t *testing.T) {

	stack := middlewareStack{}
	err := fmt.Errorf("test")
	event, _ := newEvent([]interface{}{errors.New(err, 1)}, &defaultNotifier)

	stack.OnBeforeNotify(func(e *Event, c *Configuration) error {
		panic("oops")
	})

	called := false
	b := &bytes.Buffer{}

	stack.Run(event, &Configuration{Logger: log.New(b, log.Prefix(), 0)}, func() error {
		called = true
		return nil
	})

	logged := b.String()

	if logged != "bugsnag/middleware: unexpected panic: oops\n" {
		t.Errorf("Logged: %s", logged)
	}

	if called == false {
		t.Errorf("Notify was not called when BeforeNotify panicked")
	}
}

func TestHttpRequestMiddleware(t *testing.T) {
	var req *http.Request
	rawData := []interface{}{req}

	event := &Event{RawData: rawData}
	config := &Configuration{}
	err := httpRequestMiddleware(event, config)
	if err != nil {
		t.Errorf("Should not happen")
	}
}
