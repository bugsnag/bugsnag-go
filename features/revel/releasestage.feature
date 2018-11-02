Feature: Configuring release stages and notify release stages

Scenario: Revel configuration through code does not prevent other release stages from going through
  Given I set environment variable "USE_CODE_CONFIG" to "YES"
  Given I set environment variable "NOTIFY_RELEASE_STAGES" to "development"
  Given I set environment variable "RELEASE_STAGE" to "development"
  And I work with a new 'revel' app
  And I set the "revel" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint in the config file for 'revel'
  When I run the script "features/fixtures/revel/run.sh"
  And I wait for 4 seconds
  And I go to the route "/configure"
  And I wait for 1 seconds
  Then I should receive a request
  And the event "app.releaseStage" equals "development"

Scenario: Revel configuration through code can prevent a report from being sent
  Given I set environment variable "USE_CODE_CONFIG" to "YES"
  Given I set environment variable "NOTIFY_RELEASE_STAGES" to "development"
  Given I set environment variable "RELEASE_STAGE" to "staging"
  And I work with a new 'revel' app
  And I set the "revel" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint in the config file for 'revel'
  When I run the script "features/fixtures/revel/run.sh"
  And I wait for 4 seconds
  And I go to the route "/configure"
  And I wait for 1 seconds
  Then I should receive no requests

Scenario: Revel configuration through config file does not prevent other release stages from going through
  And I work with a new 'revel' app
  And I set the "revel" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I set the "revel" config variable "bugsnag.notifyreleasestages" to "development"
  And I set the "revel" config variable "bugsnag.releasestage" to "development"
  And I configure the bugsnag endpoint in the config file for 'revel'
  When I run the script "features/fixtures/revel/run.sh"
  And I wait for 4 seconds
  And I go to the route "/configure"
  And I wait for 1 seconds
  Then I should receive a request
  And the event "app.releaseStage" equals "development"

Scenario: Revel configuration through config file can prevent a report from being sent
  And I work with a new 'revel' app
  And I set the "revel" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I set the "revel" config variable "bugsnag.notifyreleasestages" to "development"
  And I set the "revel" config variable "bugsnag.releasestage" to "staging"
  And I configure the bugsnag endpoint in the config file for 'revel'
  When I run the script "features/fixtures/revel/run.sh"
  And I wait for 4 seconds
  And I go to the route "/configure"
  And I wait for 1 seconds
  Then I should receive no requests

Scenario: A Revel session report contains the configured app version
  Given I set environment variable "DISABLE_REPORT_PAYLOADS" to "true"
  And I work with a new 'revel' app
  And I set the "revel" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I set the "revel" config variable "bugsnag.releasestage" to "development"
  And I configure the bugsnag sessions endpoint in the config file for 'revel'
  When I run the script "features/fixtures/revel/run.sh"
  And I wait for 4 seconds
  And I go to the route "/configure"
  And I wait for 1 seconds
  Then I should receive a request
  And the payload field "app.releaseStage" equals "development"
