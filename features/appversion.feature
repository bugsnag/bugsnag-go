Feature: Configuring app version

Scenario: An error report contains the configured app version
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  And I set environment variable "APP_VERSION" to "1.2.3"
  When I configure with the "app version" configuration and send an error
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "app.version" equals "1.2.3"

Scenario: An session report contains the configured app version
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  And I set environment variable "APP_VERSION" to "1.2.3"
  When I configure with the "app version" configuration and send a session
  And I wait for 1 second
  Then I should receive a request
  And the "bugsnag-api-key" header equals "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the "bugsnag-payload-version" header equals "1.0"
  And the payload field "app.version" equals "1.2.3"