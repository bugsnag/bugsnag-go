Feature: Configuring release stages and notify release stages

Scenario: An error report is sent when release stage matches notify release stages
  Given I set environment variable "BUGSNAG_NOTIFY_RELEASE_STAGES" to "stage1,stage2,stage3"
  And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
  And I set environment variable "BUGSNAG_RELEASE_STAGE" to "stage2"
  When I start the service "app"
  And I run "HandledErrorScenario"
  And I wait to receive an error
  And the event "app.releaseStage" equals "stage2"

Scenario: An error report is sent when no notify release stages are specified
  Given I set environment variable "BUGSNAG_RELEASE_STAGE" to "stage2"
  And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
  When I start the service "app"
  And I run "HandledErrorScenario"
  And I wait to receive an error
  And the event "app.releaseStage" equals "stage2"

Scenario: An error report is sent regardless of notify release stages if release stage is not set
  Given I set environment variable "BUGSNAG_NOTIFY_RELEASE_STAGES" to "stage1,stage2,stage3"
  And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
  When I start the service "app"
  And I run "HandledErrorScenario"
  And I wait to receive an error

Scenario: An error report is not sent if the release stage does not match the notify release stages
  Given I set environment variable "BUGSNAG_NOTIFY_RELEASE_STAGES" to "stage1,stage2,stage3"
  And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
  And I set environment variable "BUGSNAG_RELEASE_STAGE" to "stage4"
  When I start the service "app"
  And I run "HandledErrorScenario"
  And I should receive no errors

Scenario: A session report is sent when release stage matches notify release stages
  Given I set environment variable "BUGSNAG_NOTIFY_RELEASE_STAGES" to "stage1,stage2,stage3"
  And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "1"
  And I set environment variable "BUGSNAG_RELEASE_STAGE" to "stage2"
  When I start the service "app"
  And I run "SendSessionScenario"
  And I wait to receive 2 sessions
  And the session payload field "app.releaseStage" equals "stage2"

Scenario: A session report is sent when no notify release stages are specified
  Given I set environment variable "BUGSNAG_RELEASE_STAGE" to "stage2"
  And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "1"
  When I start the service "app"
  And I run "SendSessionScenario"
  And I wait to receive 2 sessions
  And the session payload field "app.releaseStage" equals "stage2"

Scenario: A session report is sent regardless of notify release stages if release stage is not set
  Given I set environment variable "BUGSNAG_NOTIFY_RELEASE_STAGES" to "stage1,stage2,stage3"
  And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "1"
  When I start the service "app"
  And I run "SendSessionScenario"
  And I wait to receive 2 sessions

Scenario: A session report is not sent if the release stage does not match the notify release stages
  Given I set environment variable "BUGSNAG_NOTIFY_RELEASE_STAGES" to "stage1,stage2,stage3"
  And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "1"
  And I set environment variable "BUGSNAG_RELEASE_STAGE" to "stage4"
  When I start the service "app"
  And I run "SendSessionScenario"
  And I should receive no sessions
