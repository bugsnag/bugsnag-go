Feature: Using recover

Scenario: An error report is sent when a go routine crashes but recovers
  When I start the service "app"
  And I run "RecoverAfterPanicScenario"
  And I wait to receive an error
  And the exception "errorClass" equals "*errors.errorString"
  And the exception "message" equals "Go routine killed but recovered"