package main

import (
	"fmt"

	bugsnag "github.com/bugsnag/bugsnag-go/v2"
)

func AutoconfigPanicScenario(command Command) (bugsnag.Configuration, func()) {
	config := ConfigureBugsnag(command)

	scenarioFunc := func() {
		panic("PANIQ!")
	}
	return config, scenarioFunc
}

func AutoconfigHandledScenario(command Command) (bugsnag.Configuration, func()) {
	config := ConfigureBugsnag(command)

	scenarioFunc := func() {
		bugsnag.Notify(fmt.Errorf("gone awry!"))
	}
	return config, scenarioFunc
}

func AutoconfigMetadataScenario(command Command) (bugsnag.Configuration, func()) {
	config := ConfigureBugsnag(command)

	scenarioFunc := func() {
		bugsnag.OnBeforeNotify(func(event *bugsnag.Event, config *bugsnag.Configuration) error {
			event.MetaData.Add("fruit", "Tomato", "beefsteak")
			event.MetaData.Add("snacks", "Carrot", "4")
			return nil
		})
		bugsnag.Notify(fmt.Errorf("gone awry!"))
	}
	return config, scenarioFunc
}