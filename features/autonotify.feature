Feature: Using auto notify

Scenario: An error report is sent when an AutoNotified crash occurs which later gets recovered
  When I start the service "app"
  And I run "AutonotifyPanicScenario"
  And I wait to receive an error
  And the exception "errorClass" equals "*errors.errorString"
  And the exception "message" equals "Go routine killed with auto notify"
