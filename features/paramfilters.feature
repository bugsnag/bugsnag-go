Feature: Configuring param filters

Scenario: An error report containing meta data is not filtered when the param filters are set but do not match
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "PARAMS_FILTERS" to "Name"
  When I configure with the "params filters" configuration and send an error
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "metaData.Account.Price(dollars)" equals "1 Million"

Scenario: An error report containing meta data is filtered when the param filters are set and completely match
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "PARAMS_FILTERS" to "Price(dollars)"
  When I configure with the "params filters" configuration and send an error
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "metaData.Account.Price(dollars)" equals "[FILTERED]"

Scenario: An error report containing meta data is filtered when the param filters are set and partially match
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "PARAMS_FILTERS" to "Price"
  When I configure with the "params filters" configuration and send an error
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "metaData.Account.Price(dollars)" equals "[FILTERED]"

