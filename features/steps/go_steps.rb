require 'net/http'

When(/^I wait for the app to open port "(.*)"$/) do |port|
  wait_for_port(port)
end

Then(/^the session contained the api key "(.*)"$/) do |api_key|
  step "the \"bugsnag-api-key\" header equals \"#{api_key}\""
end

Then(/^the session in request (\d+) contained the api key "(.*)"$/) do |request_index, api_key|
  step "the \"bugsnag-api-key\" header equals \"#{api_key}\" for request #{request_index}"
end

Then(/^the request (\d+) contained the api key "(.*)"$/) do |request_index, api_key|
  steps %Q{
    Then the "bugsnag-api-key" header equals "#{api_key}" for request #{request_index}
    And the payload field "apiKey" equals "#{api_key}" for request #{request_index}
  }
end

Then(/^the events unhandled sessions count equals (\d+) for request (\d+)$/) do |count, request_index|
  step "the payload field \"events.0.session.events.unhandled\" equals #{count} for request #{request_index}"
end

Then(/^the events handled sessions count equals (\d+) for request (\d+)$/) do |count, request_index|
  step "the payload field \"events.0.session.events.handled\" equals #{count} for request #{request_index}"
end

Then(/^the number of sessions started equals (\d+) in request (\d+)$/) do |count, request_index|
  step "the payload field \"sessionCounts.0.sessionsStarted\" equals #{count} for request #{request_index}"
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

When(/^I curl the URL "([^"]+)"$/) do |url|
  `curl -v #{url}`
end

