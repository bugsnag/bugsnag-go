When(/^I set environment variable "(.+)" to "(.+)"$/) do |key, value|
  set_script_env key, value
end
When(/^I run the script "(.+)"$/) do |script_path|
  run_script script_path
end
When(/^I open the URL "([^"]+)"$/) do |url|
  begin
    open(url, &:read)
  rescue OpenURI::HTTPError
  end
end
When(/^I open the URL "([^"]+)" in "([^"]+)"$/) do |url, browser|
  case RbConfig::CONFIG['host_os']
  when /mswin|mingw|cygwin/
    pending
  when /darwin/
    # `open` doesn't respect query params
    system "osascript -e 'tell application \"#{browser}\" to open location \"#{url}\"'"
  when /linux|bsd/
    pending
  end
end
When(/^I open the URL "(.+)" in a browser$/) do |url|
  case RbConfig::CONFIG['host_os']
  when /mswin|mingw|cygwin/
    system "start #{url}"
  when /darwin/
    step("I open the URL \"#{url}\" in \"Safari\"")
  when /linux|bsd/
    system "xdg-open '#{url}'"
  end
end
When(/^I wait for (\d+) seconds?$/) do |seconds|
  sleep(seconds)
end
