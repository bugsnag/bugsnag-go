Feature: Session data inside an error report using a session context

Scenario: An error report contains a session count when part of a session
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  When I configure with the "session" configuration and send an error
  And I wait for 1 second
  Then I should receive 2 requests
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the payload field "events.0.session.events.handled" equals 1 for request 0
  And the payload field "sessionCounts.0.sessionsStarted" equals 1 for request 1
