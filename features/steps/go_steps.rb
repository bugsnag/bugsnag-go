When('I configure the bugsnag endpoints') do
  steps %Q{
    When I set environment variable "NOTIFY_ENDPOINT" to "http://localhost:#{MOCK_API_PORT}"
    When I set environment variable "SESSIONS_ENDPOINT" to "http://localhost:#{MOCK_API_PORT}"
  }
end

When('I configure the bugsnag notify endpoint only') do
  steps %Q{
    When I set environment variable "NOTIFY_ENDPOINT" to "http://localhost:#{MOCK_API_PORT}"
  }
end

When('I configure the bugsnag sessions endpoint only') do
  steps %Q{
    When I set environment variable "SESSIONS_ENDPOINT" to "http://localhost:#{MOCK_API_PORT}"
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

Then("the request used payload v4 headers") do
  steps %Q{
    Then the "bugsnag-api-key" header is not null
    And the "bugsnag-payload-version" header equals "4"
  }
  # TODO: format using zulu style RFC3339
  # And the "bugsnag-sent-at" header is a timestamp
end

Then("the request contained the api key {string}") do |api_key|
  steps %Q{
    Then the "bugsnag-api-key" header equals "#{api_key}"
    And the payload field "apiKey" equals "#{api_key}"
  }
end

When('I set the legacy endpoint only') do
  steps %(
    When I set environment variable "ENDPOINT" to "http://localhost:#{MOCK_API_PORT}"
  )
end
