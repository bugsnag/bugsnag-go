package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/martini"
	"github.com/go-martini/martini"
)

func main() {
	m := martini.Classic()

	if os.Getenv("DISABLE_REPORT_PAYLOADS") != "" {
		// Increase publish rate for testing
		bugsnag.DefaultSessionPublishInterval = time.Millisecond * 20
	}

	config := bugsnag.Configuration{
		AppVersion: os.Getenv("APP_VERSION"),
		AppType:    os.Getenv("APP_TYPE"),
		APIKey:     os.Getenv("API_KEY"),
		Endpoint:   os.Getenv("ENDPOINT"),
		Endpoints: bugsnag.Endpoints{
			Notify:   os.Getenv("NOTIFY_ENDPOINT"),
			Sessions: os.Getenv("SESSIONS_ENDPOINT"),
		},
		Hostname:     os.Getenv("HOSTNAME"),
		ReleaseStage: os.Getenv("RELEASE_STAGE"),
	}

	if stages := os.Getenv("NOTIFY_RELEASE_STAGES"); stages != "" {
		config.NotifyReleaseStages = []string{stages}
	}

	if acs, _ := strconv.ParseBool(os.Getenv("AUTO_CAPTURE_SESSIONS")); acs {
		config.AutoCaptureSessions = acs
	}

	if filters := os.Getenv("PARAMS_FILTERS"); filters != "" {
		config.ParamsFilters = []string{filters}
	}

	config.Synchronous, _ = strconv.ParseBool(os.Getenv("SYNCHRONOUS"))
	bugsnag.Configure(config)

	m.Use(martini.Recovery())
	m.Use(bugsnagmartini.AutoNotify())

	m.Get("/unhandled", performUnhandledCrash)
	m.Get("/handled", performHandledError)
	m.Get("/metadata", metadata)

	m.RunOnAddr(":9030")
}

func performUnhandledCrash() {
	// Invalid type assertion, will panic
	func(a interface{}) string { return a.(string) }(struct{}{})
}

func performHandledError(r *http.Request) {
	if _, err := os.Open("nonexistent_file.txt"); err != nil {
		bugsnag.Notify(err, r.Context())
	}
}

func metadata() {
	customerData := map[string]string{"Name": "Joe Bloggs", "Age": "21"}
	bugsnag.Notify(fmt.Errorf("oops"), true, bugsnag.MetaData{
		"Scheme": {
			"Customer": customerData,
			"Level":    "Blue",
		},
	})
}
