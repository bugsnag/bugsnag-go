Feature: Configuring hostname

Scenario: An error report contains the configured hostname
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  And I set environment variable "HOSTNAME" to "server-1a"
  When I configure with the "hostname" configuration and send an error
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "device.hostname" equals "server-1a"

Scenario: An session report contains the configured hostname
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  And I set environment variable "HOSTNAME" to "server-1a"
  When I configure with the "hostname" configuration and send a session
  And I wait for 1 second
  Then I should receive a request
  And the "bugsnag-api-key" header equals "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the "bugsnag-payload-version" header equals "1.0"
  And the payload field "device.hostname" equals "server-1a"

Scenario: Revel configuration through code
    Given I set environment variable "USE_CODE_CONFIG" to "YES"
    Given I set environment variable "HOSTNAME" to "super-server-1"
    And I work with a new 'revel-0.20.0' app
    And I set the "revel-0.20.0" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
    And I configure the bugsnag endpoint in the config file for 'revel-0.20.0'
    When I run the script "features/fixtures/revel-0.20.0/run.sh"
    And I wait for 4 seconds
    And I go to the route "/configure"
    And I wait for 1 seconds
    Then I should receive a request
    And the event "device.hostname" equals "super-server-1"

Scenario: Revel configuration through config file
    And I work with a new 'revel-0.20.0' app
    And I set the "revel-0.20.0" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
    And I set the "revel-0.20.0" config variable "bugsnag.hostname" to "super-server-2"
    And I configure the bugsnag endpoint in the config file for 'revel-0.20.0'
    When I run the script "features/fixtures/revel-0.20.0/run.sh"
    And I wait for 4 seconds
    And I go to the route "/configure"
    And I wait for 1 seconds
    Then I should receive a request
    And the event "device.hostname" equals "super-server-2"
