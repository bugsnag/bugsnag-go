Feature: Configuring notify release stages

Scenario: An error report is sent when release stage matches notify release stages
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "NOTIFY_RELEASE_STAGES" to "stage1,stage2,stage3"
  And I set environment variable "RELEASE_STAGE" to "stage2"
  When I configure with the "release stage" configuration and send an error
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "app.releaseStage" equals "stage2"

Scenario: An error report is sent when no notify release stages are specified
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "RELEASE_STAGE" to "stage2"
  When I configure with the "release stage" configuration and send an error
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "app.releaseStage" equals "stage2"
  
Scenario: An error report is sent regardless of notify release stages if release stage is not set
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "NOTIFY_RELEASE_STAGES" to "stage1,stage2,stage3"
  When I configure with the "release stage" configuration and send an error
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "app.releaseStage" equals ""

Scenario: An error report is not if the release stage doesn't match the notify release stages
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "NOTIFY_RELEASE_STAGES" to "stage1,stage2,stage3"
    And I set environment variable "RELEASE_STAGE" to "stage4"
  When I configure with the "release stage" configuration and send an error
  And I wait for 1 second
  Then I should receive no requests

