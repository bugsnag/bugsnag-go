package main

import (
	"context"
	"fmt"

	"github.com/bugsnag/bugsnag-go/v2"
)

//go:noinline
func UnhandledCrashScenario(command Command) (bugsnag.Configuration, func()) {
	config := ConfigureBugsnag(command)

	scenarioFunc := func() {
		fmt.Printf("Calling panic\n")
		// Invalid type assertion, will panic
		func(a interface{}) string {
			return a.(string)
		}(struct{}{})
	}
	return config, scenarioFunc
}

func MultipleUnhandledErrorsScenario(command Command) (bugsnag.Configuration, func()) {
	config := ConfigureBugsnag(command)

	scenarioFunc := func() {
		//Make the order of the below predictable
		notifier := bugsnag.New(bugsnag.Configuration{Synchronous: true})
		notifier.FlushSessionsOnRepanic(false)

		ctx := bugsnag.StartSession(context.Background())
		defer func() { recover() }()
		defer notifier.AutoNotify(ctx)
		defer notifier.AutoNotify(ctx)
		panic("oops")
	}
	return config, scenarioFunc
}
