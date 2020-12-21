Feature: Session data inside an error report using a session context

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I have built the service "app"

Scenario: An error report contains a session count when part of a session
  When I run the go service "app" with the test case "session-and-error"
  Then I wait to receive 2 requests after the start up session
  And the request 0 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the request 1 is a valid session report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event handled sessions count equals 1 for request 0
  And the event unhandled sessions count equals 0 for request 0
  And the number of sessions started equals 1 for request 1
