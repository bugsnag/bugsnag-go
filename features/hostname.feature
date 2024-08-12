Feature: Configuring hostname

Scenario: An error report contains the configured hostname
  Given I set environment variable "BUGSNAG_HOSTNAME" to "server-1a"
  And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "false"
  When I start the service "app"
  And I run "HandledScenario"
  And I wait to receive an error
  And the event "device.hostname" equals "server-1a"

Scenario: An session report contains the configured hostname
  Given I set environment variable "BUGSNAG_HOSTNAME" to "server-1a"
  And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "true"
  When I start the service "app"
  And I run "SendSessionScenario"
  And I wait to receive a session
  And the session payload field "device.hostname" equals "server-1a"
