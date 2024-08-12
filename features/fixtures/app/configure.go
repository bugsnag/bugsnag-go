package main

import (
	bugsnag "github.com/bugsnag/bugsnag-go/v2"
)

func ConfigureBugsnag(command Command) bugsnag.Configuration {
	config := bugsnag.Configuration{}

	config.APIKey = command.APIKey
	config.Endpoints.Sessions = command.SessionsEndpoint
	config.Endpoints.Notify = command.NotifyEndpoint

	return config
}
