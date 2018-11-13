Feature: Configuring release stage

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "SERVER_PORT" to "4515"

Scenario Outline: An error report and session is sent with configured release stage
  Given I set environment variable "RELEASE_STAGE" to "my-stage"
  When I start the service "revel"
  And I wait for the app to open port "4515"
  And I wait for 4 seconds
  And I open the URL "http://localhost:4515/handled"
  Then I wait to receive 2 requests
  And the request 0 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "app.releaseStage" equals "my-stage" for request 0
  And the request 1 is a valid session report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the payload field "app.releaseStage" equals "my-stage" for request 1
  And the event handled sessions count equals 1 for request 0
  And the number of sessions started equals 1 for request 1
