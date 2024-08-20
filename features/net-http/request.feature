Feature: Capturing request information automatically

Scenario: An error report will automatically contain request information
  When I start the service "app"
  And I run "HttpServerScenario"
  And I wait for the host "localhost" to open port "4512"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4512/handled"
  Then I wait to receive an error
  And the event "request.clientIp" is not null
  And the event "request.headers.User-Agent" equals "Ruby"
  And the event "request.httpMethod" equals "GET"
  And the event "request.url" ends with "/handled"
  And the event "request.url" starts with "http://"