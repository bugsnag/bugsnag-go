Feature: Using auto notify

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I have built the service "app"

Scenario: An error report is sent when an unhandled crash occurs
  When I run the go service "app" with the test case "unhandled"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the exception "errorClass" equals "panic"
  And the exception "message" equals "interface conversion: interface {} is struct {}, not string"

Scenario: An error report is sent when a go routine crashes which is protected by auto notify
  When I run the go service "app" with the test case "autonotify"
  Then I wait to receive 2 requests
  
  And the request 0 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the exception "errorClass" equals "*errors.errorString" for request 0
  And the exception "message" equals "Go routine killed with auto notify" for request 0

  And the request 1 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the exception "errorClass" equals "panic" for request 1
  And the exception "message" equals "Go routine killed with auto notify [recovered]" for request 1