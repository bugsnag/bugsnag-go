Feature: Configuring hostname

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I have built the service "app"

Scenario: An error report contains the configured hostname
  Given I set environment variable "HOSTNAME" to "server-1a"
  When I run the go service "app" with the test case "handled"
  And I wait for 1 second
  Then I should receive a request
  And the request is valid for the error reporting API
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "device.hostname" equals "server-1a"

Scenario: An session report contains the configured hostname
  Given I set environment variable "HOSTNAME" to "server-1a"
  When I run the go service "app" with the test case "session"
  And I wait for 1 second
  Then I should receive a request
  And the request is valid for the session tracking API
  And the session contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the payload field "device.hostname" equals "server-1a"