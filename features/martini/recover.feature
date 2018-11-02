Feature: Using recover

Scenario: Recovering metadata for martini apps
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag notify endpoint only
  When I run the script "features/fixtures/martini/run.sh"
  And I go to the martini route "/recover"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "unhandled" is false
  And the exception "errorClass" equals "*runtime.TypeAssertionError"
  And the exception "message" equals "interface conversion: interface {} is struct {}, not string"
