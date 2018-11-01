Feature: Using auto notify

Scenario: An error report is sent when a go routine crashes
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  When I configure with the "auto notify" configuration and send an error with crash
  And I wait for 1 second
  Then I should receive 2 requests
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the payload field "events.0.exceptions.0.errorClass" equals "*errors.errorString" for request 0
  And the payload field "events.0.exceptions.0.message" equals "Go routine killed" for request 0
  And the payload field "events.0.exceptions.0.errorClass" equals "panic" for request 1
  And the payload field "events.0.exceptions.0.message" equals "Go routine killed [recovered]" for request 1

Scenario: Revel panics are captured automatically
  Given I work with a new 'revel-0.20.0' app
  And I set the "revel-0.20.0" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint in the config file for 'revel-0.20.0'
  When I run the script "features/fixtures/revel-0.20.0/run.sh"
  And I wait for 4 seconds
  And I go to the route "/unhandled"
  And I wait for 1 seconds
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "app.releaseStage" equals "dev"
  And the event "app.type" equals "Revel"
  And the event "context" equals "App.Unhandled"
  And the event "request.httpMethod" equals "GET"
  And the event "request.url" equals "http://localhost:9020/unhandled"
  And the event "session.events.handled" equals 0
  And the event "session.events.unhandled" equals 1
  And the event "severity" equals "error"
  And the event "severityReason.type" equals "unhandledErrorMiddleware"
  And the event "unhandled" is true
  And the exception "errorClass" equals "*runtime.TypeAssertionError"
  And the "file" of stack frame 8 equals "controllers/app.go"
  And the "lineNumber" of stack frame 8 equals 26

Scenario: An unhandled error report is sent when a go routine crashes when using net http
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  When I run the http-net test server with the "auto notify" configuration and crashes
  And I wait for 1 second
  Then I should receive 3 requests
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the payload field "events.0.exceptions.0.errorClass" equals "*errors.errorString" for request 0
  And the payload field "events.0.exceptions.0.message" equals "Go routine killed" for request 0
  And the payload field "events.0.exceptions.0.errorClass" equals "panic" for request 2
  And the payload field "events.0.exceptions.0.message" equals "Go routine killed [recovered]" for request 2
  And the payload field "events.0.session.events.unhandled" equals 1 for request 0
  And the payload field "sessionCounts.0.sessionsStarted" equals 1 for request 1
