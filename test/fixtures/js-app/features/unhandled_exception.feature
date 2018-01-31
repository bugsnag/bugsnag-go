Feature: Browser support

Background:
    When I set environment variable "BUGSNAG_API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"

Scenario Outline: ReferenceError - Undefined is not a function
    When I navigate to the "undefined-is-not-a-function" page in "<browser>"
    Then I should receive a request
    And the request is a valid for the error reporting API
    And the payload field "events" is an array with 1 element
    And the event "unhandled" is true
    And the event "context" equals "/undefined-is-not-a-function.html"
    And the script content is in event metadata
    And the exception "errorClass" equals "ReferenceError"
    And the exception "message" equals "<message>"
    And the "method" of stack frame 0 equals "start_reactor"
    And the "lineNumber" of stack frame 0 equals 22
    And the "columnNumber" of stack frame 0 equals <colNum>

    Examples:
        | browser | colNum | message                                     |
        | Chrome  | 17     | initialize_cooling_bay is not defined       |
        | Safari  | 39     | Can't find variable: initialize_cooling_bay |
        | Firefox | 17     | initialize_cooling_bay is not defined       |

Scenario Outline: Rejecting a promise with an error
    When I navigate to the "unhandled-promise-rejection" page in "<browser>"
    Then I should receive a request
    And the request is a valid for the error reporting API
    And the payload field "events" is an array with 1 element
    And the event "unhandled" is true
    And the event "context" equals "/unhandled-promise-rejection.html"
    And the script content is in event metadata
    And the exception "errorClass" equals "Error"
    And the exception "message" equals "There is no cake"
    And the "method" of stack frame 0 equals "start_reactor"
    And the "lineNumber" of stack frame 0 equals 24
    And the "columnNumber" of stack frame 0 equals <colNum>

    Examples:
        | browser | colNum |
        | Chrome  | 32     |
        | Safari  | 41     |

Scenario Outline: Rejecting a promise with a string
    When I navigate to the "unhandled-promise-rejection-string" page in "<browser>"
    Then I should receive a request
    And the request is a valid for the error reporting API
    And the payload field "events" is an array with 1 element
    And the event "unhandled" is true
    And the event "context" equals "/unhandled-promise-rejection-string.html"
    And the exception "errorClass" equals "UnhandledRejection"
    And the exception "message" equals "Rejection reason was not an Error. See "Promise" tab for more detail."
    And the event "metaData.promise.rejection reason" equals "There is no cake"
    # And the script content is in event metadata
    # And the "method" of stack frame 0 equals "start_reactor"
    # And the "columnNumber" of stack frame 0 equals <colNum>
    # And the "lineNumber" of stack frame 0 equals 24

    Examples:
        | browser | colNum |
        | Chrome  | 32     |
        | Safari  | 41     |
