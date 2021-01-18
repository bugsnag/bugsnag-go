package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/bugsnag/bugsnag-go/v2"
)

var testcase = flag.String("test", "", "the error scenario to run")

func main() {
	bugsnag.Configure(bugsnag.Configuration{})

	// Increase publish rate for testing
	bugsnag.DefaultSessionPublishInterval = time.Millisecond * 50

	flag.Parse()

	switch *testcase {
	case "panic":
		explicitPanic()
	case "handled":
		handledEvent()
	case "handled-metadata":
		handledMetadata()
	case "no-op":
		// nothing to see here
	default:
		fmt.Printf("No test case found for '%s'\n", *testcase)
	}
	time.Sleep(time.Millisecond * 100) // time to send before termination
}
