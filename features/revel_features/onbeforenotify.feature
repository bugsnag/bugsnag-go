Feature: Configuring on before notify

Scenario: Revel integration can send three bugsnags and use on before notify to drop one and modify the message of another
  And I work with a new revel app
  And I set the revel config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint in the config file for revel
  When I run the script "features/fixtures/revel/run.sh"
  And I wait for 4 seconds
  And I go to the route "/onbeforenotify"
  Then I should receive 2 requests
  And the payload field "events.0.exceptions.0.message" equals "Don't ignore this error" for request 0
  And the payload field "events.0.exceptions.0.message" equals "Error message was changed" for request 1
