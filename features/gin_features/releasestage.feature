Feature: Configuring release stage

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "SERVER_PORT" to "4511"
  And I set environment variable "RELEASE_STAGE" to "my-stage"

Scenario: An error report is sent with configured release stage
  Given I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  When I start the service "gin"
  And I wait for the app to open port "4511"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4511/handled"
  Then I wait to receive a request
  And the request 0 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "app.releaseStage" equals "my-stage" for request 0

Scenario: A session report contains the configured app type
  Given I set environment variable "AUTO_CAPTURE_SESSIONS" to "true"
  When I start the service "gin"
  And I wait for the app to open port "4511"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4511/session"
  Then I wait to receive a request after the start up session
  And the request is a valid session report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the payload field "app.releaseStage" equals "my-stage"
