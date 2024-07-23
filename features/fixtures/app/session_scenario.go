package main

import (
	"context"
	"fmt"

	"github.com/bugsnag/bugsnag-go/v2"
)

func SendSessionScenario() (bugsnag.Configuration, func()) {
	config := bugsnag.Configuration{}
	scenarioFunc := func() {
		bugsnag.StartSession(context.Background())
	}
	return config, scenarioFunc
}

func SessionAndErrorScenario() (bugsnag.Configuration, func()) {
	config := bugsnag.Configuration{}
	scenarioFunc := func() {
		ctx := bugsnag.StartSession(context.Background())
		bugsnag.Notify(fmt.Errorf("oops"), ctx)
	}
	return config, scenarioFunc
}
