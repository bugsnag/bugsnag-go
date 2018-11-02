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
