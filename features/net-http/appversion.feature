Feature: Configuring app version

Background:
  And I set environment variable "BUGSNAG_APP_VERSION" to "3.1.2"

Scenario: A error report contains the configured app type when using a net http app
  Given I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
  When I start the service "app"
  And I run "HttpServerScenario"
  And I wait for the host "localhost" to open port "4512"
  And I open the URL "http://localhost:4512/handled"
  And I wait to receive an error
  And I should receive no sessions
  And the event "app.version" equals "3.1.2"

Scenario: A session report contains the configured app type when using a net http app
  Given I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "1"
  When I start the service "app"
  And I run "HttpServerScenario"
  And I wait for the host "localhost" to open port "4512"
  And I open the URL "http://localhost:4512/session"
  And I wait to receive a session
  And the session payload field "app.version" equals "3.1.2"
