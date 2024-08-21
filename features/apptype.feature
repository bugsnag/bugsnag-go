Feature: Configuring app type

Background:
  Given I set environment variable "BUGSNAG_APP_TYPE" to "background-queue"

Scenario: An error report contains the configured app type when running a go app
  Given I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
  When I start the service "app"
  And I run "HandledErrorScenario"
  And I wait to receive an error
  And the event "app.type" equals "background-queue"

Scenario: An session report contains the configured app type when running a go app
  Given I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "1"
  When I start the service "app"
  And I run "SendSessionScenario"
  And I wait to receive 2 sessions
  And the session payload field "app.type" equals "background-queue"
