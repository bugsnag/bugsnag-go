Feature: Using recover

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  And I configure the bugsnag endpoint
  And I have built the service "app"

Scenario: An error report is sent when a go routine crashes but recovers
  When I run the go service "app" with the test case "recover"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the exception "errorClass" equals "*errors.errorString"
  And the exception "message" equals "Go routine killed but recovered"
