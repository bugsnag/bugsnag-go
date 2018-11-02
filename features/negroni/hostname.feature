Feature: Configuring hostname

Scenario: A negroni error report contains the configured hostname
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag notify endpoint only
  And I set environment variable "HOSTNAME" to "server-24"
  When I run the script "features/fixtures/negroni/run.sh"
  And I go to the negroni route "/handled"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "device.hostname" equals "server-24"
