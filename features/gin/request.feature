Feature: Capturing request information automatically

Scenario: A gin error report will automatically contain request information
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag notify endpoint only
  When I run the script "features/fixtures/gin/run.sh"
  And I go to the gin route "/handled"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "request.httpMethod" equals "GET"
