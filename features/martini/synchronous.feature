Feature: Configuring synchronous flag

Scenario: An error report is sent asynchrously but exits immediately so is not sent for martini
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag notify endpoint only
  When I run the script "features/fixtures/martini/run.sh"
  And I wait for 1 second
  And I send a request to '/async' on the martini app that might fail
  And I wait for 1 second
  Then I should receive no requests
