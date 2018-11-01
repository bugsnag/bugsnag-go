Feature: Configure auto capture sessions

Scenario: Disabling auto capture sessions for negroni
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I am working with a new negroni app
  And I configure the bugsnag endpoints
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  When I run the script "features/fixtures/negroni/run.sh"
  And I go to the negroni route "/handled"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"

Scenario: Disabling auto capture sessions for martini
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I am working with a new martini app
  And I configure the bugsnag endpoints
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  When I run the script "features/fixtures/martini/run.sh"
  And I go to the martini route "/handled"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"

Scenario: An error report and session report is sent when auto capture sessions is on using the net http framework
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "true"
  When I run the http-net test server with the "auto capture sessions" configuration
  And I wait for 1 second
  Then I should receive 2 requests
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "session" is not null
  And the payload field "events.0.session.events.handled" equals 1 for request 0
  And the payload field "sessionCounts.0.sessionsStarted" equals 1 for request 1

Scenario: Only an error report is sent when auto capture sessions is off using the net http framework
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  When I run the http-net test server with the "auto capture sessions" configuration
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "session" is null
