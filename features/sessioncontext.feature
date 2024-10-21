Feature: Session data inside an error report using a session context

Scenario: An error report contains a session count when part of a session
  When I start the service "app"
  And I run "SessionAndErrorScenario"
  Then I wait to receive 1 error
  # one session is created on start
  And I wait to receive 2 session
  And I discard the oldest session
  And the session payload has a valid sessions array