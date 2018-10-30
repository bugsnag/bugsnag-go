package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	bugsnag "github.com/bugsnag/bugsnag-go"
)

func main() {
	testcase := flag.String("case", "", "test case to run")
	send := flag.String("send", "", "whether to send a session or error")
	flag.Parse()

	// Increase publish rate for testing
	bugsnag.DefaultSessionPublishInterval = time.Millisecond * 20

	sentError := false

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
	case "on before notify":
		caseOnBeforeNotify()
		sentError = true
	case "params filters":
		caseParamsFilters()
	case "synchronous":
		caseSynchronous()
		sentError = true
	case "user data":
		caseUserData()
		sentError = true
	case "metadata":
		caseMetaData()
		sentError = true
	case "auto notify":
		caseAutoNotify()
		sentError = true
	case "recover":
		caseRecover()
		sentError = true
	case "session":
		caseSession()
		sentError = true
	default:
		panic("No valid test case: " + *testcase)
	}

	if *send == "error" {
		if !sentError {
			sendError()
		}
	} else if *send == "session" {
		bugsnag.StartSession(context.Background())

		// Give some time for the session to be sent before exiting
		time.Sleep(100 * time.Millisecond)
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
	notifier.NotifySync(fmt.Errorf("oops"), true, bugsnag.MetaData{
		"Account": {
			"Name":           "Company XYZ",
			"Price(dollars)": "1 Million",
		},
	})
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

func caseOnBeforeNotify() {
	config := newDefaultConfig()
	bugsnag.Configure(config)
	bugsnag.OnBeforeNotify(
		func(event *bugsnag.Event, config *bugsnag.Configuration) error {
			if event.Message == "Ignore this error" {
				return fmt.Errorf("not sending errors to ignore")
			}
			// continue notifying as normal
			if event.Message == "Change error message" {
				event.Message = "Error message was changed"
			}
			return nil
		})

	notifier := bugsnag.New()
	notifier.NotifySync(fmt.Errorf("Don't ignore this error"), true)
	notifier.NotifySync(fmt.Errorf("Ignore this error"), true)
	notifier.NotifySync(fmt.Errorf("Change error message"), true)
}

func caseParamsFilters() {
	config := newDefaultConfig()
	paramsFilters := os.Getenv("PARAMS_FILTERS")
	if paramsFilters != "" {
		config.ParamsFilters = strings.Split(paramsFilters, ",")
	}
	bugsnag.Configure(config)
}

func caseSynchronous() {
	config := newDefaultConfig()
	sync, err := strconv.ParseBool(os.Getenv("SYNCHRONOUS"))
	if err != nil {
		panic("Unknown synchronous flag: " + err.Error())
	}
	config.Synchronous = sync
	bugsnag.Configure(config)

	notifier := bugsnag.New()
	notifier.Notify(fmt.Errorf("Generic error"))
}

func caseUserData() {
	config := newDefaultConfig()
	bugsnag.Configure(config)

	notifier := bugsnag.New()
	notifier.NotifySync(fmt.Errorf("oops"), true, bugsnag.User{
		Id:    os.Getenv("USER_ID"),
		Name:  os.Getenv("USER_NAME"),
		Email: os.Getenv("USER_EMAIL"),
	})
}

func caseMetaData() {
	config := newDefaultConfig()
	bugsnag.Configure(config)

	notifier := bugsnag.New()

	customerData := map[string]string{"Name": "Joe Bloggs", "Age": "21"}
	notifier.NotifySync(fmt.Errorf("oops"), true, bugsnag.MetaData{
		"Scheme": {
			"Customer": customerData,
			"Level":    "Blue",
		},
	})
}

func caseAutoNotify() {
	config := newDefaultConfig()
	bugsnag.Configure(config)

	go func() {
		defer bugsnag.AutoNotify()
		panic("Go routine killed")
	}()

	time.Sleep(200 * time.Millisecond)
}

func caseRecover() {
	config := newDefaultConfig()
	bugsnag.Configure(config)

	go func() {
		defer bugsnag.Recover()
		panic("Go routine killed")
	}()

	time.Sleep(200 * time.Millisecond)
}

func caseSession() {
	config := newDefaultConfig()
	bugsnag.Configure(config)

	ctx := bugsnag.StartSession(context.Background())
	notifier := bugsnag.New()
	notifier.NotifySync(fmt.Errorf("oops"), true, ctx)

	time.Sleep(200 * time.Millisecond)
}
