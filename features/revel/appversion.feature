Feature: Configuring app version

Scenario: Revel configuration through code
  Given I set environment variable "USE_CODE_CONFIG" to "YES"
  Given I set environment variable "APP_VERSION" to "1.33.7"
  And I work with a new 'revel' app
  And I set the "revel" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint in the config file for 'revel'
  When I run the script "features/fixtures/revel/run.sh"
  And I wait for 4 seconds
  And I go to the route "/configure"
  And I wait for 1 seconds
  Then I should receive a request
  And the event "app.version" equals "1.33.7"

Scenario: Revel configuration through config file
  And I work with a new 'revel' app
  And I set the "revel" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I set the "revel" config variable "bugsnag.appversion" to "1.33.7"
  And I configure the bugsnag endpoint in the config file for 'revel'
  When I run the script "features/fixtures/revel/run.sh"
  And I wait for 4 seconds
  And I go to the route "/configure"
  And I wait for 1 seconds
  Then I should receive a request
  And the event "app.version" equals "1.33.7"

Scenario: A Revel session report contains the configured app version
  Given I set environment variable "DISABLE_REPORT_PAYLOADS" to "true"
  And I work with a new 'revel' app
  And I set the "revel" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I set the "revel" config variable "bugsnag.appversion" to "1.33.7"
  And I configure the bugsnag sessions endpoint in the config file for 'revel'
  When I run the script "features/fixtures/revel/run.sh"
  And I wait for 4 seconds
  And I go to the route "/configure"
  And I wait for 1 seconds
  Then I should receive a request
  And the payload field "app.version" equals "1.33.7"
  And the payload field "sessionCounts.0.sessionsStarted" equals 1
