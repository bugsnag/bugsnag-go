When('I configure the bugsnag endpoint') do
  steps %(
    When I set environment variable "NOTIFY_ENDPOINT" to "http://localhost:#{MOCK_API_PORT}"
  )
end

When('I run the app {string}') do |app_path|
  run_command(@script_env, "cd features/fixtures/#{app_path}; go run main.go")
end

Then('the request used payload v4 headers') do
  steps %(
    Then the "bugsnag-api-key" header is not null
    And the "bugsnag-payload-version" header equals "4"
  )
  # TODO: format using zulu style RFC3339 by adding the following line above:
  # And the "bugsnag-sent-at" header is a timestamp
end

Then('the request contained the api key {string}') do |api_key|
  steps %(
    Then the "bugsnag-api-key" header equals "#{api_key}"
    And the payload field "apiKey" equals "#{api_key}"
  )
end

When('I navigate to the route {string}') do |route|
  steps %(
    When I open the URL "http://localhost:#{port}#{route}"
    And I wait for 1 second
  )
end
