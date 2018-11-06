Feature: Session data inside an error report using a session context

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I have built the service "app"

Scenario: An error report contains a session count when part of a session
  When I run the go service "app" with the test case "session and error"
  And I wait for 1 second
  Then I should receive 2 requests
  And the request 0 is valid for the error reporting API
  And the request 0 contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the request 1 is valid for the session tracking API
  And the session in request 1 contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the events handled sessions count equals 1 for request 0
  And the events unhandled sessions count equals 0 for request 0
  And the number of sessions started equals 1 in request 1
