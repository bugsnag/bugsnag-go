NEGRONI_PORT = 9040

When('I go to the negroni route {string}') do |route|
  steps %(
    And I wait for 1 second
    And I open the URL "http://localhost:#{NEGRONI_PORT}#{route}"
    And I wait for 1 second
  )
end

When('I send a request to {string} on the negroni app that might fail') do |path|
  run_command(@script_env,
              "curl http://localhost:#{NEGRONI_PORT}#{path}",
              must_pass: false)
end
