Feature: Configuring on before notify

Scenario: Send three bugsnags and use on before notify to drop one and modify the message of another
  When I start the service "app"
  And I run "OnBeforeNotifyScenario"
  And I wait to receive 2 errors
  And the exception "message" equals "Don't ignore this error"
  And I discard the oldest error
  And the exception "message" equals "Error message was changed"
