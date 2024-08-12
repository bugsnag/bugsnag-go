Feature: Sending meta data

Scenario: An error report contains custom meta data
  When I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  And I start the service "app"
  And I run "MetadataScenario"
  And I wait to receive an error
  And the event "metaData.Scheme.Customer.Name" equals "Joe Bloggs"
  And the event "metaData.Scheme.Customer.Age" equals "21"
  And the event "metaData.Scheme.Level" equals "Blue"
