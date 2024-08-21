Feature: Reporting multiple handled and unhandled errors in the same session

Background:
  Given I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"

Scenario: Handled errors know about previous reported handled errors
  When I start the service "app"
  And I run "MultipleHandledErrorsScenario"
  And I wait to receive 2 errors
  And the event handled sessions count equals 1
  And I discard the oldest error
  And the event handled sessions count equals 2

Scenario: Unhandled errors know about previous reported handled errors
  When I start the service "app"
  And I run "MultipleUnhandledErrorsScenario"
  And I wait to receive 2 errors
  And the event unhandled sessions count equals 1
  And I discard the oldest error
  And the event unhandled sessions count equals 2