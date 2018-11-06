Feature: Configuring app type

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "APP_VERSION" to "3.1.2"
  And I have built the service "gin"
  And I stop the service "gin"
  And I set environment variable "SERVER_PORT" to "4511"

Scenario Outline: A error report contains the configured app type when using a gin app
  Given I set environment variable "GIN_VERSION" to "<gin version>"
  And I set environment variable "GO_VERSION" to "<go version>"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  When I start the service "gin"
  And I wait for the app to open port "4511"
  And I wait for 1 seconds
  And I open the URL "http://localhost:4511/handled"
  And I wait for 1 seconds
  Then I should receive a request
  And the request is valid for the error reporting API
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "app.version" equals "3.1.2"

  Examples:
  | go version | gin version |
  | 1.11       | v1.3.0      |
  # | 1.10       | v1.3.0      |
  # | 1.9        | v1.3.0      |
  # | 1.8        | v1.3.0      |
  # | 1.7        | v1.3.0      |
  # | 1.11       | v1.2        |
  # | 1.10       | v1.2        |
  # | 1.9        | v1.2        |
  # | 1.8        | v1.2        |
  # | 1.7        | v1.2        |
  # | 1.11       | v1.1        |
  # | 1.10       | v1.1        |
  # | 1.9        | v1.1        |
  # | 1.8        | v1.1        |
  # | 1.7        | v1.1        |
  # | 1.11       | v1.0        |
  # | 1.10       | v1.0        |
  # | 1.9        | v1.0        |
  # | 1.8        | v1.0        |
  # | 1.7        | v1.0        |

Scenario Outline: A session report contains the configured app type when using a gin app
  Given I set environment variable "GIN_VERSION" to "<gin version>"
  And I set environment variable "GO_VERSION" to "<go version>"
  When I start the service "gin"
  And I wait for the app to open port "4511"
  And I wait for 1 seconds
  And I open the URL "http://localhost:4511/session"
  And I wait for 1 seconds
  Then I should receive a request
  And the request is valid for the session tracking API
  And the session contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the payload field "app.version" equals "3.1.2"

  Examples:
  | go version | gin version |
  | 1.11       | v1.3.0      |
  # | 1.10       | v1.3.0      |
  # | 1.9        | v1.3.0      |
  # | 1.8        | v1.3.0      |
  # | 1.7        | v1.3.0      |
  # | 1.11       | v1.2        |
  # | 1.10       | v1.2        |
  # | 1.9        | v1.2        |
  # | 1.8        | v1.2        |
  # | 1.7        | v1.2        |
  # | 1.11       | v1.1        |
  # | 1.10       | v1.1        |
  # | 1.9        | v1.1        |
  # | 1.8        | v1.1        |
  # | 1.7        | v1.1        |
  # | 1.11       | v1.0        |
  # | 1.10       | v1.0        |
  # | 1.9        | v1.0        |
  # | 1.8        | v1.0        |
  # | 1.7        | v1.0        |

