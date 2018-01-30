When(/^I set environment variable "(.+)" to "(.+)"$/) do |key, value|
  set_script_env key, value
end
When(/^I run the script "(.+)"$/) do |script_path|
  run_script script_path
end
When(/^I open the URL "(.+)"$/) do |url|
  begin
    open(url, &:read)
  rescue OpenURI::HTTPError
  end
end
When(/^I wait for (\d+) seconds?$/) do |seconds|
  sleep(seconds)
end
