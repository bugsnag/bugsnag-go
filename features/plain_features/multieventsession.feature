Feature: Reporting multiple handled and unhandled errors in the same session

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I have built the service "app"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"

Scenario: Handled errors know about previous reported handled errors
  When I run the go service "app" with the test case "multiple-handled"
  And I wait to receive 2 requests
  And the request 0 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the request 1 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event handled sessions count equals 1 for request 0
  And the event handled sessions count equals 2 for request 1

Scenario: Unhandled errors know about previous reported handled errors
  When I run the go service "app" with the test case "multiple-unhandled"
  And I wait to receive 2 requests
  And the request 0 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the request 1 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event unhandled sessions count equals 1 for request 0
  And the event unhandled sessions count equals 2 for request 1
