Feature: Configuring app version

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "APP_VERSION" to "3.1.2"
  And I set environment variable "SERVER_PORT" to "4513"

Scenario: A error report contains the configured app type
  Given I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  When I start the service "martini"
  And I wait for the app to open port "4513"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4513/handled"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "app.version" equals "3.1.2"

Scenario: A session report contains the configured app type
  Given I set environment variable "AUTO_CAPTURE_SESSIONS" to "true"
  When I start the service "martini"
  And I wait for the app to open port "4513"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4513/session"
  Then I wait to receive a request after the start up session
  And the request is a valid session report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the payload field "app.version" equals "3.1.2"

