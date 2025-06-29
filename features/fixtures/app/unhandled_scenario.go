package main

import (
	"context"
	"fmt"

	"github.com/bugsnag/bugsnag-go/v2"
)

//go:noinline
func UnhandledCrashScenario(command Command) func() {
	scenarioFunc := func() {
		fmt.Printf("Calling panic\n")
		// Invalid type assertion, will panic
		func(a interface{}) string {
			return a.(string)
		}(struct{}{})
	}
	return scenarioFunc
}

func MultipleUnhandledErrorsScenario(command Command) func() {
	scenarioFunc := func() {
		//Make the order of the below predictable
		notifier := bugsnag.New(bugsnag.Configuration{
			Synchronous: true,
			Endpoints: bugsnag.Endpoints{
				Notify:   command.NotifyEndpoint,
				Sessions: command.SessionsEndpoint,
			},
		})
		notifier.FlushSessionsOnRepanic(false)

		ctx := bugsnag.StartSession(context.Background())
		defer func() { recover() }()
		defer notifier.AutoNotify(ctx)
		defer notifier.AutoNotify(ctx)
		panic("oops")
	}
	return scenarioFunc
}
