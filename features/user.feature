Feature: Sending user data

Scenario: An error report contains custom user data
  When I start the service "app"
  And I run "SetUserScenario"
  And I wait to receive an error
  And the event "user.id" equals "test-user-id"
  And the event "user.name" equals "test-user-name"
  And the event "user.email" equals "test-user-email"