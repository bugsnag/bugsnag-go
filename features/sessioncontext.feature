Feature: Session data inside an error report using a session context

Scenario: An error report contains a session count when part of a session
  When I start the service "app"
  And I run "SessionAndErrorScenario"
  Then I wait to receive 2 errors
  And the event handled sessions count equals 1 for request 0
  And the event unhandled sessions count equals 0 for request 0
  And the number of sessions started equals 1 for request 1
