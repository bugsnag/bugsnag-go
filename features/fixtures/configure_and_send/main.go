package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	bugsnag "github.com/bugsnag/bugsnag-go"
)

func main() {
	testcase := flag.String("case", "", "test case to run")
	send := flag.String("send", "", "whether to send a session/error or both")
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
	case "hostname":
		caseHostname()
	case "release stage":
		caseNotifyReleaseStage()
	default:
		panic("No valid test case: " + *testcase)
	}

	if *send == "error" {
		sendError()
	} else if *send == "session" {
		sendSession()
	} else {
		panic("No valid send case: " + *send)
	}
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

func sendSession() {
	bugsnag.DefaultSessionPublishInterval = time.Millisecond * 10
	bugsnag.StartSession(context.Background())
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

func caseHostname() {
	config := newDefaultConfig()
	config.Hostname = os.Getenv("HOSTNAME")
	bugsnag.Configure(config)
}

func caseNotifyReleaseStage() {
	config := newDefaultConfig()
	notifyReleaseStages := os.Getenv("NOTIFY_RELEASE_STAGES")
	if notifyReleaseStages != "" {
		config.NotifyReleaseStages = strings.Split(notifyReleaseStages, ",")
	}
	releaseStage := os.Getenv("RELEASE_STAGE")
	if releaseStage != "" {
		config.ReleaseStage = releaseStage
	}
	bugsnag.Configure(config)
}
