package main

import (
	"fmt"

	"github.com/bugsnag/bugsnag-go/v2"
)

func SetUserScenario(command Command) (bugsnag.Configuration, func()) {
	config := ConfigureBugsnag()
	config.APIKey = command.APIKey
	config.Endpoints.Sessions = command.SessionsEndpoint
	config.Endpoints.Notify = command.NotifyEndpoint

	scenarioFunc := func() {
		bugsnag.Notify(fmt.Errorf("oops"), bugsnag.User{
			Id:    "test-user-id",
			Name:  "test-user-name",
			Email: "test-user-email",
		})
	}

	return config, scenarioFunc
}
