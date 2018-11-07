Feature: Capturing request information automatically

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "SERVER_PORT" to "4511"

Scenario Outline: An error report will automatically contain request information
  Given I set environment variable "GIN_VERSION" to "<gin version>"
  When I start the service "gin"
  And I wait for the app to open port "4511"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4511/handled"
  Then I wait to receive 2 requests
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "request.clientIp" is not null
  And the event "request.headers.User-Agent" equals "Ruby"
  And the event "request.httpMethod" equals "GET"
  And the event "request.url" ends with "/handled"
  And the event "request.url" starts with "http://"
    
  Examples:
  | gin version |
  | v1.3.0      |
  | v1.2        |
  | v1.1        |
  | v1.0        |
