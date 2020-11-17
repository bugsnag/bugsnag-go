Feature: Plain handled errors

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I have built the service "app"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"

Scenario: A handled error sends a report
  When I run the go service "app" with the test case "handled"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "unhandled" is false
  And the event "severity" equals "warning"
  And the event "severityReason.type" equals "handledError"
  And the exception "errorClass" equals "*os.PathError"
  And the "file" of stack frame 0 equals "main.go"

Scenario: A handled error sends a report with a custom name
  Given I set environment variable "ERROR_CLASS" to "MyCustomErrorClass"
  When I run the go service "app" with the test case "handled"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "unhandled" is false
  And the event "severity" equals "warning"
  And the event "severityReason.type" equals "handledError"
  And the exception "errorClass" equals "MyCustomErrorClass"
  And the "file" of stack frame 0 equals "main.go"

Scenario: Sending an event using a callback to modify report contents
  When I run the go service "app" with the test case "handled with callback"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "unhandled" is false
  And the event "severity" equals "info"
  And the event "severityReason.type" equals "userCallbackSetSeverity"
  And the event "context" equals "nonfatal.go:14"
  And the "file" of stack frame 0 equals "main.go"
  And stack frame 0 contains a local function spanning 240 to 246
  And the "file" of stack frame 1 equals ">insertion<"
  And the "lineNumber" of stack frame 1 equals 0

Scenario: Marking an error as unhandled in a callback
  When I run the go service "app" with the test case "make unhandled with callback"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "unhandled" is true
  And the event "severity" equals "error"
  And the event "severityReason.type" equals "userCallbackSetSeverity"
  And the event "severityReason.unhandledOverridden" is true
  And the "file" of stack frame 0 equals "main.go"
  And stack frame 0 contains a local function spanning 252 to 255
