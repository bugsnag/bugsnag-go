Feature: Panic handling

    Background:
      Given I set environment variable "BUGSNAG_SOURCE_ROOT" to "/app/src/features/fixtures/app/"
      And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "false"

    Scenario: Capturing a panic
      When I start the service "app"
      And I run "UnhandledScenario"
      And I wait to receive an error
      And the event "unhandled" is true
      And the event "severity" equals "error"
      And the event "severityReason.type" equals "unhandledPanic"
      And the exception "errorClass" equals "panic"
      #And the exception "message" is one of:
      #  | interface conversion: interface is struct {}, not string      |
      #  | interface conversion: interface {} is struct {}, not string   |
      #And the in-project frames of the stacktrace are:
      #  | file    | method               |
      #  | main.go | unhandledCrash.func1 |
      #  | main.go | unhandledCrash       |
