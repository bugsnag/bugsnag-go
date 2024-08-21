Feature: Configuring hostname

Scenario: An error report contains the configured hostname
  Given I set environment variable "BUGSNAG_HOSTNAME" to "server-1a"
  And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
  When I start the service "app"
  And I run "HandledErrorScenario"
  And I wait to receive an error
  And the event "device.hostname" equals "server-1a"

Scenario: A session report contains the configured hostname
  Given I set environment variable "BUGSNAG_HOSTNAME" to "server-1a"
  And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "1"
  When I start the service "app"
  And I run "SendSessionScenario"
  And I wait to receive 2 sessions
  And the session payload field "device.hostname" equals "server-1a"