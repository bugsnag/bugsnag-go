Feature: Capturing request information automatically

Scenario: An error report will automatically contain request information
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  And I set environment variable "AUTO_CAPTURE_SESSION" to "true"
  When I run the http-net test server with the "default" configuration
  And I wait for 1 second
  Then I should receive 2 requests
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "request.clientIp" starts with "127.0.0.1"
  And the event "request.headers.User-Agent" equals "Go-http-client/1.1"
  And the event "request.httpMethod" equals "GET"
  And the event "request.url" ends with "/1234abcd?fish=bird"
  And the event "request.url" starts with "http://127.0.0.1:"