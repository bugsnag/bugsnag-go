Feature: Configuring on before notify

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I set environment variable "SERVER_PORT" to "4514"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"

Scenario Outline: Send three bugsnags and use on before notify to drop one and modify the message of another
  Given I set environment variable "NEGRONI_VERSION" to "<negroni version>"
  When I start the service "negroni"
  And I wait for the app to open port "4514"
  And I wait for 2 seconds
  And I open the URL "http://localhost:4514/onbeforenotify"
  Then I wait to receive 2 requests
  
  And the request 0 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the exception "message" equals "Don't ignore this error" for request 0
  
  And the request 1 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the exception "message" equals "Error message was changed" for request 1

  Examples:
  | negroni version |
  | v1.0.0          |
  | v0.3.0          |
  | v0.2.0          |
  