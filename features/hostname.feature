Feature: Configuring hostname

Scenario: A negroni error report contains the configured hostname
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I am working with a new negroni app
  And I configure the bugsnag notify endpoint only
  And I set environment variable "HOSTNAME" to "server-24"
  When I run the script "features/fixtures/negroni/run.sh"
  And I go to the negroni route "/handled"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "device.hostname" equals "server-24"

Scenario: A martini error report contains the configured hostname
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I am working with a new martini app
  And I configure the bugsnag notify endpoint only
  And I set environment variable "HOSTNAME" to "server-24"
  When I run the script "features/fixtures/martini/run.sh"
  And I go to the martini route "/handled"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "device.hostname" equals "server-24"

Scenario: An error report contains the configured hostname
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  And I set environment variable "HOSTNAME" to "server-1a"
  When I configure with the "hostname" configuration and send an error
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "device.hostname" equals "server-1a"

Scenario: An session report contains the configured hostname
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  And I set environment variable "HOSTNAME" to "server-1a"
  When I configure with the "hostname" configuration and send a session
  And I wait for 1 second
  Then I should receive a request
  And the "bugsnag-api-key" header equals "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the "bugsnag-payload-version" header equals "1.0"
  And the payload field "device.hostname" equals "server-1a"