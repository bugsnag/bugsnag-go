package main

import (
	"fmt"

	"github.com/bugsnag/bugsnag-go/v2"
)

func MetadataScenario(command Command) func() {
	scenarioFunc := func() {
		customerData := map[string]string{"Name": "Joe Bloggs", "Age": "21"}
		bugsnag.Notify(fmt.Errorf("oops"), bugsnag.MetaData{
			"Scheme": {
				"Customer": customerData,
				"Level":    "Blue",
			},
		})
	}
	return scenarioFunc
}

func FilteredMetadataScenario(command Command) func() {
	scenarioFunc := func() {
		bugsnag.Notify(fmt.Errorf("oops"), bugsnag.MetaData{
			"Account": {
				"Name":           "Company XYZ",
				"Price(dollars)": "1 Million",
			},
		})
	}
	return scenarioFunc
}
