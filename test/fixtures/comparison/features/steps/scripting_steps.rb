When(/^I send an? "(.+)"-type request$/) do |request_type|
  steps %Q{
    When I set environment variable "request_type" to "#{request_type}"
    And I run the script "features/scripts/send_request.sh"
    And I wait for 1 seconds
  }
end
