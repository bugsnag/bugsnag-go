package main

import (
	"flag"
	"fmt"
	"os"

	bugsnag "github.com/bugsnag/bugsnag-go"
)

func main() {
	testcase := flag.String("case", "", "test case to run")
	flag.Parse()

	switch *testcase {
	case "default":
		caseDefault()
	case "app version":
		caseAppVersion()
	case "app type":
		caseAppType()
	case "legacy endpoint":
		caseLegacyEndpoint()
	default:
		panic("No valid test case: " + *testcase)
	}

	sendError()
}

func newDefaultConfig() bugsnag.Configuration {
	return bugsnag.Configuration{
		APIKey: os.Getenv("API_KEY"),
		Endpoints: bugsnag.Endpoints{
			Notify:   os.Getenv("NOTIFY_ENDPOINT"),
			Sessions: os.Getenv("SESSIONS_ENDPOINT"),
		},
	}
}

func sendError() {
	notifier := bugsnag.New()
	notifier.NotifySync(fmt.Errorf("oops"), true)
}

func caseDefault() {
	config := newDefaultConfig()
	bugsnag.Configure(config)
}

func caseAppVersion() {
	config := newDefaultConfig()
	config.AppVersion = os.Getenv("APP_VERSION")
	bugsnag.Configure(config)
}

func caseAppType() {
	config := newDefaultConfig()
	config.AppType = os.Getenv("APP_TYPE")
	bugsnag.Configure(config)
}

func caseLegacyEndpoint() {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:   os.Getenv("API_KEY"),
		Endpoint: os.Getenv("NOTIFY_ENDPOINT"),
	})
}
