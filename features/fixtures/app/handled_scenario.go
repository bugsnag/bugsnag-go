package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/bugsnag/bugsnag-go/v2"
)

func HandledErrorScenario() (bugsnag.Configuration, func()) {
	config := bugsnag.Configuration{}
	scenarioFunc := func() {
		if _, err := os.Open("nonexistent_file.txt"); err != nil {
			if errClass := os.Getenv("ERROR_CLASS"); errClass != "" {
				bugsnag.Notify(err, bugsnag.ErrorClass{Name: errClass})
			} else {
				bugsnag.Notify(err)
			}
		}
	}
	return config, scenarioFunc
}

func MultipleHandledErrorsScenario() (bugsnag.Configuration, func()) {
	//Make the order of the below predictable
	config := bugsnag.Configuration{Synchronous: true}

	scenarioFunc := func() {
		ctx := bugsnag.StartSession(context.Background())
		bugsnag.Notify(fmt.Errorf("oops"), ctx)
		bugsnag.Notify(fmt.Errorf("oops"), ctx)
	}
	return config, scenarioFunc
}

func NestedHandledErrorScenario() (bugsnag.Configuration, func()) {
	config := bugsnag.Configuration{}
	scenarioFunc := func() {
		if err := Login("token " + os.Getenv("API_KEY")); err != nil {
			bugsnag.Notify(NewCustomErr("terminate process", err))
		} else {
			i := len(os.Getenv("API_KEY"))
			// Some nonsense to avoid inlining checkValue
			if val, err := CheckValue(i); err != nil {
				fmt.Printf("err: %v, val: %d", err, val)
			}
			if val, err := CheckValue(i - 46); err != nil {
				fmt.Printf("err: %v, val: %d", err, val)
			}

			log.Fatalf("This test is broken - no error was generated.")
		}
	}
	return config, scenarioFunc
}

func HandledCallbackErrorScenario() (bugsnag.Configuration, func()) {
	config := bugsnag.Configuration{}
	scenarioFunc := func() {
		bugsnag.Notify(fmt.Errorf("inadequent Prep Error"), func(event *bugsnag.Event) {
			event.Context = "nonfatal.go:14"
			event.Severity = bugsnag.SeverityInfo

			event.Stacktrace[1].File = ">insertion<"
			event.Stacktrace[1].LineNumber = 0
		})
	}
	return config, scenarioFunc
}

func HandledToUnhandledScenario() (bugsnag.Configuration, func()) {
	config := bugsnag.Configuration{}
	scenarioFunc := func() {
		bugsnag.Notify(fmt.Errorf("unknown event"), func(event *bugsnag.Event) {
			event.Unhandled = true
			event.Severity = bugsnag.SeverityError
		})
	}
	return config, scenarioFunc
}

func OnBeforeNotifyScenario() (bugsnag.Configuration, func()) {
	config := bugsnag.Configuration{}
	scenarioFunc := func() {
		bugsnag.OnBeforeNotify(
			func(event *bugsnag.Event, config *bugsnag.Configuration) error {
				if event.Message == "ignore this error" {
					return fmt.Errorf("not sending errors to ignore")
				}
				// continue notifying as normal
				if event.Message == "change error message" {
					event.Message = "Error message was changed"
				}
				return nil
			})
		bugsnag.Notify(fmt.Errorf("ignore this error"))
		bugsnag.Notify(fmt.Errorf("don't ignore this error"))
		bugsnag.Notify(fmt.Errorf("change error message"))
	}
	return config, scenarioFunc
}
