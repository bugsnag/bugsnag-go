package main

import (
	"context"
	"fmt"

	"github.com/bugsnag/bugsnag-go/v2"
)

func SendSessionScenario(command Command) func() {
	scenarioFunc := func() {
		bugsnag.StartSession(context.Background())
	}
	return scenarioFunc
}

func SessionAndErrorScenario(command Command) func() {
	scenarioFunc := func() {
		ctx := bugsnag.StartSession(context.Background())
		bugsnag.Notify(fmt.Errorf("oops"), ctx)
	}
	return scenarioFunc
}
