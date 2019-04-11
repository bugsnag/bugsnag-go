package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/negroni"
	"github.com/urfave/negroni"
)

func main() {
	config := bugsnag.Configuration{
		APIKey: os.Getenv("API_KEY"),
		Endpoints: bugsnag.Endpoints{
			Notify:   os.Getenv("BUGSNAG_ENDPOINT"),
			Sessions: os.Getenv("BUGSNAG_ENDPOINT"),
		},
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

	acs, err := strconv.ParseBool(os.Getenv("AUTO_CAPTURE_SESSIONS"))
	if err == nil {
		config.AutoCaptureSessions = acs
	}
	bugsnag.Configure(config)

	// Increase publish rate for testing
	bugsnag.DefaultSessionPublishInterval = time.Millisecond * 300

	mux := http.NewServeMux()
	mux.HandleFunc("/autonotify-then-recover", unhandledCrash)
	mux.HandleFunc("/handled", handledError)
	mux.HandleFunc("/session", session)
	mux.HandleFunc("/autonotify", autonotify)
	mux.HandleFunc("/onbeforenotify", onBeforeNotify)
	mux.HandleFunc("/recover", dontDie)
	mux.HandleFunc("/user", user)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	// Add bugsnag handler after negroni.NewRecovery() to ensure panics get picked up
	n.Use(bugsnagnegroni.AutoNotify())
	n.UseHandler(mux)
	http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), n)
}

func unhandledCrash(w http.ResponseWriter, r *http.Request) {
	// Invalid type assertion, will panic
	func(a interface{}) string {
		return a.(string)
	}(struct{}{})
}

func handledError(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Open("nonexistent_file.txt"); err != nil {
		if errClass := os.Getenv("ERROR_CLASS"); errClass != "" {
			bugsnag.Notify(err, r.Context(), bugsnag.ErrorClass{Name: errClass})
		} else {
			bugsnag.Notify(err, r.Context())
		}
	}
}

func session(w http.ResponseWriter, r *http.Request) {
	log.Println("single session")
}

func dontDie(w http.ResponseWriter, r *http.Request) {
	defer bugsnag.Recover(r.Context())
	panic("Request killed but recovered")
}

func user(w http.ResponseWriter, r *http.Request) {
	bugsnag.Notify(fmt.Errorf("oops"), bugsnag.User{
		Id:    "test-user-id",
		Name:  "test-user-name",
		Email: "test-user-email",
	})
}

func onBeforeNotify(w http.ResponseWriter, r *http.Request) {
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

func autonotify(w http.ResponseWriter, r *http.Request) {
	go func(ctx context.Context) {
		defer func() { recover() }()
		defer bugsnag.AutoNotify(ctx)
		panic("Go routine killed with auto notify")
	}(r.Context())
}
