package bugsnag_test

import (
	"fmt"
	"strings"
	"testing"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/bugsnag/bugsnag-go/v2"
	"github.com/bugsnag/bugsnag-go/v2/errors"
	. "github.com/bugsnag/bugsnag-go/v2/testutil"
)

var bugsnaggedReports chan []byte

func notifierSetup(url string) *bugsnag.Notifier {
	return bugsnag.New(bugsnag.Configuration{
		APIKey:    TestAPIKey,
		Endpoints: bugsnag.Endpoints{Notify: url, Sessions: url + "/sessions"},
	})
}

func crash(s interface{}) int {
	return s.(int)
}

func TestStackframesAreSkippedCorrectly(t *testing.T) {
	ts, reports := Setup()
	bugsnaggedReports = reports
	defer ts.Close()
	notifier := notifierSetup(ts.URL)

	bugsnag.Configure(bugsnag.Configuration{
		APIKey:    TestAPIKey,
		Endpoints: bugsnag.Endpoints{Notify: ts.URL, Sessions: ts.URL + "/sessions"},
	})

	t.Run("notifier.Notify", func(st *testing.T) {
		notifier.Notify(fmt.Errorf("oopsie"))
		assertStackframesMatch(t, []errors.StackFrame{
			errors.StackFrame{Name: "TestStackframesAreSkippedCorrectly.func1", File: "notifier_test.go"},
		})
	})
	t.Run("bugsnag.Notify", func(st *testing.T) {
		bugsnag.Notify(fmt.Errorf("oopsie"))
		assertStackframesMatch(t, []errors.StackFrame{
			errors.StackFrame{Name: "TestStackframesAreSkippedCorrectly.func2", File: "notifier_test.go"},
		})
	})

	t.Run("notifier.NotifySync", func(st *testing.T) {
		notifier.NotifySync(fmt.Errorf("oopsie"), true)
		assertStackframesMatch(t, []errors.StackFrame{
			errors.StackFrame{Name: "TestStackframesAreSkippedCorrectly.func3", File: "notifier_test.go"},
		})
	})

	t.Run("notifier.AutoNotify", func(st *testing.T) {
		func() {
			defer func() { recover() }()
			defer notifier.AutoNotify()
			crash("NaN")
		}()
		assertStackframesMatch(t, []errors.StackFrame{
			errors.StackFrame{Name: "TestStackframesAreSkippedCorrectly.func4.1", File: "notifier_test.go"},
			errors.StackFrame{Name: "TestStackframesAreSkippedCorrectly.func4", File: "notifier_test.go"},
		})
	})
	t.Run("bugsnag.AutoNotify", func(st *testing.T) {
		func() {
			defer func() { recover() }()
			defer bugsnag.AutoNotify()
			crash("NaN")
		}()
		assertStackframesMatch(t, []errors.StackFrame{
			errors.StackFrame{Name: "TestStackframesAreSkippedCorrectly.func5.1", File: "notifier_test.go"},
			errors.StackFrame{Name: "TestStackframesAreSkippedCorrectly.func5", File: "notifier_test.go"},
		})
	})

	// Expect the following frames to be present for *.Recover
	/*
		{ "file": "runtime/panic.go", "method": "gopanic" },
		{ "file": "runtime/iface.go", "method": "panicdottypeE" },
		{ "file": "$GOPATH/src/github.com/bugsnag/bugsnag-go/notifier_test.go", "method": "TestStackframesAreSkippedCorrectly.func4.1" },
		{ "file": "$GOPATH/src/github.com/bugsnag/bugsnag-go/notifier_test.go", "method": "TestStackframesAreSkippedCorrectly.func4" },
		{ "file": "testing/testing.go", "method": "tRunner" },
		{ "file": "runtime/asm_amd64.s", "method": "goexit" }
	*/
	t.Run("notifier.Recover", func(st *testing.T) {
		func() {
			defer notifier.Recover()
			crash("NaN")
		}()
		assertStackframesMatch(t, []errors.StackFrame{
			errors.StackFrame{Name: "TestStackframesAreSkippedCorrectly.func6.1", File: "notifier_test.go"},
			errors.StackFrame{Name: "TestStackframesAreSkippedCorrectly.func6", File: "notifier_test.go"},
		})
	})
	t.Run("bugsnag.Recover", func(st *testing.T) {
		func() {
			defer bugsnag.Recover()
			crash("NaN")
		}()
		assertStackframesMatch(t, []errors.StackFrame{
			errors.StackFrame{Name: "TestStackframesAreSkippedCorrectly.func7.1", File: "notifier_test.go"},
			errors.StackFrame{Name: "TestStackframesAreSkippedCorrectly.func7", File: "notifier_test.go"},
		})
	})
}

