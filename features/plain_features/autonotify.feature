Feature: Using auto notify

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I have built the service "app"

Scenario: An error report is sent when an AutoNotified crash occurs which later gets recovered
  When I run the go service "app" with the test case "autonotify"
  Then I wait for 3 seconds
  And the request 1 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the exception "errorClass" equals "*errors.errorString" for request 1
  And the exception "message" equals "Go routine killed with auto notify" for request 1
