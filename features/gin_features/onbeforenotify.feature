Feature: Configuring on before notify

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I have built the service "gin"
  And I stop the service "gin"
  And I set environment variable "SERVER_PORT" to "4511"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"

Scenario Outline: Send three bugsnags and use on before notify to drop one and modify the message of another
  Given I set environment variable "GIN_VERSION" to "<gin version>"
  And I set environment variable "GO_VERSION" to "<go version>"
  When I start the service "gin"
  And I wait for the app to open port "4511"
  And I wait for 1 seconds
  And I open the URL "http://localhost:4511/onbeforenotify"
  And I wait for 1 seconds
  Then I should receive 2 requests
  
  And the request 0 is valid for the error reporting API
  And the request 0 contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the exception "message" equals "Don't ignore this error" for request 0
  
  And the request 1 is valid for the error reporting API
  And the request 1 contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the exception "message" equals "Error message was changed" for request 1

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
  