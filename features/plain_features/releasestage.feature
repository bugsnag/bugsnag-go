Feature: Configuring release stages and notify release stages

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I have built the service "app"

Scenario: An error report is sent when release stage matches notify release stages
  Given I set environment variable "NOTIFY_RELEASE_STAGES" to "stage1,stage2,stage3"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  And I set environment variable "RELEASE_STAGE" to "stage2"
  When I run the go service "app" with the test case "handled"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "app.releaseStage" equals "stage2"

Scenario: An error report is sent when no notify release stages are specified
  Given I set environment variable "RELEASE_STAGE" to "stage2"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  When I run the go service "app" with the test case "handled"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "app.releaseStage" equals "stage2"

Scenario: An error report is sent regardless of notify release stages if release stage is not set
  Given I set environment variable "NOTIFY_RELEASE_STAGES" to "stage1,stage2,stage3"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  When I run the go service "app" with the test case "handled"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"

Scenario: An error report is not sent if the release stage does not match the notify release stages
  Given I set environment variable "NOTIFY_RELEASE_STAGES" to "stage1,stage2,stage3"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  And I set environment variable "RELEASE_STAGE" to "stage4"
  When I run the go service "app" with the test case "handled"
  And I wait for 3 second
  Then I should receive no requests

Scenario: An session report is sent when release stage matches notify release stages
  Given I set environment variable "NOTIFY_RELEASE_STAGES" to "stage1,stage2,stage3"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "true"
  And I set environment variable "RELEASE_STAGE" to "stage2"
  When I run the go service "app" with the test case "session"
  Then I wait to receive a request after the start up session
  And the request is a valid session report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the payload field "app.releaseStage" equals "stage2"

Scenario: An session report is sent when no notify release stages are specified
  Given I set environment variable "RELEASE_STAGE" to "stage2"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "true"
  When I run the go service "app" with the test case "session"
  Then I wait to receive a request after the start up session
  And the request is a valid session report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the payload field "app.releaseStage" equals "stage2"

Scenario: An session report is sent regardless of notify release stages if release stage is not set
  Given I set environment variable "NOTIFY_RELEASE_STAGES" to "stage1,stage2,stage3"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "true"
  When I run the go service "app" with the test case "session"
  Then I wait to receive a request after the start up session
  And the request is a valid session report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"

Scenario: An session report is not sent if the release stage does not match the notify release stages
  Given I set environment variable "NOTIFY_RELEASE_STAGES" to "stage1,stage2,stage3"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "true"
  And I set environment variable "RELEASE_STAGE" to "stage4"
  When I run the go service "app" with the test case "session"
  And I wait for 3 second
  Then I should receive no requests

