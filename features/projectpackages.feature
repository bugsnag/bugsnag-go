Feature: Configuring project packages

Scenario: An error report contains the configured app type
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "PROJECT_PACKAGES" to "bob"
  When I configure with the "project packages" configuration and send an error
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the "file" of stack frame 2 equals "bob"

