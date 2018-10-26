When('I configure the bugsnag endpoint') do
  steps %Q{
    When I set environment variable "NOTIFY_ENDPOINT" to "http://localhost:#{MOCK_API_PORT}"
  }
end

When('I configure with the {string} configuration and send an error') do |testcase|
  run_command(@script_env, "cd features/fixtures/configure_and_send; go run main.go -case=\"#{testcase}\"")
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
