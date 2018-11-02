Feature: Configuring on before notify

Scenario: A negroni error report contains the configured hostname
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  Given I set environment variable "SYNCHRONOUS" to "true"
  And I configure the bugsnag notify endpoint only
  When I run the script "features/fixtures/negroni/run.sh"
  And I go to the negroni route "/onbeforenotify"
  Then I should receive 2 requests
  And the "bugsnag-api-key" header equals "a35a2a72bd230ac0aa0f52715bbdc6aa" for request 0
  And the "bugsnag-api-key" header equals "a35a2a72bd230ac0aa0f52715bbdc6aa" for request 1
  And the payload field "events.0.exceptions.0.message" equals "Don't ignore this error" for request 0
  And the payload field "events.0.exceptions.0.message" equals "Error message was changed" for request 1
