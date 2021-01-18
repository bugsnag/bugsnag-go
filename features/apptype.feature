Feature: Configuring app type

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "APP_TYPE" to "background-queue"
  And I have built the service "app"

Scenario: An error report contains the configured app type when running a go app
  Given I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  When I run the go service "app" with the test case "handled"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "app.type" equals "background-queue"

Scenario: An session report contains the configured app type when running a go app
  Given I set environment variable "AUTO_CAPTURE_SESSIONS" to "true"
  When I run the go service "app" with the test case "session"
  Then I wait to receive a request after the start up session
  And the request is a valid session report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the payload field "app.type" equals "background-queue"

