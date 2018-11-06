Feature: Configuring app version

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "APP_VERSION" to "3.1.2"
  And I have built the service "app"

Scenario: An error report contains the configured app type when running a go app
  When I run the go service "app" with the test case "handled"
  Then I wait to receive a request
  And the request is valid for the error reporting API
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "app.version" equals "3.1.2"

Scenario: An session report contains the configured app type when running a go app
  When I run the go service "app" with the test case "session"
  Then I wait to receive a request
  And the request is valid for the session tracking API
  And the session contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the payload field "app.version" equals "3.1.2"