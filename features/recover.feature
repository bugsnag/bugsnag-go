Feature: Using recover

Scenario: Recovering for negroni apps
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I am working with a new negroni app
  And I configure the bugsnag notify endpoint only
  When I run the script "features/fixtures/negroni/run.sh"
  And I go to the negroni route "/recover"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "unhandled" is false
  And the exception "errorClass" equals "*runtime.TypeAssertionError"
  And the exception "message" equals "interface conversion: interface {} is struct {}, not string"

Scenario: Recovering metadata for martini apps
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I am working with a new martini app
  And I configure the bugsnag notify endpoint only
  When I run the script "features/fixtures/martini/run.sh"
  And I go to the martini route "/recover"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "unhandled" is false
  And the exception "errorClass" equals "*runtime.TypeAssertionError"
  And the exception "message" equals "interface conversion: interface {} is struct {}, not string"

Scenario: An error report is sent when a go routine crashes
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  When I configure with the "recover" configuration and send an error with crash
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the exception "errorClass" equals "*errors.errorString"
  And the exception "message" equals "Go routine killed"

Scenario: An error report is sent when a go routine crashes when using net http
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  When I run the http-net test server with the "recover" configuration
  And I wait for 1 second
  Then I should receive 2 requests
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the exception "errorClass" equals "*errors.errorString"
  And the exception "message" equals "Go routine killed"
  And the payload field "events.0.session.events.handled" equals 1 for request 0
  And the payload field "sessionCounts.0.sessionsStarted" equals 1 for request 1
