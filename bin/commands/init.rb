require 'fileutils'

FileUtils.mkdir_p([
  "features/fixtures",
  "features/scripts",
  "features/support",
  "features/steps"
])
if File.exist? "Gemfile"
  contents = open("Gemfile", &:read)
  unless contents.include? 'bugsnag-maze-runner'
    File.open("Gemfile", "a+") do |file|
      file.puts "gem 'bugsnag-maze-runner', :git => 'git@github.com:bugsnag/maze-runner'"
    end
  end
else
  open("Gemfile", "w") do |file|
    file.puts <<-CONTENTS
source 'https://rubygems.org'

gem 'bugsnag-maze-runner', :git => 'git@github.com:bugsnag/maze-runner'
CONTENTS
  end
end

unless File.exist? "features/support/env.rb"
  open("features/support/env.rb", "w") do |file|
    file.puts <<-CONTENTS
# Any 'run once' setup should go here as this file is evaluated
# when the environment loads.
# Any helper functions added here will be available in step
# definitions

# Scenario hooks
Before do
# Runs before every Scenario
end

After do
# Runs after every Scenario
end

at_exit do
# Runs when the test run is completed
end
CONTENTS
  end
end
# Add sample steps file
unless File.exist? "features/scripts/send_request.sh"
  open("features/scripts/send_request.sh", "w") do |file|
    file.puts <<-CONTENTS
#!/usr/bin/env ruby

require 'net/http'

# Sends a request to the mock server running on the port
# specified by the MOCK_API_PORT environment variable
http = Net::HTTP.new('localhost', ENV['MOCK_API_PORT'])
request = Net::HTTP::Post.new('/')
request['Content-Type'] = 'application/json'
request.body = '{"dessert":"' + ENV['menu_item'] + '"}'
http.request(request)
CONTENTS
  end
  FileUtils.chmod("+x", "features/scripts/send_request.sh")
end
unless File.exist? "features/validation.feature"
  open("features/validation.feature", "w") do |file|
    file.puts <<-CONTENTS
Feature: Ordering anything with lemon

Scenario: Lemon cake is in the online menu
  When I select "lemon cake" on the website
  Then I should receive a request
  And the payload body matches the JSON fixture in "features/fixtures/lemon.json"

Scenario: Lemon meringue is in the online menu
  When I select "lemon meringue" on the website
  Then I should receive a request
  And the payload body matches the JSON fixture in "features/fixtures/lemon.json"
CONTENTS
  end
end
unless File.exist? "features/fixtures/lemon.json"
  open("features/fixtures/lemon.json", "w") do |file|
    file.puts <<-CONTENTS
{
"dessert": "^lemon"
}
CONTENTS
  end
end
unless File.exist? "features/steps/website_steps.rb"
  open("features/steps/website_steps.rb", "w") do |file|
    file.puts <<-CONTENTS
When("I select {string} on the website") do |menu_item|
steps %Q{
  When I set environment variable "menu_item" to "\#{menu_item}"
  And I run the script "features/scripts/send_request.sh"
  And I wait for 1 second
}
end
CONTENTS
  end
end
puts "Initialized sample maze runner configuration to 'features/'"
unless ARGV.include? "--skip-install"
  puts "Installing dependencies..."
  `gem install bundler`
  `bundle install`
end
puts "Done! Run `bundle exec bugsnag-maze-runner` to test your configuration."
