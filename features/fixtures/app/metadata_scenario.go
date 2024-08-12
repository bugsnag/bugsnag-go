package main

import (
	"fmt"

	"github.com/bugsnag/bugsnag-go/v2"
)

func MetadataScenario(command Command) (bugsnag.Configuration, func()) {
	config := ConfigureBugsnag(command)

	scenarioFunc := func() {
		customerData := map[string]string{"Name": "Joe Bloggs", "Age": "21"}
		bugsnag.Notify(fmt.Errorf("oops"), bugsnag.MetaData{
			"Scheme": {
				"Customer": customerData,
				"Level":    "Blue",
			},
		})
	}
	return config, scenarioFunc
}

func FilteredMetadataScenario(command Command) (bugsnag.Configuration, func()) {
	config := ConfigureBugsnag(command)

	scenarioFunc := func() {
		bugsnag.Notify(fmt.Errorf("oops"), bugsnag.MetaData{
			"Account": {
				"Name":           "Company XYZ",
				"Price(dollars)": "1 Million",
			},
		})
	}
	return config, scenarioFunc
}
