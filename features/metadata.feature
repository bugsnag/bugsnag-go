Feature: Sending meta data

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  And I configure the bugsnag endpoint
  And I have built the service "app"

Scenario: An error report contains custom meta data
  When I run the go service "app" with the test case "metadata"
  And I wait to receive a request
  And the request is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "metaData.Scheme.Customer.Name" equals "Joe Bloggs"
  And the event "metaData.Scheme.Customer.Age" equals "21"
  And the event "metaData.Scheme.Level" equals "Blue"
