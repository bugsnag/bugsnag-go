require 'net/http'

Given("I set environment variable {string} to the app directory") do |key|
  step("I set environment variable \"#{key}\" to \"/app/src/test/\"")
end

Given("I set environment variable {string} to the notify endpoint") do |key|
  step("I set environment variable \"#{key}\" to \"http://#{current_ip}:#{MOCK_API_PORT}\"")
end

Given("I set environment variable {string} to the sessions endpoint") do |key|
  # they're the same picture dot gif
  # split them out for the future endpoint splitting work
  step("I set environment variable \"#{key}\" to \"http://#{current_ip}:#{MOCK_API_PORT}\"")
end

Then(/^the request(?: (\d+))? is a valid error report with api key "(.*)"$/) do |request_index, api_key|
  request_index ||= 0
  steps %Q{
    And the request #{request_index} is valid for the error reporting API
    And the "bugsnag-api-key" header equals "#{api_key}" for request #{request_index}
    And the payload field "apiKey" equals "#{api_key}" for request #{request_index}
  }
end

Then(/^the exception is a PathError for request (\d+)$/) do |request_index|
  body = find_request(request_index)[:body]
  error_class = body["events"][0]["exceptions"][0]["errorClass"]
  if ['1.11', '1.12', '1.13', '1.14', '1.15'].include? ENV['GO_VERSION']
    assert_equal(error_class, '*os.PathError')
  else
    assert_equal(error_class, '*fs.PathError')
  end
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

When("I run the go service {string} with the test case {string}") do  |service, testcase|
  run_service_with_command(service, "./run.sh go run . -test=\"#{testcase}\"")
end

Then(/^I wait to receive a request after the start up session$/) do
  step "I wait to receive 1 requests after the start up session"
end

Then(/^I wait to receive (\d+) requests after the start up session?$/) do |request_count|
  max_attempts = 50
  attempts = 0
  start_up_message_received = false
  start_up_message_removed = false
  received = false
  until (attempts >= max_attempts) || received
    attempts += 1
    start_up_message_received ||= (stored_requests.size == 1)
    if start_up_message_received && !start_up_message_removed
      step 'the request is a valid session report with api key "a35a2a72bd230ac0aa0f52715bbdc6aa"'
      stored_requests.shift
      start_up_message_removed = true
      next
    end
    received = (stored_requests.size == request_count)
    sleep 0.2
  end
  raise "Requests not received in 10s (received #{stored_requests.size})" unless received
  # Wait an extra second to ensure there are no further requests
  sleep 1
  assert_equal(request_count, stored_requests.size, "#{stored_requests.size} requests received")
end

Then(/^(\d+) requests? (?:was|were) received$/) do |request_count|
  sleep 1
  assert_equal(request_count, stored_requests.size, "#{stored_requests.size} requests received")
end

Then("the in-project frames of the stacktrace are:") do |table|
  body = find_request(0)[:body]
  stacktrace = body["events"][0]["exceptions"][0]["stacktrace"]
  found = 0 # counts matching frames and ensures ordering is correct
  expected = table.hashes.length
  stacktrace.each do |frame|
    if found < expected and frame["inProject"] and
        frame["file"] == table.hashes[found]["file"] and
        frame["method"] == table.hashes[found]["method"]
      found = found + 1
    end
  end
  assert_equal(found, expected, "expected #{expected} matching frames but found #{found}. stacktrace:\n#{stacktrace}")
end

Then("stack frame {int} contains a local function spanning {int} to {int}") do |frame, val, old_val|
  # Old versions of Go put the line number on the end of the function
  if ['1.7', '1.8'].include? ENV['GO_VERSION']
    step "the \"lineNumber\" of stack frame #{frame} equals #{old_val}"
  else
    step "the \"lineNumber\" of stack frame #{frame} equals #{val}"
  end
end

Then("the exception {string} is one of:") do |key, table|
  body = find_request(0)[:body]
  exception = body["events"][0]["exceptions"][0]
  options = table.raw.flatten
  assert(options.include?(exception[key]), "expected '#{key}' to be one of #{options}")
end
