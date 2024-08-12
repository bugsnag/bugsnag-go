Feature: Reporting multiple handled and unhandled errors in the same session

Background:
  Given I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"

Scenario: Handled errors know about previous reported handled errors
  When I start the service "app"
  And I run "MultipleHandledScenario"
  And I wait to receive 2 errors
  And the event handled sessions count equals 1 for request 0
  And the event handled sessions count equals 2 for request 1

Scenario: Unhandled errors know about previous reported handled errors
  When I start the service "app"
  And I run "MultipleUnhandledScenario"
  And I wait to receive 2 errors
  And the event unhandled sessions count equals 1 for request 0
  And the event unhandled sessions count equals 2 for request 1
