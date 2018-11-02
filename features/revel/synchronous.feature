Feature: Configuring synchronous flag

Scenario: Revel configuration through code
  Given I set environment variable "USE_CODE_CONFIG" to "YES"
  Given I set environment variable "SYNCHRONOUS" to "true"
  And I work with a new 'revel' app
  And I set the "revel" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint in the config file for 'revel'
  When I run the script "features/fixtures/revel/run.sh"
  And I wait for 4 seconds
  And I go to the route "/synchronous"
  And I wait for 1 seconds
  Then I should receive a request

Scenario: Revel configuration through config file
  And I work with a new 'revel' app
  And I set the "revel" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I set the "revel" config variable "bugsnag.synchronous" to "true"
  And I configure the bugsnag endpoint in the config file for 'revel'
  When I run the script "features/fixtures/revel/run.sh"
  And I wait for 4 seconds
  And I go to the route "/synchronous"
  And I wait for 1 seconds
  Then I should receive a request
