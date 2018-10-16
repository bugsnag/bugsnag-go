package bugsnag_test

import (
	"fmt"
	"testing"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/bugsnag/bugsnag-go"
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

	// Expect the following frames to be present for *.Notify
	/*
		{ "file": "$GOPATH/src/github.com/bugsnag/bugsnag-go/notifier_test.go", "method": "TestStackframesAreSkippedCorrectly.func1" },
		{ "file": "testing/testing.go", "method": "tRunner" },
		{ "file": "runtime/asm_amd64.s", "method": "goexit" }
	*/

	t.Run("notifier.Notify", func(st *testing.T) {
		notifier.Notify(fmt.Errorf("oopsie"))
		assertStackframeCount(st, 3)
	})
	t.Run("bugsnag.Notify", func(st *testing.T) {
		bugsnag.Notify(fmt.Errorf("oopsie"))
		assertStackframeCount(st, 3)
	})

	// Expect the following frames to be present for notifier.NotifySync
	/*
		{ "file": "$GOPATH/src/github.com/bugsnag/bugsnag-go/notifier_test.go", "method": "TestStackframesAreSkippedCorrectly.func2" },
		{ "file": "testing/testing.go", "method": "tRunner" },
		{ "file": "runtime/asm_amd64.s", "method": "goexit" }
	*/

	t.Run("notifier.NotifySync", func(st *testing.T) {
		notifier.NotifySync(fmt.Errorf("oopsie"), true)
		assertStackframeCount(st, 3)
	})

	// Expect the following frames to be present for *.AutoNotify
	/*
		{ "file": "runtime/panic.go", "method": "gopanic" },
		{ "file": "runtime/iface.go", "method": "panicdottypeE" },
		{ "file": "$GOPATH/src/github.com/bugsnag/bugsnag-go/notifier_test.go", "method": "TestStackframesAreSkippedCorrectly.func2.1" },
		{ "file": "$GOPATH/src/github.com/bugsnag/bugsnag-go/notifier_test.go", "method": "TestStackframesAreSkippedCorrectly.func3" },
		{ "file": "testing/testing.go", "method": "tRunner" },
		{ "file": "runtime/asm_amd64.s", "method": "goexit" }
	*/
	t.Run("notifier.AutoNotify", func(st *testing.T) {
		func() {
			defer func() { recover() }()
			defer notifier.AutoNotify()
			crash("NaN")
		}()
		assertStackframeCount(st, 6)
	})
	t.Run("bugsnag.AutoNotify", func(st *testing.T) {
		func() {
			defer func() { recover() }()
			defer bugsnag.AutoNotify()
			crash("NaN")
		}()
		assertStackframeCount(st, 6)
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
		assertStackframeCount(st, 6)
	})
	t.Run("bugsnag.Recover", func(st *testing.T) {
		func() {
			defer bugsnag.Recover()
			crash("NaN")
		}()
		assertStackframeCount(st, 6)
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
