Feature: Using auto notify

Scenario: An error report is sent when a go routine crashes
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  When I configure with the "auto notify" configuration and send an error with crash
  And I wait for 1 second
  Then I should receive 2 requests
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the payload field "events.0.exceptions.0.errorClass" equals "*errors.errorString" for request 0
  And the payload field "events.0.exceptions.0.message" equals "Go routine killed" for request 0
  And the payload field "events.0.exceptions.0.errorClass" equals "panic" for request 1
  And the payload field "events.0.exceptions.0.message" equals "Go routine killed [recovered]" for request 1