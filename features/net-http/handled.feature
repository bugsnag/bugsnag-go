Feature: Handled errors

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  Given I set environment variable "BUGSNAG_SOURCE_ROOT" to the app directory
  And I configure the bugsnag endpoint
  And I set environment variable "SERVER_PORT" to "4512"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  And I have built the service "nethttp"

Scenario: A handled error sends a report
  When I start the service "nethttp"
  And I wait for the app to open port "4512"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4512/handled"
  Then I wait to receive a request
  And the request 0 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "unhandled" is false for request 0
  And the event "severity" equals "warning" for request 0
  And the event "severityReason.type" equals "handledError" for request 0
  And the exception is a PathError for request 0
  And the "file" of stack frame 0 equals "main.go" for request 0

Scenario: A handled error sends a report with a custom name
  Given I set environment variable "ERROR_CLASS" to "MyCustomErrorClass"
  When I start the service "nethttp"
  And I wait for the app to open port "4512"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4512/handled"
  Then I wait to receive a request
  And the request 0 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "unhandled" is false for request 0
  And the event "severity" equals "warning" for request 0
  And the event "severityReason.type" equals "handledError" for request 0
  And the exception "errorClass" equals "MyCustomErrorClass" for request 0
  And the "file" of stack frame 0 equals "main.go" for request 0
