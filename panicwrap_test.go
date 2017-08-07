// +build !appengine

package bugsnag

import (
	"github.com/bitly/go-simplejson"
	"github.com/kardianos/osext"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestPanicHandler(t *testing.T) {
	startTestServer()

	exePath, err := osext.Executable()
	if err != nil {
		t.Fatal(err)
	}

	// Use the same trick as panicwrap() to re-run ourselves.
	// In the init() block below, we will then panic.
	cmd := exec.Command(exePath, os.Args[1:]...)
	cmd.Env = append(os.Environ(), "BUGSNAG_API_KEY="+testAPIKey, "BUGSNAG_ENDPOINT="+testEndpoint, "please_panic=please_panic")

	if err = cmd.Start(); err != nil {
		t.Fatal(err)
	}

	if err = cmd.Wait(); err.Error() != "exit status 2" {
		t.Fatal(err)
	}

	json, err := simplejson.NewJson(<-postedJSON)
	if err != nil {
		t.Fatal(err)
	}

	event := json.Get("events").GetIndex(0)

	if event.Get("severity").MustString() != "error" {
		t.Errorf("severity should be error")
	}
	exception := event.Get("exceptions").GetIndex(0)

	message := exception.Get("message").MustString()
	if message != "ruh roh" {
		t.Errorf("caught wrong panic message: '%s'", message)
	}

	errorClass := exception.Get("errorClass").MustString()
	if errorClass != "*errors.errorString" {
		t.Errorf("caught wrong panic errorClass: '%s'", errorClass)
	}

	stacktrace := exception.Get("stacktrace")

	// Yeah, we just caught a panic from the init() function below and sent it to the server running above (mindblown)
	frame := stacktrace.GetIndex(1)
	if frame.Get("inProject").MustBool() != true ||
		frame.Get("file").MustString() != "panicwrap_test.go" ||
		frame.Get("lineNumber").MustInt() == 0 {
		t.Errorf("stack frame seems wrong at index 1: %v", frame)
	}
}

func init() {
	if os.Getenv("please_panic") != "" {
		Configure(Configuration{APIKey: os.Getenv("BUGSNAG_API_KEY"), Endpoint: os.Getenv("BUGSNAG_ENDPOINT"), ProjectPackages: []string{"github.com/bugsnag/bugsnag-go"}})
		go func() {
			defer AutoNotify()

			panick()
		}()
		// Plenty of time to crash, it shouldn't need any of it.
		time.Sleep(1 * time.Second)
	}
}

func panick() {
	panic("ruh roh")
}
