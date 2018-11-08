Feature: Configuring using the config file

Background:
  Given I configure the bugsnag endpoint

Scenario Outline: A error report contains the variables set in the config file
  Given I set environment variable "REVEL_VERSION" to "<revel version>"
  And I set environment variable "REVEL_CMD_VERSION" to "<revel cmd>"
  When I start the service "revel"
  And I wait for the app to open port "4515"
  And I wait for 4 seconds
  And I open the URL "http://localhost:4515/handled"
  Then I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "app.version" equals "4.5.6"
  And the event "app.type" equals "config-file-app-type"
  And the event "app.releaseStage" equals "config-test"
  And the event "device.hostname" equals "config-file-server"

  Examples:
  | revel version | revel cmd |
  | v0.21.0       | v0.21.1   |
  | v0.20.0       | v0.20.2   |
  | v0.19.1       | v0.19.0   |
  | v0.18.0       | v0.18.0   |