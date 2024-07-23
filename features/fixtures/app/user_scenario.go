package main

import (
	"fmt"

	"github.com/bugsnag/bugsnag-go/v2"
)

func SetUserScenario() (bugsnag.Configuration, func()) {
	config := bugsnag.Configuration{}
	scenarioFunc := func() {
		bugsnag.Notify(fmt.Errorf("oops"), bugsnag.User{
			Id:    "test-user-id",
			Name:  "test-user-name",
			Email: "test-user-email",
		})
	}

	return config, scenarioFunc
}
