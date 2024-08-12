Feature: Configuring app type

Background:
  Given I set environment variable "APP_TYPE" to "background-queue"

Scenario: An error report contains the configured app type when running a go app
  Given I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  When I start the service "app"
  And I run "HandledScenario"
  And I wait to receive an error
  And the event "app.type" equals "background-queue"

Scenario: An session report contains the configured app type when running a go app
  Given I set environment variable "AUTO_CAPTURE_SESSIONS" to "true"
  When I start the service "app"
  And I run "SendSessionScenario"
  And I wait to receive a session
  And the session payload field "app.type" equals "background-queue"
