Feature: Using recover

Scenario: An error report is sent when a goroutine crashes in a revel controller method
  Given I work with a new revel app
  And I set the revel config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint in the config file for revel
  When I run the script "features/fixtures/revel/run.sh"
  And I wait for 4 seconds
  And I go to the route "/recover"
  And I wait for 1 seconds
  Then I should receive a request
  And the exception "errorClass" equals "*errors.errorString"
  And the exception "message" equals "Go routine killed"
