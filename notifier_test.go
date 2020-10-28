package bugsnag_test

import (
	"fmt"
	"strings"
	"testing"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/errors"
	. "github.com/bugsnag/bugsnag-go/testutil"
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

func assertStackframeCount(t *testing.T, expCount int) {
	report, _ := simplejson.NewJson(<-bugsnaggedReports)
	stacktrace := GetIndex(GetIndex(report, "events", 0), "exceptions", 0).Get("stacktrace")
	if s := stacktrace.MustArray(); len(s) != expCount {
		t.Errorf("Expected %d stackframe(s), but there were %d stackframes", expCount, len(s))
		s, _ := stacktrace.EncodePretty()
		t.Errorf(string(s))
	}
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
