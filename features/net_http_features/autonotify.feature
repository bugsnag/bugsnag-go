Feature: Using auto notify

Background:
  Given I set environment variable "API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And I configure the bugsnag endpoint
  And I have built the service "nethttp"
  And I stop the service "nethttp"
  And I set environment variable "SERVER_PORT" to "4512"
  

Scenario: An error report is sent when an unhandled crash occurs
  Given I set environment variable "AUTO_CAPTURE_SESSIONS" to "false"
  When I start the service "nethttp"
  And I wait for the app to open port "4512"
  And I wait for 1 seconds
  And I open the URL "http://localhost:4512/unhandled"
  Then I wait to receive 1 request
  And the request is valid for the error reporting API
  And the request contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
  And the event "unhandled" is true 
  And the exception "errorClass" equals "*runtime.TypeAssertionError"
  And the exception "message" equals "interface conversion: interface {} is struct {}, not string"

# Scenario: An error report is sent when a go routine crashes which is protected by auto notify
#   When I start the service "nethttp"
#   And I wait for the app to open port "4512"
#   And I wait for 1 seconds
#   And I open the URL "http://localhost:4512/autonotify"
#   And I wait for 1 seconds
#   Then I should receive 3 requests
#   And the request 0 is valid for the error reporting API
#   And the request 0 contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
#   And the request 1 is valid for the session tracking API
#   And the session in request 1 contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
#   And the request 2 is valid for the error reporting API
#   And the request 2 contained the api key "a35a2a72bd230ac0aa0f52715bbdc6aa"
#   And the event "unhandled" is true for request 0
#   And the exception "errorClass" equals "*errors.errorString" for request 0
#   And the exception "message" equals "Go routine killed with auto notify" for request 0
#   And the event "unhandled" is true for request 2
#   And the exception "errorClass" equals "panic" for request 2
#   And the exception "message" equals "Go routine killed with auto notify [recovered]" for request 2
#   And the events unhandled sessions count equals 1 for request 0
#   And the number of sessions started equals 1 in request 1