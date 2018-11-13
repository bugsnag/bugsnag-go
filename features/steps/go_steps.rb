require 'net/http'

When(/^I wait for the app to open port "(.*)"$/) do |port|
  wait_for_port(port)
end

Then(/^the request(?: (\d+))? is a valid error report with api key "(.*)"$/) do |request_index, api_key|
  request_index ||= 0
  steps %Q{
    And the request #{request_index} is valid for the error reporting API
    And the "bugsnag-api-key" header equals "#{api_key}" for request #{request_index}
    And the payload field "apiKey" equals "#{api_key}" for request #{request_index}
  }
end

Then(/^the request(?: (\d+))? is a valid session report with api key "(.*)"$/) do |request_index, api_key|
  request_index ||= 0
  steps %Q{
    And the request #{request_index} is valid for the session tracking API
    And the "bugsnag-api-key" header equals "#{api_key}" for request #{request_index}
  }
end

Then(/^the event unhandled sessions count equals (\d+) for request (\d+)$/) do |count, request_index|
  step "the payload field \"events.0.session.events.unhandled\" equals #{count} for request #{request_index}"
end

Then(/^the event handled sessions count equals (\d+) for request (\d+)$/) do |count, request_index|
  step "the payload field \"events.0.session.events.handled\" equals #{count} for request #{request_index}"
end

Then(/^the number of sessions started equals (\d+) for request (\d+)$/) do |count, request_index|
  step "the payload field \"sessionCounts.0.sessionsStarted\" equals #{count} for request #{request_index}"
end

Then(/^I wait to receive a request$/) do
  step "I wait to receive 1 request"
end

Then(/^I wait to receive (\d+) requests?$/) do |request_count|
  max_attempts = 50
  attempts = 0
  received = false
  until (attempts >= max_attempts) || received
    attempts += 1
    received = (stored_requests.size == request_count)
    sleep 0.2
  end
  raise "Requests not received in 10s (received #{stored_requests.size})" unless received
  # Wait an extra second to ensure there are no further requests
  sleep 1
  assert_equal(request_count, stored_requests.size, "#{stored_requests.size} requests received")
end

When("I run the go service {string} with the test case {string}") do |service, testcase|
  run_service_with_command(service, "go run main.go -test=\"#{testcase}\"")
end
