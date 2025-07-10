Feature: Plain handled errors

Background:
  Given I set environment variable "BUGSNAG_SOURCE_ROOT" to "/app/src/features/fixtures/app/"
  And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"

Scenario: A handled error sends a report
  When I start the service "app"
  And I run "HandledErrorScenario"
  And I wait to receive an error
  And the event "unhandled" is false
  And the event "severity" equals "warning"
  And the event "severityReason.type" equals "handledError"
  And the exception "errorClass" matches "\*os.PathError|\*fs.PathError"
  And the "file" of stack frame 0 equals "handled_scenario.go"

Scenario: A handled error sends a report with a custom name
  Given I set environment variable "ERROR_CLASS" to "MyCustomErrorClass"
  When I start the service "app"
  And I run "HandledErrorScenario"
  And I wait to receive an error
  And the event "unhandled" is false
  And the event "severity" equals "warning"
  And the event "severityReason.type" equals "handledError"
  And the exception "errorClass" equals "MyCustomErrorClass"
  And the "file" of stack frame 0 equals "handled_scenario.go"

Scenario: Sending an event using a callback to modify report contents
  When I start the service "app"
  And I run "HandledCallbackErrorScenario"
  And I wait to receive an error
  And the event "unhandled" is false
  And the event "severity" equals "info"
  And the event "severityReason.type" equals "userCallbackSetSeverity"
  And the event "context" equals "nonfatal.go:14"
  And the "file" of stack frame 0 equals "handled_scenario.go"
  And the "lineNumber" of stack frame 0 equals 65
  And the "file" of stack frame 1 equals ">insertion<"
  And the "lineNumber" of stack frame 1 equals 0

Scenario: Marking an error as unhandled in a callback
  When I start the service "app"
  And I run "HandledToUnhandledScenario"
  And I wait to receive an error
  And the event "unhandled" is true
  And the event "severity" equals "error"
  And the event "severityReason.type" equals "userCallbackSetSeverity"
  And the event "severityReason.unhandledOverridden" is true
  And the "file" of stack frame 0 equals "handled_scenario.go"
  And the "lineNumber" of stack frame 0 equals 78

Scenario: Unwrapping the causes of a handled error
  When I start the service "app"
  And I run "NestedHandledErrorScenario"
  And I wait to receive an error
  And the event "unhandled" is false
  And the event "severity" equals "warning"
  And the event "exceptions.0.message" equals "terminate process"
  And the "lineNumber" of stack frame 0 equals 46
  And the "file" of stack frame 0 equals "handled_scenario.go"
  And the "method" of stack frame 0 equals "NestedHandledErrorScenario.func1"
  And the event "exceptions.1.message" equals "login failed"
  And the event "exceptions.1.stacktrace.0.file" equals "utils.go"
  And the event "exceptions.1.stacktrace.0.lineNumber" equals 39
  And the event "exceptions.2.message" equals "invalid token"
  And the event "exceptions.2.stacktrace.0.file" equals "utils.go"
  And the event "exceptions.2.stacktrace.0.lineNumber" equals 47