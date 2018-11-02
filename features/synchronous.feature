Feature: Configuring synchronous flag

Scenario: An error report is sent asynchrously but exits immediately so is not sent
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  And I set environment variable "SYNCHRONOUS" to "false"
  When I configure with the "synchronous" configuration and send an error
  And I wait for 1 second
  Then I should receive no requests

Scenario: An error report is report synchronously so it will send before exiting
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  And I set environment variable "SYNCHRONOUS" to "true"
  When I configure with the "synchronous" configuration and send an error
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
