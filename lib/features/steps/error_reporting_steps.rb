Then("the request is a valid for the error reporting API") do
  steps %Q{
    Then the "Bugsnag-API-Key" header is not null
    And the "Content-Type" header equals "application/json"
    And the "Bugsnag-Payload-Version" header equals "4.0"
    And the "Bugsnag-Sent-At" header is a timestamp

    And the payload field "notifier.name" is not null
    And the payload field "notifier.url" is not null
    And the payload field "notifier.version" is not null
    And the payload field "events" is a non-empty array

    And each element in payload field "events" has "severity"
    And each element in payload field "events" has "severityReason.type"
    And each element in payload field "events" has "unhandled"
    And each element in payload field "events" has "exceptions"
  }
end
Then("the event {string} is true") do |field|
  step "the payload field \"events.0.#{field}\" is true"
end
Then("the event {string} is false") do |field|
  step "the payload field \"events.0.#{field}\" is false"
end
Then(/^the event "(.+)" equals "(.+)"$/) do |field, string_value|
  step "the payload field \"events.0.#{field}\" equals \"#{string_value}\""
end
Then(/^the event "(.+)" starts with "(.+)"$/) do |field, string_value|
  step "the payload field \"events.0.#{field}\" starts with \"#{string_value}\""
end
Then(/^the event "(.+)" ends with "(.+)"$/) do |field, string_value|
  step "the payload field \"events.0.#{field}\" ends with \"#{string_value}\""
end

Then(/^the exception "(.+)" starts with "(.+)"$/) do |field, string_value|
  step "the payload field \"events.0.exceptions.0.#{field}\" starts with \"#{string_value}\""
end
Then(/^the exception "(.+)" ends with "(.+)"$/) do |field, string_value|
  step "the payload field \"events.0.exceptions.0.#{field}\" ends with \"#{string_value}\""
end
Then(/^the exception "(.+)" equals "(.+)"$/) do |field, string_value|
  step "the payload field \"events.0.exceptions.0.#{field}\" equals \"#{string_value}\""
end

Then(/^the "(.+)" of stack frame (\d+) equals "(.+)"$/) do |key, num, value|
  field = "events.0.exceptions.0.stacktrace.#{num}.#{key}"
  step "the payload field \"#{field}\" equals \"#{value}\""
end
Then(/^the "(.+)" of stack frame (\d+) ends with "(.+)"$/) do |key, num, value|
  field = "events.0.exceptions.0.stacktrace.#{num}.#{key}"
  step "the payload field \"#{field}\" ends with \"#{value}\""
end
