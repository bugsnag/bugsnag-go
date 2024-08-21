package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bugsnag/bugsnag-go/v2"
)

func HttpServerScenario(Command) func() {
	scenarioFunc := func() {
		http.HandleFunc("/handled", handledError)
		http.HandleFunc("/autonotify-then-recover", unhandledCrash)
		http.HandleFunc("/session", session)
		http.HandleFunc("/autonotify", autonotify)
		http.HandleFunc("/onbeforenotify", onBeforeNotify)
		http.HandleFunc("/recover", dontdie)
		http.HandleFunc("/user", user)

		http.ListenAndServe(":4512", recoverWrap(bugsnag.Handler(nil)))
	}

	return scenarioFunc
}

// Simple wrapper to send internal server error on panics
func recoverWrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r != nil {
				http.Error(w, "", http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
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

func unhandledCrash(w http.ResponseWriter, r *http.Request) {
	// Invalid type assertion, will panic
	func(a interface{}) string {
		return a.(string)
	}(struct{}{})
}

func session(w http.ResponseWriter, r *http.Request) {
	log.Println("single session")
}

func autonotify(w http.ResponseWriter, r *http.Request) {
	go func(ctx context.Context) {
		defer func() { recover() }()
		defer bugsnag.AutoNotify(ctx)
		panic("Go routine killed with auto notify")
	}(r.Context())
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

func dontdie(w http.ResponseWriter, r *http.Request) {
	defer bugsnag.Recover(r.Context())
	panic("Request killed but recovered")
}

func user(w http.ResponseWriter, r *http.Request) {
	bugsnag.Notify(fmt.Errorf("oops"), r.Context(), bugsnag.User{
		Id:    "test-user-id",
		Name:  "test-user-name",
		Email: "test-user-email",
	})
}