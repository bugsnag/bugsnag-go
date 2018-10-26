Feature: Plain handled errors

Scenario: A handled error sends a report
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  When I configure with the "default" configuration and send an error
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "unhandled" is false
  And the event "severity" equals "warning"
  And the event "severityReason.type" equals "handledError"
  And the exception "errorClass" equals "*errors.errorString"
  And the "file" of stack frame 0 equals "main.go"
