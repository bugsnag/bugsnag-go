Feature: Sending user data

Scenario: An error report contains custom user data
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  And I set environment variable "USER_ID" to "test-user-id"
  And I set environment variable "USER_NAME" to "test-user-name"
  And I set environment variable "USER_EMAIL" to "test-user-email"
  When I configure with the "user data" configuration and send an error
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "user.id" equals "test-user-id"
  And the event "user.name" equals "test-user-name"
  And the event "user.email" equals "test-user-email"