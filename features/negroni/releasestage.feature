Feature: Configuring release stages and notify release stages

Scenario: An error report is sent when release stage matches notify release stages for negroni
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag notify endpoint only
  And I set environment variable "RELEASE_STAGE" to "staging"
  And I set environment variable "NOTIFY_RELEASE_STAGES" to "staging"
  When I run the script "features/fixtures/negroni/run.sh"
  And I go to the negroni route "/handled"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"

Scenario: An error report is sent when no notify release stages are specified for negroni
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag notify endpoint only
  And I set environment variable "RELEASE_STAGE" to "staging"
  When I run the script "features/fixtures/negroni/run.sh"
  And I go to the negroni route "/handled"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"

Scenario: An error report is sent regardless of notify release stages if release stage is not set for negroni
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag notify endpoint only
  And I set environment variable "NOTIFY_RELEASE_STAGES" to "staging"
  When I run the script "features/fixtures/negroni/run.sh"
  And I go to the negroni route "/handled"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"

Scenario: An error report is not sent if the release stage does not match the notify release stages for negroni
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag notify endpoint only
  And I set environment variable "RELEASE_STAGE" to "staging"
  And I set environment variable "NOTIFY_RELEASE_STAGES" to "production"
  When I run the script "features/fixtures/negroni/run.sh"
  And I go to the negroni route "/handled"
  And I wait for 1 second
  Then I should receive no requests
