Feature: Configuring param filters

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  And I configure the bugsnag endpoint
  And I have built the service "app"

Scenario: An error report containing meta data is not filtered when the param filters are set but do not match
  Given I set environment variable "PARAMS_FILTERS" to "Name"
  When I run the go service "app" with the test case "filtered"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "metaData.Account.Price(dollars)" equals "1 Million"

Scenario: An error report containing meta data is filtered when the param filters are set and completely match
  Given I set environment variable "PARAMS_FILTERS" to "Price(dollars)"
  When I run the go service "app" with the test case "filtered"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "metaData.Account.Price(dollars)" equals "[FILTERED]"

Scenario: An error report containing meta data is filtered when the param filters are set and partially match
  Given I set environment variable "PARAMS_FILTERS" to "Price"
  When I run the go service "app" with the test case "filtered"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "metaData.Account.Price(dollars)" equals "[FILTERED]"

