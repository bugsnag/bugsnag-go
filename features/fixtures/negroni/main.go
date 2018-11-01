package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/negroni"
	"github.com/urfave/negroni"
)

func main() {
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

	mux := http.NewServeMux()
	mux.HandleFunc("/unhandled", unhandledCrash)
	mux.HandleFunc("/handled", handledError)
	mux.HandleFunc("/metadata", metadata)
	mux.HandleFunc("/onbeforenotify", onbeforenotify)
	mux.HandleFunc("/recover", dontdie)
	mux.HandleFunc("/async", async)
	mux.HandleFunc("/user", user)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	// Add bugsnag handler after negroni.NewRecovery() to ensure panics get picked up
	n.Use(bugsnagnegroni.AutoNotify())
	n.UseHandler(mux)

	http.ListenAndServe(":9040", n)
}

func unhandledCrash(w http.ResponseWriter, r *http.Request) {
	// Invalid type assertion, will panic
	func(a interface{}) string { return a.(string) }(struct{}{})
}

func handledError(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Open("nonexistent_file.txt"); err != nil {
		bugsnag.Notify(err, r.Context())
	}
}

func metadata(w http.ResponseWriter, r *http.Request) {
	customerData := map[string]string{"Name": "Joe Bloggs", "Age": "21"}
	bugsnag.Notify(fmt.Errorf("oops"), bugsnag.MetaData{
		"Scheme": {
			"Customer": customerData,
			"Level":    "Blue",
		},
	})
}

func dontdie(w http.ResponseWriter, r *http.Request) {
	defer bugsnag.Recover()
	func(a interface{}) string { return a.(string) }(struct{}{})
}

func async(w http.ResponseWriter, r *http.Request) {
	bugsnag.Notify(fmt.Errorf("If I show up it means I was sent synchronously"))
	defer os.Exit(0)
}

func user(w http.ResponseWriter, r *http.Request) {
	bugsnag.Notify(fmt.Errorf("oops"), bugsnag.User{
		Id:    "test-user-id",
		Name:  "test-user-name",
		Email: "test-user-email",
	})
}

func onbeforenotify(w http.ResponseWriter, r *http.Request) {
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
	bugsnag.Notify(fmt.Errorf("Don't ignore this error"))
	bugsnag.Notify(fmt.Errorf("Change error message"))
}
