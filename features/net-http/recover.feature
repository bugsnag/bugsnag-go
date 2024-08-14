Feature: Using recover

Scenario: An error report is sent when request crashes but is recovered
  When I start the service "app"
  And I run "HttpServerScenario"
  And I wait for the host "localhost" to open port "4512"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4512/recover"
  Then I wait to receive an error
  And the exception "errorClass" equals "*errors.errorString"
  And the exception "message" equals "Request killed but recovered"
