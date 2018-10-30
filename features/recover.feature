Feature: Using auto notify

Scenario: An error report is sent when a go routine crashes
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  When I configure with the "recover" configuration and send an error with crash
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the exception "errorClass" equals "*errors.errorString"
  And the exception "message" equals "Go routine killed"