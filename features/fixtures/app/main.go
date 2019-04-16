package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	bugsnag "github.com/bugsnag/bugsnag-go"
)

func configureBasicBugsnag(testcase string) {
	config := bugsnag.Configuration{
		APIKey:     os.Getenv("API_KEY"),
		AppVersion: os.Getenv("APP_VERSION"),
		AppType:    os.Getenv("APP_TYPE"),
		Hostname:   os.Getenv("HOSTNAME"),
	}

	if notifyReleaseStages := os.Getenv("NOTIFY_RELEASE_STAGES"); notifyReleaseStages != "" {
		config.NotifyReleaseStages = strings.Split(notifyReleaseStages, ",")
	}

	if releaseStage := os.Getenv("RELEASE_STAGE"); releaseStage != "" {
		config.ReleaseStage = releaseStage
	}

	if filters := os.Getenv("PARAMS_FILTERS"); filters != "" {
		config.ParamsFilters = []string{filters}
	}

	sync, err := strconv.ParseBool(os.Getenv("SYNCHRONOUS"))
	if err == nil {
		config.Synchronous = sync
	}

	acs, err := strconv.ParseBool(os.Getenv("AUTO_CAPTURE_SESSIONS"))
	if err == nil {
		config.AutoCaptureSessions = acs
	}

	switch testcase {
	case "endpoint legacy":
		config.Endpoint = os.Getenv("BUGSNAG_ENDPOINT")
	case "endpoint notify":
		config.Endpoints = bugsnag.Endpoints{Notify: os.Getenv("BUGSNAG_ENDPOINT")}
	case "endpoint session":
		config.Endpoints = bugsnag.Endpoints{Sessions: os.Getenv("BUGSNAG_ENDPOINT")}
	default:
		config.Endpoints = bugsnag.Endpoints{
			Notify:   os.Getenv("BUGSNAG_ENDPOINT"),
			Sessions: os.Getenv("BUGSNAG_ENDPOINT"),
		}
	}
	bugsnag.Configure(config)

	time.Sleep(200 * time.Millisecond)
	// Increase publish rate for testing
	bugsnag.DefaultSessionPublishInterval = time.Millisecond * 100
}

func main() {

	test := flag.String("test", "handled", "what the app should send, either handled, unhandled, session, autonotify")
	flag.Parse()

	configureBasicBugsnag(*test)
	time.Sleep(100 * time.Millisecond) // Ensure tests are less flaky by ensuring the start-up session gets sent

	switch *test {
	case "unhandled":
		unhandledCrash()
	case "handled", "endpoint legacy", "endpoint notify", "endpoint session":
		handledError()
	case "session":
		session()
	case "autonotify":
		autonotify()
	case "metadata":
		metadata()
	case "onbeforenotify":
		onBeforeNotify()
	case "filtered":
		filtered()
	case "recover":
		dontDie()
	case "session and error":
		sessionAndError()
	case "send and exit":
		sendAndExit()
	case "user":
		user()
	case "multiple handled":
		multipleHandled()
	case "multiple unhandled":
		multipleUnhandled()
	default:
		log.Println("Not a valid test flag: " + *test)
		os.Exit(1)
	}

}

func multipleHandled() {
	//Make the order of the below predictable
	bugsnag.Configure(bugsnag.Configuration{Synchronous: true})

	ctx := bugsnag.StartSession(context.Background())
	bugsnag.Notify(fmt.Errorf("oops"), ctx)
	bugsnag.Notify(fmt.Errorf("oops"), ctx)
}

func multipleUnhandled() {
	//Make the order of the below predictable
	notifier := bugsnag.New(bugsnag.Configuration{Synchronous: true})
	notifier.FlushSessionsOnRepanic(false)

	ctx := bugsnag.StartSession(context.Background())
	defer func() { recover() }()
	defer notifier.AutoNotify(ctx)
	defer notifier.AutoNotify(ctx)
	panic("oops")
}

func unhandledCrash() {
	// Invalid type assertion, will panic
	func(a interface{}) string {
		return a.(string)
	}(struct{}{})
}

func handledError() {
	if _, err := os.Open("nonexistent_file.txt"); err != nil {
		if errClass := os.Getenv("ERROR_CLASS"); errClass != "" {
			bugsnag.Notify(err, bugsnag.ErrorClass{Name: errClass})
		} else {
			bugsnag.Notify(err)
		}
	}
	// Give some time for the error to be sent before exiting
	time.Sleep(200 * time.Millisecond)
}

func session() {
	bugsnag.StartSession(context.Background())

	// Give some time for the session to be sent before exiting
	time.Sleep(200 * time.Millisecond)
}

func autonotify() {
	go func() {
		defer bugsnag.AutoNotify()
		panic("Go routine killed with auto notify")
	}()

	// Give enough time for the panic to happen
	time.Sleep(100 * time.Millisecond)
}

func metadata() {
	customerData := map[string]string{"Name": "Joe Bloggs", "Age": "21"}
	bugsnag.Notify(fmt.Errorf("oops"), bugsnag.MetaData{
		"Scheme": {
			"Customer": customerData,
			"Level":    "Blue",
		},
	})
	time.Sleep(200 * time.Millisecond)
}

func filtered() {
	bugsnag.Notify(fmt.Errorf("oops"), bugsnag.MetaData{
		"Account": {
			"Name":           "Company XYZ",
			"Price(dollars)": "1 Million",
		},
	})
	time.Sleep(200 * time.Millisecond)
}

func onBeforeNotify() {
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
	bugsnag.Notify(fmt.Errorf("Ignore this error"))
	time.Sleep(100 * time.Millisecond)
	bugsnag.Notify(fmt.Errorf("Don't ignore this error"))
	time.Sleep(100 * time.Millisecond)
	bugsnag.Notify(fmt.Errorf("Change error message"))
	time.Sleep(100 * time.Millisecond)
}

func dontDie() {
	go func() {
		defer bugsnag.Recover()
		panic("Go routine killed but recovered")
	}()
	time.Sleep(100 * time.Millisecond)
}

func sessionAndError() {
	ctx := bugsnag.StartSession(context.Background())
	bugsnag.Notify(fmt.Errorf("oops"), ctx)

	time.Sleep(200 * time.Millisecond)
}

func sendAndExit() {
	bugsnag.Notify(fmt.Errorf("oops"))
}

func user() {
	bugsnag.Notify(fmt.Errorf("oops"), bugsnag.User{
		Id:    "test-user-id",
		Name:  "test-user-name",
		Email: "test-user-email",
	})

	time.Sleep(200 * time.Millisecond)
}
