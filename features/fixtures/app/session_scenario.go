package main

import (
	"context"
	"fmt"

	"github.com/bugsnag/bugsnag-go/v2"
)

func SendSessionScenario(command Command) (bugsnag.Configuration, func()) {
	config := ConfigureBugsnag(command)

	scenarioFunc := func() {
		bugsnag.StartSession(context.Background())
	}
	return config, scenarioFunc
}

func SessionAndErrorScenario(command Command) (bugsnag.Configuration, func()) {
	config := ConfigureBugsnag(command)
	scenarioFunc := func() {
		ctx := bugsnag.StartSession(context.Background())
		bugsnag.Notify(fmt.Errorf("oops"), ctx)
	}
	return config, scenarioFunc
}
