Feature: Custom ErrorClass and Message fields

Background:
  Given I set environment variable "BUGSNAG_SOURCE_ROOT" to "/app/src/features/fixtures/app/"
  And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"

Scenario: A handled error with custom ErrorClass
  When I start the service "app"
  And I run "CustomErrorClassScenario"
  And I wait to receive an error
  And the event "unhandled" is false
  And the event "severity" equals "warning"
  And the exception "errorClass" equals "CustomErrorClass"
  And the exception "message" matches ".*"

Scenario: A handled error with custom Message
  When I start the service "app"
  And I run "CustomMessageScenario"
  And I wait to receive an error
  And the event "unhandled" is false
  And the event "severity" equals "warning"
  And the exception "errorClass" matches "\*os.PathError|\*fs.PathError"
  And the exception "message" equals "This is a custom message"

Scenario: A handled error with both custom ErrorClass and Message
  When I start the service "app"
  And I run "CustomBothScenario"
  And I wait to receive an error
  And the event "unhandled" is false
  And the event "severity" equals "warning"
  And the exception "errorClass" equals "CustomBothError"
  And the exception "message" equals "Custom message with both fields"

Scenario: Multiple errors with different custom ErrorClass and Message combinations
  When I start the service "app"
  And I run "MultipleCustomFieldsScenario"
  And I wait to receive 2 errors
  And the exception "errorClass" equals "FirstCustomError"
  And the exception "message" equals "First error custom message"
  And I discard the oldest error
  And the exception "errorClass" equals "SecondCustomError"
  And the exception "message" equals "Second error custom message"
Scenario: Multiple ErrorClass objects - last one wins
  When I start the service "app"
  And I run "MultipleErrorClassScenario"
  And I wait to receive an error
  And the event "unhandled" is false
  And the event "severity" equals "warning"
  And the exception "errorClass" equals "FinalErrorClass"
  And the exception "message" matches ".*"

Scenario: Multiple Message objects - last one wins
  When I start the service "app"
  And I run "MultipleMessageScenario"
  And I wait to receive an error
  And the event "unhandled" is false
  And the event "severity" equals "warning"
  And the exception "errorClass" matches "\*os.PathError|\*fs.PathError"
  And the exception "message" equals "Final custom message"

Scenario: Multiple ErrorClass and Message objects mixed - last of each wins
  When I start the service "app"
  And I run "MultipleOfBothScenario"
  And I wait to receive an error
  And the event "unhandled" is false
  And the event "severity" equals "warning"
  And the exception "errorClass" equals "SecondErrorClass"
  And the exception "message" equals "Second message"