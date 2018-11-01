Feature: Configuring synchronous flag

Scenario: An error report is sent asynchrously but exits immediately so is not sent
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  And I set environment variable "SYNCHRONOUS" to "false"
  When I configure with the "synchronous" configuration and send an error
  And I wait for 1 second
  Then I should receive no requests

Scenario: An error report is report synchronously so it will send before exiting
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  And I set environment variable "SYNCHRONOUS" to "true"
  When I configure with the "synchronous" configuration and send an error
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"

Scenario: Revel configuration through code
    Given I set environment variable "USE_CODE_CONFIG" to "YES"
    Given I set environment variable "SYNCHRONOUS" to "true"
    And I work with a new 'revel-0.20.0' app
    And I set the "revel-0.20.0" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
    And I configure the bugsnag endpoint in the config file for 'revel-0.20.0'
    When I run the script "features/fixtures/revel-0.20.0/run.sh"
    And I wait for 4 seconds
    And I go to the route "/synchronous"
    And I wait for 1 seconds
    Then I should receive a request

Scenario: Revel configuration through config file
    And I work with a new 'revel-0.20.0' app
    And I set the "revel-0.20.0" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
    And I set the "revel-0.20.0" config variable "bugsnag.synchronous" to "true"
    And I configure the bugsnag endpoint in the config file for 'revel-0.20.0'
    When I run the script "features/fixtures/revel-0.20.0/run.sh"
    And I wait for 4 seconds
    And I go to the route "/synchronous"
    And I wait for 1 seconds
    Then I should receive a request
