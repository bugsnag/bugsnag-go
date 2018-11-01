NEGRONI_PORT = 9030

When('I go to the negroni route {string}') do |route|
  steps %(
    And I wait for 1 second
    And I open the URL "http://localhost:#{NEGRONI_PORT}#{route}"
    And I wait for 1 second
  )
end

When('I set the legacy endpoint only') do
  steps %(
    When I set environment variable "ENDPOINT" to "http://localhost:#{MOCK_API_PORT}"
  )
end

When('I am working with a new negroni app') do
  run_command('killall negroni || true')
end

When('I send a request to {string} on the negroni app that might fail') do |path|
  run_command(@script_env,
              "curl http://localhost:#{NEGRONI_PORT}#{path}",
              must_pass: false)
end
