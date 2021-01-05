Feature: Configuring synchronous flag

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  And I configure the bugsnag endpoint
  And I have built the service "app"

Scenario: An error report is sent asynchrously but exits immediately so is not sent
  Given I set environment variable "SYNCHRONOUS" to "false"
  When I run the go service "app" with the test case "send-and-exit"
  And I wait for 3 second
  Then I should receive no requests

Scenario: An error report is report synchronously so it will send before exiting
  Given I set environment variable "SYNCHRONOUS" to "true"
  When I run the go service "app" with the test case "send-and-exit"
  Then I wait to receive 1 requests
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"

