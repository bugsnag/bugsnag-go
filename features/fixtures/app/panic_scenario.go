package main

import (
	"github.com/bugsnag/bugsnag-go/v2"
)

func AutonotifyPanicScenario()(bugsnag.Configuration, func())  {
	config := bugsnag.Configuration{}
	scenarioFunc := func() {
		defer bugsnag.AutoNotify()
		panic("Go routine killed with auto notify")
	}

	return config, scenarioFunc
}

func RecoverAfterPanicScenario() (bugsnag.Configuration, func()) {
	config := bugsnag.Configuration{}
	scenarioFunc := func() {
		defer bugsnag.Recover()
		panic("Go routine killed but recovered")
	}
	return config, scenarioFunc
}