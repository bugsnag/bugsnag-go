Feature: Configuring param filters

Scenario: Filtering metadata for martini apps
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  Given I set environment variable "PARAMS_FILTERS" to "Name"
  And I am working with a new martini app
  And I configure the bugsnag notify endpoint only
  When I run the script "features/fixtures/martini/run.sh"
  And I go to the martini route "/metadata"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "metaData.Scheme.Customer.Name" equals "[FILTERED]"

Scenario: An error report containing meta data is not filtered when the param filters are set but do not match
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  And I set environment variable "PARAMS_FILTERS" to "Name"
  When I configure with the "params filters" configuration and send an error
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "metaData.Account.Price(dollars)" equals "1 Million"

Scenario: An error report containing meta data is filtered when the param filters are set and completely match
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  And I set environment variable "PARAMS_FILTERS" to "Price(dollars)"
  When I configure with the "params filters" configuration and send an error
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "metaData.Account.Price(dollars)" equals "[FILTERED]"

Scenario: An error report containing meta data is filtered when the param filters are set and partially match
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  And I set environment variable "PARAMS_FILTERS" to "Price"
  When I configure with the "params filters" configuration and send an error
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "metaData.Account.Price(dollars)" equals "[FILTERED]"

