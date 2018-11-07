Feature: Sending user data

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "SERVER_PORT" to "4511"

Scenario Outline: An error report contains custom user data
  Given I set environment variable "GIN_VERSION" to "<gin version>"
  And I set environment variable "USER_ID" to "test-user-id"
  And I set environment variable "USER_NAME" to "test-user-name"
  And I set environment variable "USER_EMAIL" to "test-user-email"
  When I start the service "gin"
  And I wait for the app to open port "4511"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4511/user"
  Then I wait to receive 2 requests
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "user.id" equals "test-user-id"
  And the event "user.name" equals "test-user-name"
  And the event "user.email" equals "test-user-email"

  Examples:
  | gin version |
  | v1.3.0      |
  | v1.2        |
  | v1.1        |
  | v1.0        |
