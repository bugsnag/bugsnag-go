Feature: Handled errors

Background:
  Given I set environment variable "BUGSNAG_SOURCE_ROOT" to "/app/src/features/fixtures/app/"

Scenario: A handled error sends a report
  When I start the service "app"
  And I run "HttpServerScenario"
  And I wait for the host "localhost" to open port "4512"
  And I open the URL "http://localhost:4512/handled"
  Then I wait to receive an error
  And the event "unhandled" is false
  And the event "severity" equals "warning"
  And the event "severityReason.type" equals "handledError"
  And the exception "errorClass" equals "*os.PathError"
  And the "file" of stack frame 0 equals "nethttp_scenario.go"

Scenario: A handled error sends a report with a custom name
  Given I set environment variable "ERROR_CLASS" to "MyCustomErrorClass"
  When I start the service "app"
  And I run "HttpServerScenario"
  And I wait for the host "localhost" to open port "4512"
  And I open the URL "http://localhost:4512/handled"
  Then I wait to receive an error
  And the event "unhandled" is false
  And the event "severity" equals "warning"
  And the event "severityReason.type" equals "handledError"
  And the exception "errorClass" equals "MyCustomErrorClass"
  And the "file" of stack frame 0 equals "nethttp_scenario.go"