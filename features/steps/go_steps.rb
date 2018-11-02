require 'net/http'

When('I configure the bugsnag endpoints') do
  steps %Q{
    When I set environment variable "NOTIFY_ENDPOINT" to "http://#{current_ip}:#{MOCK_API_PORT}"
    When I set environment variable "SESSIONS_ENDPOINT" to "http://#{current_ip}:#{MOCK_API_PORT}"
  }
end

When('I configure the bugsnag notify endpoint only') do
  steps %Q{
    When I set environment variable "NOTIFY_ENDPOINT" to "http://#{current_ip}:#{MOCK_API_PORT}"
  }
end

When('I configure the bugsnag sessions endpoint only') do
  steps %Q{
    When I set environment variable "SESSIONS_ENDPOINT" to "http://#{current_ip}:#{MOCK_API_PORT}"
  }
end

When('I configure with the {string} configuration and send an error') do |testcase|
  run_command(@script_env, "cd features/fixtures/configure_and_send; go run main.go -case=\"#{testcase}\" -send=error")
end

When('I configure with the {string} configuration and send an error with crash') do |testcase|
  run_command(@script_env, "cd features/fixtures/configure_and_send; go run main.go -case=\"#{testcase}\" -send=error",  must_pass: false)
end

When('I configure with the {string} configuration and send a session') do |testcase|
  run_command(@script_env, "cd features/fixtures/configure_and_send; go run main.go -case=\"#{testcase}\" -send=session")
end

When("I run the http-net test server with the {string} configuration") do |testcase|
  run_command(@script_env, "cd features/fixtures/http_net; go run main.go -case=\"#{testcase}\"")
end

When("I run the http-net test server with the {string} configuration and crashes") do |testcase|
  run_command(@script_env, "cd features/fixtures/http_net; go run main.go -case=\"#{testcase}\"",  must_pass: false)
end

When(/^I wait for the app to open port "(.*)"$/) do |port|
  wait_for_port(port)
end