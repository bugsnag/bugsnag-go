When("I start a local server") do
  steps %Q{
    When I run the script "features/fixtures/local_server.sh"
    And I wait for 6 seconds
  }
end
