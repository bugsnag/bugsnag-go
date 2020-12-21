Feature: Panic handling

    Background:
      Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
      And I configure the bugsnag endpoint
      And I have built the service "app"
      And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"

    Scenario: Capturing a panic
      When I run the go service "app" with the test case "unhandled"
      Then I wait to receive a request
      And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
      And the event "unhandled" is true
      And the event "severity" equals "error"
      And the event "severityReason.type" equals "unhandledPanic"
      And the exception "errorClass" equals "panic"
      And the exception "message" is one of:
        | interface conversion: interface is struct {}, not string      |
        | interface conversion: interface {} is struct {}, not string   |
      And the in-project frames of the stacktrace are:
        | file    | method               |
        | main.go | unhandledCrash.func1 |
        | main.go | unhandledCrash       |