func TestModifyingEventsWithCallbacks(t *testing.T) {
	server, eventQueue := Setup()
	defer server.Close()
	notifier := notifierSetup(server.URL)

	bugsnag.Configure(bugsnag.Configuration{
		APIKey:    TestAPIKey,
		Endpoints: bugsnag.Endpoints{Notify: server.URL, Sessions: server.URL + "/sessions"},
	})

	t.Run("bugsnag.Notify change unhandled in block", func(st *testing.T) {
		notifier.Notify(fmt.Errorf("ahoy"), func(event *bugsnag.Event) {
			event.Unhandled = true
		})
		json, _ := simplejson.NewJson(<-eventQueue)
		event := GetIndex(json, "events", 0)
		exception := GetIndex(event, "exceptions", 0)
		message := exception.Get("message").MustString()
		unhandled := event.Get("unhandled").MustBool()
		overridden := event.Get("severityReason").Get("unhandledOverridden").MustBool()
		if message != "ahoy" {
			st.Errorf("incorrect error message '%s'", message)
		}
		if !unhandled {
			st.Errorf("failed to change handled-ness in block")
		}
		if !overridden {
			st.Errorf("failed to set handledness change in block")
		}
	})

	t.Run("bugsnag.Notify with block", func(st *testing.T) {
		notifier.Notify(fmt.Errorf("bnuuy"), bugsnag.Context{String: "should be overridden"}, func(event *bugsnag.Event) {
			event.Context = "known unknowns"
		})
		json, _ := simplejson.NewJson(<-eventQueue)
		event := GetIndex(json, "events", 0)
		context := event.Get("context").MustString()
		exception := GetIndex(event, "exceptions", 0)
		class := exception.Get("errorClass").MustString()
		message := exception.Get("message").MustString()
		if class != "*errors.errorString" {
			st.Errorf("incorrect error class '%s'", class)
		}
		if message != "bnuuy" {
			st.Errorf("incorrect error message '%s'", message)
		}
		if context != "known unknowns" {
			st.Errorf("failed to change context in block. '%s'", context)
		}
		if event.Get("unhandled").MustBool() {
			st.Errorf("error is unexpectedly unhandled")
		}
		if overridden, err := event.Get("severityReason").Get("unhandledOverridden").Bool(); err == nil {
			// if err == nil, then the value existed in the payload. the expectation
			// is that unhandledOverridden is not sent when handled-ness is not changed.
			st.Errorf("error unexpectedly has unhandledOverridden: %v", overridden)
		}
	})
}

func assertStackframesMatch(t *testing.T, expected []errors.StackFrame) {
	var lastmatch int = 0
	var matched int = 0
	event, _ := simplejson.NewJson(<-bugsnaggedReports)
	json := GetIndex(event, "events", 0)
	stacktrace := GetIndex(json, "exceptions", 0).Get("stacktrace")
	for i := 0; i < len(stacktrace.MustArray()); i++ {
		actualFrame := stacktrace.GetIndex(i)
		file := actualFrame.Get("file").MustString()
		method := actualFrame.Get("method").MustString()
		for index, expectedFrame := range expected {
			if index < lastmatch {
				continue
			}
			if strings.HasSuffix(file, expectedFrame.File) && expectedFrame.Name == method {
				lastmatch = index
				matched++
			}
		}
	}

	if matched != len(expected) {
		s, _ := stacktrace.EncodePretty()
		t.Errorf("failed to find matches for %d frames: '%v'\ngot: '%v'", len(expected)-matched, expected[matched:], string(s))
	}
}
