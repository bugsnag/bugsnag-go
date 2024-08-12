package main

import (
	"github.com/bugsnag/bugsnag-go/v2"
)

func AutonotifyPanicScenario(command Command) (bugsnag.Configuration, func()) {
	config := ConfigureBugsnag()
	config.APIKey = command.APIKey
	config.Endpoints.Sessions = command.SessionsEndpoint
	config.Endpoints.Notify = command.NotifyEndpoint

	scenarioFunc := func() {
		defer bugsnag.AutoNotify()
		panic("Go routine killed with auto notify")
	}

	return config, scenarioFunc
}

func RecoverAfterPanicScenario(command Command) (bugsnag.Configuration, func()) {
	config := ConfigureBugsnag()
	config.APIKey = command.APIKey
	config.Endpoints.Sessions = command.SessionsEndpoint
	config.Endpoints.Notify = command.NotifyEndpoint

	scenarioFunc := func() {
		defer bugsnag.Recover()
		panic("Go routine killed but recovered")
	}
	return config, scenarioFunc
}
