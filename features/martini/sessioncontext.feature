Feature: Session data inside an error report using a session context

Scenario: An error report contains a session count when part of a session for martini
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag notify endpoint only
  When I run the script "features/fixtures/martini/run.sh"
  And I go to the martini route "/handled"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the payload field "events.0.session.events.handled" equals 1 for request 0
