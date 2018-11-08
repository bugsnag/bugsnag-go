Feature: Configuring app version

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "APP_VERSION" to "3.1.2"
  And I set environment variable "SERVER_PORT" to "4515"
  And I set environment variable "USE_CODE_CONFIG" to "true"

Scenario Outline: A error report contains the configured app type when using a gin
  Given I set environment variable "REVEL_VERSION" to "<revel version>"
  And I set environment variable "REVEL_CMD_VERSION" to "<revel cmd>"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  When I start the service "revel"
  And I wait for the app to open port "4515"
  And I wait for 4 seconds
  And I open the URL "http://localhost:4515/handled"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "app.version" equals "3.1.2"

  Examples:
  | revel version | revel cmd |
  | v0.21.0       | v0.21.1   |
  | v0.20.0       | v0.20.2   |
  | v0.19.1       | v0.19.0   |
  | v0.18.0       | v0.18.0   |

Scenario Outline: A session report contains the configured app type when using gin
  Given I set environment variable "REVEL_VERSION" to "<revel version>"
  And I set environment variable "REVEL_CMD_VERSION" to "<revel cmd>"
  When I start the service "revel"
  And I wait for the app to open port "4515"
  And I wait for 4 seconds
  And I open the URL "http://localhost:4515/session"
  Then I wait to receive a request
  And the request is a valid session report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the payload field "app.version" equals "3.1.2"

  Examples:
  | revel version | revel cmd |
  | v0.21.0       | v0.21.1   |
  | v0.20.0       | v0.20.2   |
  | v0.19.1       | v0.19.0   |
  | v0.18.0       | v0.18.0   |