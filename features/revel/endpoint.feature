Feature: Configuring endpoint

Scenario: Revel legacy endpoint is picked up in code config
  Given I set environment variable "USE_CODE_CONFIG" to "YES"
  And I configure the legacy bugsnag endpoint in as an environment variable
  And I work with a new 'revel-0.20.0' app
  And I set the "revel-0.20.0" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  When I run the script "features/fixtures/revel-0.20.0/run.sh"
  And I wait for 4 seconds
  And I go to the route "/configure"
  And I wait for 1 seconds
  Then I should receive a request

Scenario: Revel configuration through config file
  Given I work with a new 'revel-0.20.0' app
  And I set the "revel-0.20.0" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the legacy bugsnag endpoint in the config file for "revel-0.20.0"
  When I run the script "features/fixtures/revel-0.20.0/run.sh"
  And I wait for 4 seconds
  And I go to the route "/configure"
  And I wait for 1 seconds
  Then I should receive a request
