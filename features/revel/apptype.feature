Feature: Configuring app type

Scenario: Revel configuration through code
  Given I set environment variable "USE_CODE_CONFIG" to "YES"
  Given I set environment variable "APP_TYPE" to "Core API"
  And I work with a new 'revel' app
  And I set the "revel" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint in the config file for 'revel'
  When I run the script "features/fixtures/revel/run.sh"
  And I wait for 8 seconds
  And I go to the route "/configure"
  And I wait for 1 seconds
  Then I should receive a request
  And the event "app.type" equals "Core API"

Scenario: Revel configuration through config file
  And I work with a new 'revel' app
  And I set the "revel" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I set the "revel" config variable "bugsnag.apptype" to "Core API"
  And I configure the bugsnag endpoint in the config file for 'revel'
  When I run the script "features/fixtures/revel/run.sh"
  And I wait for 4 seconds
  And I go to the route "/configure"
  And I wait for 1 seconds
  Then I should receive a request
  And the event "app.type" equals "Core API"
