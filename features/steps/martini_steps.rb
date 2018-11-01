MARTINI_PORT = 9030

When('I go to the martini route {string}') do |route|
  steps %(
    And I wait for 1 second
    And I open the URL "http://localhost:#{MARTINI_PORT}#{route}"
    And I wait for 1 second
  )
end

When('I am working with a new martini app') do
  run_command('killall martini || true')
end

When('I send a request to {string} on the martini app that might fail') do |path|
  run_command(@script_env,
              "curl http://localhost:#{MARTINI_PORT}#{path}",
              must_pass: false)
end
