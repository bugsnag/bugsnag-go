Feature: Configure integration with environment variables

    The library should be configurable using environment variables to support
    single-line and reusable configuration

    Background:
        Given I set environment variable "BUGSNAG_API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
        And I set environment variable "BUGSNAG_NOTIFY_ENDPOINT" to the notify endpoint
        And I set environment variable "BUGSNAG_SESSIONS_ENDPOINT" to the sessions endpoint
        And I have built the service "autoconfigure"

    Scenario Outline: Adding content to handled events through env variables
        Given I set environment variable "<variable>" to "<value>"
        And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
        When I run the go service "autoconfigure" with the test case "<testcase>"
        Then I wait to receive a request
        And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
        And the event "<field>" equals "<value>"

        Examples:
            | testcase | variable                              | value           | field                         |
            | panic    | BUGSNAG_APP_VERSION                   | 1.4.34          | app.version                   |
            | panic    | BUGSNAG_APP_TYPE                      | mailer-daemon   | app.type                      |
            | panic    | BUGSNAG_RELEASE_STAGE                 | beta1           | app.releaseStage              |
            | panic    | BUGSNAG_HOSTNAME                      | dream-machine-2 | device.hostname               |
            | panic    | BUGSNAG_METADATA_device_instance      | kube2-33-A      | metaData.device.instance      |
            | panic    | BUGSNAG_METADATA_framework_version    | v3.1.0          | metaData.framework.version    |
            | panic    | BUGSNAG_METADATA_device_runtime_level | 1C              | metaData.device.runtime_level |
            | panic    | BUGSNAG_METADATA_Carrot               | orange          | metaData.custom.Carrot        |

            | handled  | BUGSNAG_APP_VERSION                   | 1.4.34          | app.version                   |
            | handled  | BUGSNAG_APP_TYPE                      | mailer-daemon   | app.type                      |
            | handled  | BUGSNAG_RELEASE_STAGE                 | beta1           | app.releaseStage              |
            | handled  | BUGSNAG_HOSTNAME                      | dream-machine-2 | device.hostname               |
            | handled  | BUGSNAG_METADATA_device_instance      | kube2-33-A      | metaData.device.instance      |
            | handled  | BUGSNAG_METADATA_framework_version    | v3.1.0          | metaData.framework.version    |
            | handled  | BUGSNAG_METADATA_device_runtime_level | 1C              | metaData.device.runtime_level |
            | handled  | BUGSNAG_METADATA_Carrot               | orange          | metaData.custom.Carrot        |

    Scenario: Configuring project packages
        Given I set environment variable "BUGSNAG_PROJECT_PACKAGES" to "main,test"
        And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
        When I run the go service "autoconfigure" with the test case "panic"
        Then I wait to receive a request
        And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
        And the in-project frames of the stacktrace are:
            | file     | method        | lineNumber |
            | cases.go | explicitPanic | 22         |
            | main.go  | main          | 11         |

    Scenario: Configuring source root
        Given I set environment variable "BUGSNAG_SOURCE_ROOT" to the app directory
        And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
        And I run the go service "autoconfigure" with the test case "panic"
        Then I wait to receive a request
        And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
        And the in-project frames of the stacktrace are:
            | file     | method        | lineNumber |
            | cases.go | explicitPanic | 22         |
            | main.go  | main          | 11         |

    Scenario: Delivering events filtering through notify release stages
        Given I set environment variable "BUGSNAG_NOTIFY_RELEASE_STAGES" to "prod,beta"
        And I set environment variable "BUGSNAG_RELEASE_STAGE" to "beta"
        And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
        And I run the go service "autoconfigure" with the test case "panic"
        Then I wait to receive a request
        And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"

    Scenario: Suppressing events through notify release stages
        Given I set environment variable "BUGSNAG_NOTIFY_RELEASE_STAGES" to "prod,beta"
        And I set environment variable "BUGSNAG_RELEASE_STAGE" to "dev"
        And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
        And I run the go service "autoconfigure" with the test case "panic"
        Then 0 requests were received

    Scenario: Suppressing events using panic handler
        Given I set environment variable "BUGSNAG_DISABLE_PANIC_HANDLER" to "1"
        And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
        And I run the go service "autoconfigure" with the test case "panic"
        And I wait for 2 seconds
        Then 0 requests were received

    Scenario: Enabling synchronous event delivery
        Given I set environment variable "BUGSNAG_SYNCHRONOUS" to "1"
        And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
        When I run the go service "autoconfigure" with the test case "handled"
        Then 1 request was received

    Scenario: Filtering metadata
        Given I set environment variable "BUGSNAG_PARAMS_FILTERS" to "tomato,pears"
        And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
        When I run the go service "autoconfigure" with the test case "handled-metadata"
        Then I wait to receive a request
        And the event "metaData.fruit.Tomato" equals "[FILTERED]"
        And the event "metaData.snacks.Carrot" equals "4"
