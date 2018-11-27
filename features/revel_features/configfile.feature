Feature: Configuring using the config file

Background:
  Given I configure the bugsnag endpoint
  And I set environment variable "USE_PROPERTIES_FILE_CONFIG" to "true"

Scenario: A error report contains the variables set in the config file
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
