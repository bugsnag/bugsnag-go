Feature: Revel 0.20.0 configured through code

Scenario: Setting the app version throuh code is reflected in reports
    Given I set environment variable "USE_CODE_CONFIG" to "YES"
    Given I set environment variable "APP_VERSION" to "1.33.7"
    And I work with a new 'revel-0.20.0' app
    And I set the "revel-0.20.0" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
    And I configure the bugsnag endpoint in the config file for 'revel-0.20.0'
    When I run the script "features/fixtures/revel-0.20.0/run.sh"
    And I wait for 2 seconds
    And I go to the route "/configure"
    And I wait for 1 seconds
    Then I should receive a request
    And the event "app.version" equals "1.33.7"

