GIN_PORT = 9050

When('I go to the gin route {string}') do |route|
  steps %(
    And I wait for 1 second
    And I open the URL "http://localhost:#{GIN_PORT}#{route}"
    And I wait for 1 second
  )
end

When('I send a request to {string} on the gin app that might fail') do |path|
  run_command(@script_env,
              "curl http://localhost:#{GIN_PORT}#{path}",
              must_pass: false)
end
