Feature: Configuring app version

Background:
  And I set environment variable "BUGSNAG_APP_VERSION" to "3.1.2"

Scenario: A error report contains the configured app type when using a net http app
  Given I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "false"
  When I start the service "app"
  And I run "HttpServerScenario"
  And I wait for the host "localhost" to open port "4512"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4512/handled"
  And I wait to receive an error
  And I should receive no sessions
  And the error is valid for the error reporting API version "4" for the "Bugsnag Go" notifier
  And the event "app.version" equals "3.1.2"

Scenario: A session report contains the configured app type when using a net http app
  Given I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "true"
  When I start the service "app"
  And I run "HttpServerScenario"
  And I wait for the host "localhost" to open port "4512"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4512/session"
  And I wait to receive a session
  And the session is valid for the session reporting API version "1.0" for the "Bugsnag Go" notifier
  And the session payload field "sessions.0.app.version" equals "3.1.2"

