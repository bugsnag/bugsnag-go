Feature: Sending meta data

Scenario: An error report can add metadata for negroni
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I am working with a new negroni app
  And I configure the bugsnag notify endpoint only
  When I run the script "features/fixtures/negroni/run.sh"
  And I go to the negroni route "/metadata"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "metaData.Scheme.Customer.Name" equals "Joe Bloggs"
  And the event "metaData.Scheme.Customer.Age" equals "21"
  And the event "metaData.Scheme.Level" equals "Blue"

Scenario: An error report can add metadata for martini
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I am working with a new martini app
  And I configure the bugsnag notify endpoint only
  When I run the script "features/fixtures/martini/run.sh"
  And I go to the martini route "/metadata"
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "metaData.Scheme.Customer.Name" equals "Joe Bloggs"
  And the event "metaData.Scheme.Customer.Age" equals "21"
  And the event "metaData.Scheme.Level" equals "Blue"

Scenario: An error report contains custom meta data
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  When I configure with the "metadata" configuration and send an error
  And I wait for 1 second
  Then I should receive a request
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "metaData.Scheme.Customer.Name" equals "Joe Bloggs"
  And the event "metaData.Scheme.Customer.Age" equals "21"
  And the event "metaData.Scheme.Level" equals "Blue"

Scenario: An error report contains custom meta data when using the net http framework
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoints
  When I run the http-net test server with the "meta data" configuration
  And I wait for 1 second
  Then I should receive 2 requests
  And the request used payload v4 headers
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "metaData.Scheme.Customer.Name" equals "Joe Bloggs"
  And the event "metaData.Scheme.Customer.Age" equals "21"
  And the event "metaData.Scheme.Level" equals "Blue"
