require 'test/unit'
require 'minitest'
require 'open-uri'
require 'json'

include Test::Unit::Assertions

def find_request(request_index)
  request_index ||= 0
  return stored_requests[request_index]
end

Then(/^I should receive (\d+) requests?$/) do |request_count|
  assert_equal(request_count, stored_requests.size, "#{stored_requests.size} requests received")
end
Then(/^I should receive a request$/) do
  step "I should receive 1 request"
end
Then(/^I should receive no requests$/) do
  step "I should receive 0 request"
end
Then(/^the "(.+)" header is not null(?: for request (\d+))?$/) do |header_name, request_index|
  assert_not_nil(find_request(request_index)[:request][header_name],
                "The '#{header_name}' header should not be null")
end
Then(/^the "(.+)" header equals "(.+)"(?: for request (\d+))?$/) do |header_name, header_value, request_index|
  assert_equal(header_value, find_request(request_index)[:request][header_name])
end

Then(/^the "(.+)" header is a timestamp(?: for request (\d+))?$/) do |header_name, request_index|
  header = find_request(request_index)[:request][header_name]
  assert_match(/^\d{4}\-\d{2}\-\d{2}T\d{2}:\d{2}:[\d\.]+Z?$/, header)
end

Then(/^the request (\d+) is valid for the Android Mapping API$/) do |request_index|
  parts = find_request(request_index)[:body]
  assert_equal(6, parts.size)
  assert_not_nil(parts["proguard"])
  assert_not_nil(parts["apiKey"])
  assert_not_nil(parts["appId"])
  assert_not_nil(parts["versionCode"])
  assert_not_nil(parts["buildUUID"])
  assert_not_nil(parts["versionName"])
end

Then(/^the request (\d+) has (\d+) parts$/) do |request_index, part_count|
  parts = find_request(request_index)[:body]
  assert_equal(part_count, parts.size)
end

Then(/^the part "(.+)" for request (\d+) is not null$/) do |part_key, request_index|
  parts = find_request(request_index)[:body]
  assert_not_nil(parts[part_key], "The field '#{part_key}' should not be null")
end

Then(/^the part "(.+)" for request (\d+) equals "(.+)"$/) do |part_key, request_index, expected_value|
  parts = find_request(request_index)[:body]
  assert_not_nil(parts[part_key], expected_value)
end

Then(/^the payload body does not match the JSON fixture in "(.+)"(?: for request (\d+))?$/) do |fixture_path, request_index|
  payload_value = find_request(request_index)[:body]
  expected_value = JSON.parse(open(fixture_path, &:read))
  result = value_compare(expected_value, payload_value)
  assert_false(result.equal?, "Payload:\n#{payload_value}\nExpected:#{expected_value}")
end
Then(/^the payload body matches the JSON fixture in "(.+)"(?: for request (\d+))?$/) do |fixture_path, request_index|
  payload_value = find_request(request_index)[:body]
  expected_value = JSON.parse(open(fixture_path, &:read))
  result = value_compare(expected_value, payload_value)
  assert_true(result.equal?, "The payload field '#{result.keypath}' does not match the fixture:\n #{result.reasons.join('\n')}")
end
Then(/^the payload field "(.+)" matches the JSON fixture in "(.+)"(?: for request (\d+))?$/) do |field_path, fixture_path, request_index|
  payload_value = read_key_path(find_request(request_index)[:body], field_path)
  expected_value = JSON.parse(open(fixture_path, &:read))
  result = value_compare(expected_value, payload_value)
  assert_true(result.equal?, "The payload field '#{result.keypath}' does not match the fixture:\n #{result.reasons.join('\n')}")
end
Then(/^the payload field "(.+)" is true(?: for request (\d+))?$/) do |field_path, request_index|
  assert_equal(true, read_key_path(find_request(request_index)[:body], field_path))
end
Then(/^the payload field "(.+)" is false(?: for request (\d+))?$/) do |field_path, request_index|
  assert_equal(false, read_key_path(find_request(request_index)[:body], field_path))
end

Then(/^the payload field "(.+)" is null(?: for request (\d+))?$/) do |field_path, request_index|
  value = read_key_path(find_request(request_index)[:body], field_path)
  assert_nil(value, "The field '#{field_path}' should be null but is #{value}")
end
Then(/^the payload field "(.+)" is not null(?: for request (\d+))?$/) do |field_path, request_index|
  assert_not_nil(read_key_path(find_request(request_index)[:body], field_path),
                "The field '#{field_path}' should not be null")
end
Then(/^the payload field "(.+)" equals (\d+)(?: for request (\d+))?$/) do |field_path, int_value, request_index|
  assert_equal(int_value, read_key_path(find_request(request_index)[:body], field_path))
end
Then(/^the payload field "(.+)" equals "(.+)"(?: for request (\d+))?$/) do |field_path, string_value, request_index|
  assert_equal(string_value, read_key_path(find_request(request_index)[:body], field_path))
end
Then(/^the payload field "(.+)" starts with "(.+)"(?: for request (\d+))?$/) do |field_path, string_value, request_index|
  value = read_key_path(find_request(request_index)[:body], field_path)
  assert_kind_of String, value
  assert(value.start_with?(string_value), "Field '#{field_path}' value ('#{value}') does not start with '#{string_value}'")
end
Then(/^the payload field "(.+)" ends with "(.+)"(?: for request (\d+))?$/) do |field_path, string_value, request_index|
  value = read_key_path(find_request(request_index)[:body], field_path)
  assert_kind_of String, value
  assert(value.end_with? string_value, "Field '#{field_path}' does not end with '#{string_value}'")
end
Then(/^the payload field "(.+)" is an array with (\d+) elements?(?: for request (\d+))?$/) do |field, count, request_index|
  value = read_key_path(find_request(request_index)[:body], field)
  assert_kind_of Array, value
  assert_equal(count, value.length)
end
Then(/^the payload field "(.+)" is a non-empty array(?: for request (\d+))?$/) do |field, request_index|
  value = read_key_path(find_request(request_index)[:body], field)
  assert_kind_of Array, value
  assert(value.length > 0, "the field '#{field}' must be a non-empty array")
end
Then(/^each element in payload field "(.+)" has "(.+)"(?: for request (\d+))?$/) do |key_path, element_key_path, request_index|
  value = read_key_path(find_request(request_index)[:body], key_path)
  assert_kind_of Array, value
  value.each do |element|
    assert_not_nil(read_key_path(element, element_key_path),
           "Each element in '#{key_path}' must have '#{element_key_path}'")
  end
end
