Feature: Configuring param filters

Scenario: An error report containing meta data is not filtered when the param filters are set but do not match
  Given I set environment variable "BUGSNAG_PARAMS_FILTERS" to "Name"
  When I start the service "app"
  And I run "FilteredMetadataScenario"
  And I wait to receive an error
  And the event "metaData.Account.Price(dollars)" equals "1 Million"

Scenario: An error report containing meta data is filtered when the param filters are set and completely match
  Given I set environment variable "BUGSNAG_PARAMS_FILTERS" to "Price(dollars)"
  When I start the service "app"
  And I run "FilteredMetadataScenario"
  And I wait to receive an error
  And the event "metaData.Account.Price(dollars)" equals "[FILTERED]"

Scenario: An error report containing meta data is filtered when the param filters are set and partially match
  Given I set environment variable "BUGSNAG_PARAMS_FILTERS" to "Price"
  When I start the service "app"
  And I run "FilteredMetadataScenario"
  And I wait to receive an error
  And the event "metaData.Account.Price(dollars)" equals "[FILTERED]"
