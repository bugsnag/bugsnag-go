Feature: Configuring app version

Background:
  And I set environment variable "APP_VERSION" to "3.1.2"

Scenario: An error report contains the configured app type when running a go app
  Given I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  When I start the service "app"
  And I run HandledScenario
  And I wait to receive an error
  And the event "app.version" equals "3.1.2"

Scenario: A session report contains the configured app type when running a go app
  Given I set environment variable "AUTO_CAPTURE_SESSIONS" to "true"
  When I start the service "app"
  And I run SendSessionScenario
  And I wait to receive a session
  And the session payload field "app.version" equals "3.1.2"
