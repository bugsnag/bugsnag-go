Feature: Using recover

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I have built the service "gin"
  And I stop the service "gin"
  And I set environment variable "SERVER_PORT" to "4511"

Scenario: An error report and session is sent when request crashes but is recovered
  When I start the service "gin"
  And I wait for the app to open port "4511"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4511/recover"
  Then I wait to receive 2 requests
  And the request 0 is valid for the error reporting API
  And the request 0 contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"

  And the exception "errorClass" equals "*errors.errorString"
  And the exception "message" equals "Request killed but recovered"

  And the request 1 is valid for the session tracking API
  And the session in request 1 contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"

  And the events handled sessions count equals 1 for request 0
  And the number of sessions started equals 1 in request 1
