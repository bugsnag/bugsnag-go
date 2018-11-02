Feature: Configuring endpoint

Scenario: An error report is sent when a negroni is configured to use the legacy endpoint
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I set the legacy endpoint only
  And I am working with a new negroni app
  When I run the script "features/fixtures/negroni/run.sh"
  And I go to the negroni route "/handled"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
