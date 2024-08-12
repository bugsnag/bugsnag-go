Feature: Session data inside an error report using a session context

Scenario: An error report contains a session count when part of a session
  When I start the service "app"
  And I run "SessionAndErrorScenario"
  Then I wait to receive 1 error
  And I wait to receive 1 session
  And the error is valid for the error reporting API version "4" for the "Bugsnag Go" notifier
  And the session is valid for the session reporting API version "1.0" for the "Bugsnag Go" notifier
  And the session payload has a valid sessions array