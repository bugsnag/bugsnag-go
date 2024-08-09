Feature: Breadcrumbs

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  Given I configure the bugsnag endpoint
  Given I have built the service "app"

Scenario: Disabling breadcrumbs
  Given I set environment variable "ENABLED_BREADCRUMB_TYPES" to "[]"
  When I run the go service "app" with the test case "disable-breadcrumbs"
  When I wait to receive 2 requests after the start up session

  Then the request 0 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  Then the payload field "events.0.breadcrumbs" is null for request 0

  Then the request 1 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  Then the payload field "events.0.breadcrumbs" is null for request 1

Scenario: Automatic breadcrumbs
  When I run the go service "app" with the test case "automatic-breadcrumbs"
  When I wait to receive 2 requests after the start up session

  Then the request 0 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  Then the payload field "events.0.breadcrumbs" is an array with 1 elements for request 0
  Then the payload field "events.0.breadcrumbs.0.name" equals "Bugsnag loaded" for request 0
  Then the payload field "events.0.breadcrumbs.0.type" equals "state" for request 0

  Then the request 1 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"  
  Then the payload field "events.0.breadcrumbs" is an array with 2 elements for request 1
  Then the payload field "events.0.breadcrumbs.0.name" equals "oops" for request 1
  Then the payload field "events.0.breadcrumbs.0.type" equals "error" for request 1
  Then the payload field "events.0.breadcrumbs.1.name" equals "Bugsnag loaded" for request 1
  Then the payload field "events.0.breadcrumbs.1.type" equals "state" for request 1

Scenario: Setting max breadcrumbs
  Given I set environment variable "ENABLED_BREADCRUMB_TYPES" to "[]"
  Given I set environment variable "MAXIMUM_BREADCRUMBS" to "5"
  When I run the go service "app" with the test case "maximum-breadcrumbs"
  When I wait to receive 1 requests after the start up session

  Then the request 0 is a valid error report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  Then the payload field "events.0.breadcrumbs" is an array with 5 elements
  Then the payload field "events.0.breadcrumbs.0.name" equals "Crumb 9"
  Then the payload field "events.0.breadcrumbs.1.name" equals "Crumb 8"
  Then the payload field "events.0.breadcrumbs.2.name" equals "Crumb 7"
  Then the payload field "events.0.breadcrumbs.3.name" equals "Crumb 6"
  Then the payload field "events.0.breadcrumbs.4.name" equals "Crumb 5"
