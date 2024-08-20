Feature: Configuring release stage

Background:
  Given I set environment variable "BUGSNAG_RELEASE_STAGE" to "my-stage"

Scenario: An error report is sent with configured release stage
  Given I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "false"
  When I start the service "app"
  And I run "HttpServerScenario"
  And I wait for the host "localhost" to open port "4512"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4512/handled"
  Then I wait to receive an error
  And the event "app.releaseStage" equals "my-stage"

Scenario: A session report contains the configured app type
  Given I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "true"
  When I start the service "app"
  And I run "HttpServerScenario"
  And I wait for the host "localhost" to open port "4512"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4512/session"
  Then I wait to receive a session
  And I wait to receive an error
  And the session payload field "sessions.0.app.releaseStage" equals "my-stage"