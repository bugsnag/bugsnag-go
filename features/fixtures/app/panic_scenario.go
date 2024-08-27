package main

import (
	"github.com/bugsnag/bugsnag-go/v2"
)

func AutonotifyPanicScenario(command Command) func() {
	scenarioFunc := func() {
		defer bugsnag.AutoNotify()
		panic("Go routine killed with auto notify")
	}

	return scenarioFunc
}

func RecoverAfterPanicScenario(command Command) func() {
	scenarioFunc := func() {
		defer bugsnag.Recover()
		panic("Go routine killed but recovered")
	}
	return scenarioFunc
}