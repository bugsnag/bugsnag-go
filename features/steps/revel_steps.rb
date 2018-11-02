When('I set the {string} config variable {string} to {string}') do |fixture, prop_name, prop_value|
  replace_revel_conf(fixture: fixture,
                     property_name: prop_name,
                     property_value: prop_value)
end

When('I work with a new {string} app') do |fixture|
  skip_this_scenario if go_version_is_unsupported
  conf_path = "features/fixtures/#{fixture}/conf/app.conf"
  run_command("cp #{conf_path}-default #{conf_path}")
end

When('I configure the bugsnag sessions endpoint in the config file for {string}') do |fixture|
  replace_revel_conf(fixture: fixture,
                     property_name: 'bugsnag.endpoints.sessions',
                     property_value: "http:\\/\\/localhost:#{MOCK_API_PORT}")
  replace_revel_conf(fixture: fixture,
                     property_name: 'bugsnag.endpoints.notify',
                     property_value: 'http:\\/\\/localhost:80\\/somewhere-else')
end

When('I configure the bugsnag endpoint in the config file for {string}') do |fixture|
  replace_revel_conf(fixture: fixture,
                     property_name: 'bugsnag.endpoints.notify',
                     property_value: "http:\\/\\/localhost:#{MOCK_API_PORT}")
end

When('I configure the legacy bugsnag endpoint in the config file for {string}') do |fixture|
  replace_revel_conf(fixture: fixture,
                     property_name: 'bugsnag.endpoint',
                     property_value: "http:\\/\\/localhost:#{MOCK_API_PORT}")
end

When('I configure the legacy bugsnag endpoint in as an environment variable') do
  steps %(
    Given I set environment variable "ENDPOINT" to "http://localhost:#{MOCK_API_PORT}"
  )
end

When('I go to the route {string}') do |route|
  steps %(
    When I open the URL "http://localhost:#{REVEL_PORT}#{route}"
    And I wait for 1 second
  )
end
