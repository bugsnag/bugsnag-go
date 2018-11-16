Feature: Configure auto capture sessions

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "SERVER_PORT" to "4511"

Scenario: A session is not sent if auto capture sessions is off
  Given I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  When I start the service "gin"
  And I wait for the app to open port "4511"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4511/session"
  And I wait for 2 seconds
  Then I should receive no requests

Scenario: A session is sent if auto capture sessions is on
  Given I set environment variable "AUTO_CAPTURE_SESSIONS" to "true"
  When I start the service "gin"
  And I wait for the app to open port "4511"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4511/session"
  Then I wait to receive a request after the start up session
  And the request is a valid session report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
