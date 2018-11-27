Feature: Sending user data

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  Given I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  And I configure the bugsnag endpoint
  And I have built the service "app"

Scenario: An error report contains custom user data
  Given I set environment variable "USER_ID" to "test-user-id"
  And I set environment variable "USER_NAME" to "test-user-name"
  And I set environment variable "USER_EMAIL" to "test-user-email"
  When I run the go service "app" with the test case "user"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "user.id" equals "test-user-id"
  And the event "user.name" equals "test-user-name"
  And the event "user.email" equals "test-user-email"
