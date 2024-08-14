Feature: Sending user data

Scenario: An error report contains custom user data
  When I start the service "app"
  And I run "HttpServerScenario"
  And I wait for the host "localhost" to open port "4512"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4512/user"
  Then I wait to receive an error
  And the event "user.id" equals "test-user-id"
  And the event "user.name" equals "test-user-name"
  And the event "user.email" equals "test-user-email"
