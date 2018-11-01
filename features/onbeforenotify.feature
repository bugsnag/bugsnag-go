Feature: Configuring on before notify

Scenario: Send three bugsnags and use on before notify to drop one and modify the message of another
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  When I configure with the "on before notify" configuration and send an error
  And I wait for 1 second
  Then I should receive 2 requests
  And the "bugsnag-api-key" header equals "a35a2a72bd230ac0aa0f52715bbdc6aa" for request 0
  And the "bugsnag-api-key" header equals "a35a2a72bd230ac0aa0f52715bbdc6aa" for request 1
  And the payload field "events.0.exceptions.0.message" equals "Don't ignore this error" for request 0
  And the payload field "events.0.exceptions.0.message" equals "Error message was changed" for request 1

Scenario: Revel integration can send three bugsnags and use on before notify to drop one and modify the message of another
  And I work with a new 'revel-0.20.0' app
  And I set the "revel-0.20.0" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint in the config file for 'revel-0.20.0'
  When I run the script "features/fixtures/revel-0.20.0/run.sh"
  And I wait for 4 seconds
  And I go to the route "/onbeforenotify"
  Then I should receive 2 requests
  And the payload field "events.0.exceptions.0.message" equals "Don't ignore this error" for request 0
  And the payload field "events.0.exceptions.0.message" equals "Error message was changed" for request 1

Scenario: Send three bugsnags and use on before notify to drop one and modify the message of another using net http
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  When I run the http-net test server with the "on before notify" configuration
  And I wait for 1 second
  Then I should receive 3 requests
  And the "bugsnag-api-key" header equals "a35a2a72bd230ac0aa0f52715bbdc6aa" for request 0 
  And the "bugsnag-api-key" header equals "a35a2a72bd230ac0aa0f52715bbdc6aa" for request 1
  And the payload field "events.0.exceptions.0.message" equals "Don't ignore this error" for request 0
  And the payload field "events.0.exceptions.0.message" equals "Error message was changed" for request 1
