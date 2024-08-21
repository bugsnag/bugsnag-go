package main

import (
	"fmt"

	"github.com/bugsnag/bugsnag-go/v2"
)

func SetUserScenario(command Command) func() {
	scenarioFunc := func() {
		bugsnag.Notify(fmt.Errorf("oops"), bugsnag.User{
			Id:    "test-user-id",
			Name:  "test-user-name",
			Email: "test-user-email",
		})
	}

	return scenarioFunc
}