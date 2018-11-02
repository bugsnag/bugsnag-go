Feature: Configure auto capture sessions

Scenario: Disabling auto capture sessions for martini
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  When I run the script "features/fixtures/martini/run.sh"
  And I wait for 1 second
  And I go to the martini route "/handled"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
