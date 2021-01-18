Feature: Configuring endpoint

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I have built the service "app"

Scenario: An error report is sent successfully using the notify endpoint only
  When I run the go service "app" with the test case "endpoint-notify"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"

Scenario: Configuring Bugsnag will panic if the sessions endpoint is configured without the notify endpoint
  When I run the go service "app" with the test case "endpoint-session"
  And I wait for 3 second
  Then I should receive no requests
