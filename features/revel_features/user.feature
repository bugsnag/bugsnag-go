Feature: Sending user data

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "SERVER_PORT" to "4515"
  And I set environment variable "USE_CODE_CONFIG" to "true"

Scenario Outline: An error report contains custom user data
  Given I set environment variable "REVEL_VERSION" to "<revel version>"
  And I set environment variable "REVEL_CMD_VERSION" to "<revel cmd>"
  And I set environment variable "USER_ID" to "test-user-id"
  And I set environment variable "USER_NAME" to "test-user-name"
  And I set environment variable "USER_EMAIL" to "test-user-email"
  When I start the service "revel"
  And I wait for the app to open port "4515"
  And I wait for 4 seconds
  And I open the URL "http://localhost:4515/user"
  Then I wait to receive 2 requests
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "user.id" equals "test-user-id"
  And the event "user.name" equals "test-user-name"
  And the event "user.email" equals "test-user-email"

  Examples:
  | revel version | revel cmd |
  | v0.21.0       | v0.21.1   |
  | v0.20.0       | v0.20.2   |
  | v0.19.1       | v0.19.0   |
  | v0.18.0       | v0.18.0   |
