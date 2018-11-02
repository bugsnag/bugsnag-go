Feature: Session data inside an error report using a session context

Scenario: Revel reports contains a session count for handled errors
  And I work with a new revel app
  And I set the revel config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint in the config file for revel
  When I run the script "features/fixtures/revel/run.sh"
  And I wait for 4 seconds
  And I go to the route "/handled"
  And I wait for 1 seconds
  Then I should receive a request
  And the event "session.events.handled" equals 1
  And the event "session.events.unhandled" equals 0

Scenario: Revel reports contains a session count for panics
  And I work with a new revel app
  And I set the revel config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint in the config file for revel
  When I run the script "features/fixtures/revel/run.sh"
  And I wait for 4 seconds
  And I go to the route "/unhandled"
  And I wait for 1 seconds
  Then I should receive a request
  And the event "session.events.handled" equals 0
  And the event "session.events.unhandled" equals 1
