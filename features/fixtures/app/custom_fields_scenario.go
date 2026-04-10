package main

import (
	"os"

	"github.com/bugsnag/bugsnag-go/v2"
)

// CustomErrorClassScenario tests notifying an error with a custom ErrorClass override
func CustomErrorClassScenario(command Command) func() {
	scenarioFunc := func() {
		if _, err := os.Open("nonexistent_file.txt"); err != nil {
			bugsnag.Notify(
				err,
				bugsnag.ErrorClass{Name: "CustomErrorClass"},
			)
		}
	}
	return scenarioFunc
}

// CustomMessageScenario tests notifying an error with a custom Message override
func CustomMessageScenario(command Command) func() {
	scenarioFunc := func() {
		if _, err := os.Open("nonexistent_file.txt"); err != nil {
			bugsnag.Notify(
				err,
				bugsnag.Message{String: "This is a custom message"},
			)
		}
	}
	return scenarioFunc
}

// CustomBothScenario tests notifying an error with both custom ErrorClass and Message
func CustomBothScenario(command Command) func() {
	scenarioFunc := func() {
		if _, err := os.Open("nonexistent_file.txt"); err != nil {
			bugsnag.Notify(
				err,
				bugsnag.ErrorClass{Name: "CustomBothError"},
				bugsnag.Message{String: "Custom message with both fields"},
			)
		}
	}
	return scenarioFunc
}

// MultipleCustomFieldsScenario tests notifying multiple errors with different custom field combinations
func MultipleCustomFieldsScenario(command Command) func() {
	scenarioFunc := func() {
		// Configure to send errors synchronously so each Notify() creates a separate request
		bugsnag.Configure(bugsnag.Configuration{
			Synchronous: true,
		})

		// First error with custom ErrorClass and Message
		if _, err := os.Open("nonexistent_file_1.txt"); err != nil {
			bugsnag.Notify(
				err,
				bugsnag.ErrorClass{Name: "FirstCustomError"},
				bugsnag.Message{String: "First error custom message"},
			)
		}

		// Second error with different custom ErrorClass and Message
		if _, err := os.Open("nonexistent_file_2.txt"); err != nil {
			bugsnag.Notify(
				err,
				bugsnag.ErrorClass{Name: "SecondCustomError"},
				bugsnag.Message{String: "Second error custom message"},
			)
		}
	}
	return scenarioFunc
}
// MultipleErrorClassScenario tests notifying an error with multiple ErrorClass objects
// Expected behavior: the last ErrorClass object should be used
func MultipleErrorClassScenario(command Command) func() {
	scenarioFunc := func() {
		if _, err := os.Open("nonexistent_file.txt"); err != nil {
			bugsnag.Notify(
				err,
				bugsnag.ErrorClass{Name: "FirstErrorClass"},
				bugsnag.ErrorClass{Name: "SecondErrorClass"},
				bugsnag.ErrorClass{Name: "FinalErrorClass"},
			)
		}
	}
	return scenarioFunc
}

// MultipleMessageScenario tests notifying an error with multiple Message objects
// Expected behavior: the last Message object should be used
func MultipleMessageScenario(command Command) func() {
	scenarioFunc := func() {
		if _, err := os.Open("nonexistent_file.txt"); err != nil {
			bugsnag.Notify(
				err,
				bugsnag.Message{String: "First custom message"},
				bugsnag.Message{String: "Second custom message"},
				bugsnag.Message{String: "Final custom message"},
			)
		}
	}
	return scenarioFunc
}

// MultipleOfBothScenario tests notifying an error with multiple ErrorClass and Message objects
// Expected behavior: the last of each type should be used
func MultipleOfBothScenario(command Command) func() {
	scenarioFunc := func() {
		if _, err := os.Open("nonexistent_file.txt"); err != nil {
			bugsnag.Notify(
				err,
				bugsnag.ErrorClass{Name: "FirstErrorClass"},
				bugsnag.Message{String: "First message"},
				bugsnag.ErrorClass{Name: "SecondErrorClass"},
				bugsnag.Message{String: "Second message"},
			)
		}
	}
	return scenarioFunc
}