Feature: Configuring on before notify

Scenario: Send three bugsnags and use on before notify to drop one and modify the message of another
  When I start the service "app"
  And I run "HttpServerScenario"
  And I wait for the host "localhost" to open port "4512"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4512/onbeforenotify"
  Then I wait to receive 2 errors
  And the exception "message" equals "Don't ignore this error"
  And I discard the oldest error
  And the exception "message" equals "Error message was changed"