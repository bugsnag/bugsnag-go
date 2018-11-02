Feature: Configuring param filters

Scenario: Filtering metadata for gin apps
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  Given I set environment variable "PARAMS_FILTERS" to "Name"
  And I configure the bugsnag notify endpoint only
  When I run the script "features/fixtures/gin/run.sh"
  And I go to the gin route "/metadata"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "metaData.Scheme.Customer.Name" equals "[FILTERED]"
