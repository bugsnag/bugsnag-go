Feature: Panic handling

    Background:
      Given I set environment variable "BUGSNAG_SOURCE_ROOT" to "/app/src/features/fixtures/app/"
      And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"

    Scenario: Capturing a panic
      When I start the service "app"
      And I run "UnhandledCrashScenario"
      And I wait to receive an error
      And the event "unhandled" is true
      And the event "severity" equals "error"
      And the event "severityReason.type" equals "unhandledPanic"
      And the exception "errorClass" equals "panic"
      And the exception "message" matches "^interface conversion: interface.*?is struct {}, not string"
      And the "file" of stack frame 0 equals "unhandled_scenario.go"
      And the "method" of stack frame 0 equals "UnhandledCrashScenario.func1.1"
      And the "file" of stack frame 1 equals "unhandled_scenario.go"
      And the "method" of stack frame 1 equals "UnhandledCrashScenario.func1"