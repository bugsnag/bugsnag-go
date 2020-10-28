package bugsnag

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/bugsnag/bugsnag-go/sessions"
	"github.com/kardianos/osext"
)

// Test the panic handler by launching a new process which runs the init()
// method in this file and causing a handled panic
func TestPanicHandlerHandledPanic(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()

	startPanickingProcess(t, "handled", ts.URL)

	json, err := simplejson.NewJson(<-reports)
	if err != nil {
		t.Fatal(err)
	}

	assertPayload(t, json, eventJSON{
		App:            &appJSON{},
		Context:        "",
		Device:         &deviceJSON{Hostname: "web1"},
		GroupingHash:   "",
		Session:        &sessionJSON{Events: sessions.EventCounts{Handled: 0, Unhandled: 1}},
		Severity:       "error",
		SeverityReason: &severityReasonJSON{Type: SeverityReasonHandledPanic},
		Unhandled:      true,
		Request:        &RequestJSON{},
		User:           &User{},
		Exceptions:     []exceptionJSON{{ErrorClass: "*errors.errorString", Message: "ruh roh"}},
	})

	event := getIndex(json, "events", 0)
	assertValidSession(t, event, true)

	stacktrace := getIndex(event, "exceptions", 0).Get("stacktrace")
	found := false
	for i := 0; i < len(stacktrace.MustArray()); i++ {
		frame := stacktrace.GetIndex(i)
		if strings.HasSuffix(getString(frame, "file"), "panicwrap_test.go") && getInt(frame, "lineNumber") != 0 {
			found = true
			break
		}
	}
	if !found {
		s, _ := stacktrace.EncodePretty()
		t.Errorf("no stack frame found matching this file in stack trace: %v", string(s))
	}
}

// Test the panic handler by launching a new process which runs the init()
// method in this file and causing an unhandled panic
func TestPanicHandlerUnhandledPanic(t *testing.T) {
	ts, reports := setup()
	defer ts.Close()

	startPanickingProcess(t, "unhandled", ts.URL)
	json, err := simplejson.NewJson(<-reports)
	if err != nil {
		t.Fatal(err)
	}
	assertPayload(t, json, eventJSON{
		App:            &appJSON{},
		Context:        "",
		Device:         &deviceJSON{Hostname: "web1"},
		GroupingHash:   "",
		Session:        &sessionJSON{Events: sessions.EventCounts{Handled: 0, Unhandled: 1}},
		Severity:       "error",
		SeverityReason: &severityReasonJSON{Type: SeverityReasonUnhandledPanic},
		Unhandled:      true,
		Request:        &RequestJSON{},
		User:           &User{},
		Exceptions:     []exceptionJSON{{ErrorClass: "panic", Message: "ruh roh"}},
	})
}

func startPanickingProcess(t *testing.T, variant string, endpoint string) {
	exePath, err := osext.Executable()
	if err != nil {
		t.Fatal(err)
	}

	// Use the same trick as panicwrap() to re-run ourselves.
	// In the init() block below, we will then panic.
	cmd := exec.Command(exePath, os.Args[1:]...)
	cmd.Env = append(os.Environ(), "BUGSNAG_API_KEY="+testAPIKey, "BUGSNAG_NOTIFY_ENDPOINT="+endpoint, "please_panic="+variant)

	// Gift for the debugging developer:
	// As these tests shell out we don't see, or even want to see, the output
	// of these tests by default.  The following two lines may be uncommented
	// in order to see what this command would print to stdout and stderr.
	/*
		bytes, _ := cmd.CombinedOutput()
		fmt.Println(string(bytes))
	*/

	if err = cmd.Start(); err != nil {
		t.Fatal(err)
	}

	if err = cmd.Wait(); err.Error() != "exit status 2" {
		t.Fatal(err)
	}
}

func init() {
	if os.Getenv("please_panic") == "handled" {
		Configure(Configuration{
			APIKey:          os.Getenv("BUGSNAG_API_KEY"),
			Endpoints:       Endpoints{Notify: os.Getenv("BUGSNAG_NOTIFY_ENDPOINT")},
			Hostname:        "web1",
			ProjectPackages: []string{"github.com/bugsnag/bugsnag-go"}})
		go func() {
			ctx := StartSession(context.Background())
			defer AutoNotify(ctx)

			panick()
		}()
		// Plenty of time to crash, it shouldn't need any of it.
		time.Sleep(1 * time.Second)
	} else if os.Getenv("please_panic") == "unhandled" {
		Configure(Configuration{
			APIKey:          os.Getenv("BUGSNAG_API_KEY"),
			Endpoints:       Endpoints{Notify: os.Getenv("BUGSNAG_NOTIFY_ENDPOINT")},
			Hostname:        "web1",
			Synchronous:     true,
			ProjectPackages: []string{"github.com/bugsnag/bugsnag-go"}})
		panick()
	}
}

func panick() {
	panic("ruh roh")
}
