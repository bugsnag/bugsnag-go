Feature: Using auto notify

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "SERVER_PORT" to "4511"

Scenario: An error report is sent when an AutoNotified crash occurs which later gets recovered
  Given I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  When I start the service "gin"
  And I wait for the app to open port "4511"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4511/autonotify-then-recover"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "unhandled" is true
  And the event "severityReason.type" equals "unhandledErrorMiddleware" for request 0
  And the event "severityReason.attributes.framework" equals "Gin" for request 0
  And the exception "errorClass" equals "*runtime.TypeAssertionError"
  And the exception "message" matches "interface conversion: interface ({} )?is struct {}, not string"

Scenario: An error report is sent when a go routine crashes which is reported through auto notify
  Given I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  When I start the service "gin"
  And I wait for the app to open port "4511"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4511/autonotify"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "unhandled" is true
  And the event "severityReason.type" equals "handledPanic" for request 0
  And the exception "errorClass" equals "*errors.errorString"
  And the exception "message" equals "Go routine killed with auto notify"
