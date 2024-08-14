package main

import (
	"fmt"

	bugsnag "github.com/bugsnag/bugsnag-go/v2"
)

func AutoconfigPanicScenario(command Command) func() {
	scenarioFunc := func() {
		panic("PANIQ!")
	}
	return scenarioFunc
}

func AutoconfigHandledScenario(command Command) func() {
	scenarioFunc := func() {
		bugsnag.Notify(fmt.Errorf("gone awry!"))
	}
	return scenarioFunc
}

func AutoconfigMetadataScenario(command Command) func() {
	scenarioFunc := func() {
		bugsnag.OnBeforeNotify(func(event *bugsnag.Event, config *bugsnag.Configuration) error {
			event.MetaData.Add("fruit", "Tomato", "beefsteak")
			event.MetaData.Add("snacks", "Carrot", "4")
			return nil
		})
		bugsnag.Notify(fmt.Errorf("gone awry!"))
	}
	return scenarioFunc
}
