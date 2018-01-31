require 'test/unit'
require 'minitest'
require 'open-uri'
require 'json'

include Test::Unit::Assertions

Then(/^I should receive (\d+) requests?$/) do |request_count|
  assert_equal(request_count, stored_requests.size, "#{stored_requests.size} requests received")
end
Then(/^I should receive a request$/) do
  step "I should receive 1 request"
end
Then(/^the "(.+)" header is not null$/) do |header_name|
  assert_not_nil(stored_requests.first[:request][header_name],
                "The '#{header_name}' header should not be null")
end
Then(/^the "(.+)" header equals "(.+)"$/) do |header_name, header_value|
  assert_equal(header_value, stored_requests.first[:request][header_name])
end
Then("the {string} header is a timestamp") do |header_name|
  header = stored_requests.first[:request][header_name]
  assert_match(/^\d{4}\-\d{2}\-\d{2}T\d{2}:\d{2}:[\d\.]+Z?$/, header)
end
Then("the payload field {string} is true") do |field_path|
  assert_equal(true, read_key_path(stored_requests.first[:body], field_path))
end
Then("the payload field {string} is false") do |field_path|
  assert_equal(false, read_key_path(stored_requests.first[:body], field_path))
end
Then(/^the payload field "(.+)" is not null$/) do |field_path|
  assert_not_nil(read_key_path(stored_requests.first[:body], field_path),
                "The field '#{field_path}' should not be null")
end
Then(/^the payload field "(.+)" equals (\d+)$/) do |field_path, int_value|
  assert_equal(int_value, read_key_path(stored_requests.first[:body], field_path))
end
Then(/^the payload field "(.+)" equals "(.+)"$/) do |field_path, string_value|
  assert_equal(string_value, read_key_path(stored_requests.first[:body], field_path))
end
Then(/^the payload field "(.+)" starts with "(.+)"$/) do |field_path, string_value|
  value = read_key_path(stored_requests.first[:body], field_path)
  assert_kind_of String, value
  assert(value.start_with?(string_value), "Field '#{field_path}' value ('#{value}') does not start with '#{string_value}'")
end
Then(/^the payload field "(.+)" ends with "(.+)"$/) do |field_path, string_value|
  value = read_key_path(stored_requests.first[:body], field_path)
  assert_kind_of String, value
  assert(value.end_with? string_value, "Field '#{field_path}' does not end with '#{string_value}'")
end
Then(/^the payload field "(.+)" is an array with (\d+) elements?$/) do |field, count|
  value = read_key_path(stored_requests.first[:body], field)
  assert_kind_of Array, value
  assert_equal(count, value.length)
end
Then(/^the payload field "(.+)" is a non-empty array$/) do |field|
  value = read_key_path(stored_requests.first[:body], field)
  assert_kind_of Array, value
  assert(value.length > 0, "the field '#{field}' must be a non-empty array")
end
Then(/^each element in payload field "(.+)" has "(.+)"$/) do |key_path, element_key_path|
  value = read_key_path(stored_requests.first[:body], key_path)
  assert_kind_of Array, value
  value.each do |element|
    assert_not_nil(read_key_path(element, element_key_path),
           "Each element in '#{key_path}' must have '#{element_key_path}'")
  end
end
