// +build !appengine

package tests

import (
	"os/exec"
	"testing"
)

// Starts an app, sends a request, and tests that the resulting bugsnag
// error report has the correct values.

func TestRevelRequestPanic(t *testing.T) {
	startTestServer()
	body := startRevelApp(t, "default")
	assertSeverityReasonEqual(t, body, "error", "unhandledErrorMiddleware", true)
	pkill("revel")
}

func TestRevelRequestPanicCallbackAltered(t *testing.T) {
	startTestServer()
	body := startRevelApp(t, "beforenotify")
	assertSeverityReasonEqual(t, body, "info", "userCallbackSetSeverity", true)
	pkill("revel")
}

func pkill(process string) {
	cmd := exec.Command("pkill", "-x", process)
	cmd.Start()
	cmd.Wait()
}
