package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"time"

	bugsnag "github.com/bugsnag/bugsnag-go"
)

func main() {
	testcase := flag.String("case", "", "test case to run")
	flag.Parse()

	// Increase publish rate for testing
	bugsnag.DefaultSessionPublishInterval = time.Millisecond * 20

	switch *testcase {
	case "default":
		caseDefault()
	case "auto notify":
		caseAutoNotify()
	case "meta data":
		caseMetaData()
	case "on before notify":
		caseOnBeforeNotify()
	case "recover":
		caseRecover()
	case "user data":
		caseUserData()
	case "auto capture sessions":
		caseAutoCaptureSessions()

	default:
		panic("No valid test case: " + *testcase)
	}
}

func startTestServer() *httptest.Server {
	ts := httptest.NewServer(bugsnag.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bugsnag.Notify(fmt.Errorf("oops"), r.Context())
	})))
	return ts
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

func caseDefault() {
	config := newDefaultConfig()
	bugsnag.Configure(config)

	ts := startTestServer()
	defer ts.Close()

	http.Get(ts.URL + "/1234abcd?fish=bird")
	time.Sleep(200 * time.Millisecond)
}

func caseAutoNotify() {
	config := newDefaultConfig()
	bugsnag.Configure(config)

	ts := httptest.NewServer(bugsnag.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		go func(ctx context.Context) {
			defer bugsnag.AutoNotify(ctx)
			panic("Go routine killed")
		}(r.Context())
	})))
	defer ts.Close()

	http.Get(ts.URL + "/1234abcd?fish=bird")
	time.Sleep(200 * time.Millisecond)
}

func caseMetaData() {
	config := newDefaultConfig()
	bugsnag.Configure(config)

	ts := httptest.NewServer(bugsnag.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customerData := map[string]string{"Name": "Joe Bloggs", "Age": "21"}
		bugsnag.Notify(fmt.Errorf("oops"), r.Context(), bugsnag.MetaData{
			"Scheme": {
				"Customer": customerData,
				"Level":    "Blue",
			},
		})
	})))
	defer ts.Close()

	http.Get(ts.URL + "/1234abcd?fish=bird")
	time.Sleep(200 * time.Millisecond)
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

	requestCount := 0
	notifier := bugsnag.New()
	ts := httptest.NewServer(bugsnag.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requestCount == 0 {
			notifier.NotifySync(fmt.Errorf("Don't ignore this error"), true)
		} else if requestCount == 1 {
			notifier.NotifySync(fmt.Errorf("Ignore this error"), true)
		} else if requestCount == 2 {
			notifier.NotifySync(fmt.Errorf("Change error message"), true)
		}
		requestCount++
	})))
	defer ts.Close()

	http.Get(ts.URL + "/1234abcd?fish=bird")
	http.Get(ts.URL + "/1234abcd?fish=bird")
	http.Get(ts.URL + "/1234abcd?fish=bird")
	time.Sleep(200 * time.Millisecond)
}

func caseRecover() {
	config := newDefaultConfig()
	bugsnag.Configure(config)

	ts := httptest.NewServer(bugsnag.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		go func(ctx context.Context) {
			defer bugsnag.Recover(ctx)
			panic("Go routine killed")
		}(r.Context())
	})))
	defer ts.Close()

	http.Get(ts.URL + "/1234abcd?fish=bird")
	time.Sleep(200 * time.Millisecond)
}

func caseUserData() {
	config := newDefaultConfig()
	bugsnag.Configure(config)

	ts := httptest.NewServer(bugsnag.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bugsnag.Notify(fmt.Errorf("oops"), bugsnag.User{
			Id:    os.Getenv("USER_ID"),
			Name:  os.Getenv("USER_NAME"),
			Email: os.Getenv("USER_EMAIL"),
		})
	})))
	defer ts.Close()

	http.Get(ts.URL + "/1234abcd?fish=bird")
	time.Sleep(200 * time.Millisecond)
}

func caseAutoCaptureSessions() {
	config := newDefaultConfig()
	acsFlag, err := strconv.ParseBool(os.Getenv("AUTO_CAPTURE_SESSIONS"))
	if err != nil {
		panic(err)
	}
	config.AutoCaptureSessions = acsFlag
	bugsnag.Configure(config)

	ts := startTestServer()
	defer ts.Close()

	http.Get(ts.URL + "/1234abcd?fish=bird")
	time.Sleep(200 * time.Millisecond)
}
