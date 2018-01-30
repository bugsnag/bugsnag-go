Feature: Sinatra support

Scenario: Sinatra unhandled exception
    When I set environment variable "BUGSNAG_API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
    And I start a Sinatra app
    And I navigate to the route "/syntax-error"
    Then I should receive a request
    And the request is a valid for the error reporting API
    And the request used the Ruby notifier
    And the "Bugsnag-API-Key" header equals "a35a2a72bd230ac0aa0f52715bbdc6aa"
    And the payload field "events" is an array with 1 element
    And the event "unhandled" is true
    And the event "app.version" equals "2.5.1"
    And the event "context" equals "GET /syntax-error"
    And the exception "errorClass" equals "NoMethodError"
    And the exception "message" starts with "undefined method `rt' for #<Sinatra::Application"
    And the "method" of stack frame 0 equals "make_a_syntax_error"
    And the "method" of stack frame 1 equals "block in <main>"

Scenario: Sinatra handled exception
    When I set environment variable "BUGSNAG_API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
    And I start a Sinatra app
    And I navigate to the route "/notify"
    Then I should receive a request
    And the request is a valid for the error reporting API
    And the request used the Ruby notifier
    And the "Bugsnag-API-Key" header equals "a35a2a72bd230ac0aa0f52715bbdc6aa"
    And the event "context" equals "GET /notify"
    And the event "unhandled" is false
    And the exception "errorClass" equals "InvariantException"
    And the exception "message" starts with "The cake was rotten"
    And the "method" of stack frame 0 equals "send_manual_notify"
    And the "method" of stack frame 1 equals "block in <main>"

Scenario: Sinatra override context
    When I set environment variable "BUGSNAG_API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
    And I start a Sinatra app
    And I navigate to the route "/notify?context=foo"
    Then I should receive a request
    And the event "context" equals "foo"
