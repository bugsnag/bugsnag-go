Feature: Using auto notify

Scenario: An error report is sent when an AutoNotified crash occurs which later gets recovered
  Given I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
  When I start the service "app"
  And I run "HttpServerScenario"
  And I wait for the host "localhost" to open port "4512"
  And I open the URL "http://localhost:4512/autonotify-then-recover"
  Then I wait to receive an error
  And the event "unhandled" is true
  And the exception "errorClass" equals "*runtime.TypeAssertionError"
  And the exception "message" matches "interface conversion: interface ({} )?is struct {}, not string"

Scenario: An error report is sent when a go routine crashes which is reported through auto notify
  Given I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
  When I start the service "app"
  And I run "HttpServerScenario"
  And I wait for the host "localhost" to open port "4512"
  And I open the URL "http://localhost:4512/autonotify"
  Then I wait to receive an error
  And the event "unhandled" is true
  And the exception "errorClass" equals "*errors.errorString"
  And the exception "message" equals "Go routine killed with auto notify"