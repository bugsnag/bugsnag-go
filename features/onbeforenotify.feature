Feature: Configuring on before notify

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I have built the service "app"

Scenario: Send three bugsnags and use on before notify to drop one and modify the message of another
  When I run the go service "app" with the test case "onbeforenotify"
  Then I wait to receive 2 requests after the start up session
  And the request 0 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the exception "message" equals "Don't ignore this error" for request 0
  And the request 1 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the exception "message" equals "Error message was changed" for request 1
