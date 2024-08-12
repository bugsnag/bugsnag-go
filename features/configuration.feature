Feature: Configure integration with environment variables

    The library should be configurable using environment variables to support
    single-line and reusable configuration

    Background:
        Given I set environment variable "BUGSNAG_API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"

    Scenario Outline: Adding content to handled events through env variables
        Given I set environment variable "<variable>" to "<value>"
        And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
        When I start the service "app"
        And I run "<testcase>"
        And I wait to receive an error
        And the event "<field>" equals "<value>"

        Examples:
            | testcase | variable                              | value           | field                         |
            | AutoconfigPanicScenario    | BUGSNAG_APP_VERSION                   | 1.4.34          | app.version                   |
            | AutoconfigPanicScenario    | BUGSNAG_APP_TYPE                      | mailer-daemon   | app.type                      |
            | AutoconfigPanicScenario    | BUGSNAG_RELEASE_STAGE                 | beta1           | app.releaseStage              |
            | AutoconfigPanicScenario    | BUGSNAG_HOSTNAME                      | dream-machine-2 | device.hostname               |
            | AutoconfigPanicScenario    | BUGSNAG_METADATA_device_instance      | kube2-33-A      | metaData.device.instance      |
            | AutoconfigPanicScenario    | BUGSNAG_METADATA_framework_version    | v3.1.0          | metaData.framework.version    |
            | AutoconfigPanicScenario    | BUGSNAG_METADATA_device_runtime_level | 1C              | metaData.device.runtime_level |
            | AutoconfigPanicScenario    | BUGSNAG_METADATA_Carrot               | orange          | metaData.custom.Carrot        |

            | AutoconfigHandledScenario  | BUGSNAG_APP_VERSION                   | 1.4.34          | app.version                   |
            | AutoconfigHandledScenario  | BUGSNAG_APP_TYPE                      | mailer-daemon   | app.type                      |
            | AutoconfigHandledScenario  | BUGSNAG_RELEASE_STAGE                 | beta1           | app.releaseStage              |
            | AutoconfigHandledScenario  | BUGSNAG_HOSTNAME                      | dream-machine-2 | device.hostname               |
            | AutoconfigHandledScenario  | BUGSNAG_METADATA_device_instance      | kube2-33-A      | metaData.device.instance      |
            | AutoconfigHandledScenario  | BUGSNAG_METADATA_framework_version    | v3.1.0          | metaData.framework.version    |
            | AutoconfigHandledScenario  | BUGSNAG_METADATA_device_runtime_level | 1C              | metaData.device.runtime_level |
            | AutoconfigHandledScenario  | BUGSNAG_METADATA_Carrot               | orange          | metaData.custom.Carrot        |

    Scenario: Configuring project packages
        Given I set environment variable "BUGSNAG_PROJECT_PACKAGES" to "main,test"
        And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
        When I start the service "app"
        And I run "AutoconfigPanicScenario"
        And I wait to receive an error
        And the in-project frames of the stacktrace are:
            | file     | method        | lineNumber |
            | cases.go | explicitPanic | 22         |
            | main.go  | main          | 11         |

    Scenario: Configuring source root
        Given I set environment variable "BUGSNAG_SOURCE_ROOT" to "/app/src/features/fixtures/app/"
        And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
        When I start the service "app"
        And I run "AutoconfigPanicScenario"
        And I wait to receive an error
        And the in-project frames of the stacktrace are:
            | file     | method        | lineNumber |
            | cases.go | explicitPanic | 22         |
            | main.go  | main          | 11         |

    Scenario: Delivering events filtering through notify release stages
        Given I set environment variable "BUGSNAG_NOTIFY_RELEASE_STAGES" to "prod,beta"
        And I set environment variable "BUGSNAG_RELEASE_STAGE" to "beta"
        And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
        When I start the service "app"
        And I run "AutoconfigPanicScenario"
        And I wait to receive an error

    Scenario: Suppressing events through notify release stages
        Given I set environment variable "BUGSNAG_NOTIFY_RELEASE_STAGES" to "prod,beta"
        And I set environment variable "BUGSNAG_RELEASE_STAGE" to "dev"
        And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
        When I start the service "app"
        And I run "AutoconfigPanicScenario"
        Then I should receive no errors

    Scenario: Suppressing events using panic handler
        Given I set environment variable "BUGSNAG_DISABLE_PANIC_HANDLER" to "1"
        And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
        When I start the service "app"
        And I run "AutoconfigPanicScenario"
        And I wait for 2 seconds
        Then I should receive no errors

    Scenario: Enabling synchronous event delivery
        Given I set environment variable "BUGSNAG_SYNCHRONOUS" to "1"
        And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
        When I start the service "app"
        And I run "AutoconfigHandledScenario"
        And I wait to receive an error

    Scenario: Filtering metadata
        Given I set environment variable "BUGSNAG_PARAMS_FILTERS" to "tomato,pears"
        And I set environment variable "BUGSNAG_AUTO_CAPTURE_SESSIONS" to "0"
        When I start the service "app"
        And I run "AutoconfigMetadataScenario"
        And I wait to receive an error
        And the event "metaData.fruit.Tomato" equals "[FILTERED]"
        And the event "metaData.snacks.Carrot" equals "4"
