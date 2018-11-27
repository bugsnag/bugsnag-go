Feature: Using recover

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "SERVER_PORT" to "4512"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"

Scenario: An error report is sent when request crashes but is recovered
  When I start the service "nethttp"
  And I wait for the app to open port "4512"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4512/recover"
  Then I wait to receive a request
  And the request 0 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the exception "errorClass" equals "*errors.errorString" for request 0
  And the exception "message" equals "Request killed but recovered" for request 0
