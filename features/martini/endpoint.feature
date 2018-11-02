Feature: Configuring endpoint

Scenario: An error report is sent when a martini is configured to use the legacy endpoint
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I set the legacy endpoint only
  When I run the script "features/fixtures/martini/run.sh"
  And I go to the martini route "/handled"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
