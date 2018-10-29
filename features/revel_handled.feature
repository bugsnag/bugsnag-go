Feature: Revel 0.20.0 integration

Scenario: A handled error contains request data when given context
    Given I work with a new 'revel-0.20.0' app
    And I set the "revel-0.20.0" config variable "bugsnag.apikey" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
    And I configure the bugsnag endpoint in the config file for 'revel-0.20.0'
    When I run the script "features/fixtures/revel-0.20.0/run.sh"
    And I wait for 2 seconds
    And I go to the route "/handled"
    And I wait for 1 seconds
    Then I should receive a request
    And the request used payload v4 headers
    And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
    And the event "app.releaseStage" equals "dev"
    And the event "app.type" equals "Revel"
    And the event "context" equals "/handled"
    And the event "request.httpMethod" equals "GET"
    And the event "request.headers.X-Forwarded-For" equals "::1"
    And the event "request.url" equals "http://localhost:9020/handled"
    And the event "session.events.handled" equals 1
    And the event "session.events.unhandled" equals 0
    And the event "unhandled" is false
    And the event "severity" equals "warning"
    And the event "severityReason.type" equals "handledError"
    And the exception "errorClass" equals "*errors.errorString"
    And the "file" of stack frame 0 equals "controllers/app.go"
    And the "lineNumber" of stack frame 0 equals 19
