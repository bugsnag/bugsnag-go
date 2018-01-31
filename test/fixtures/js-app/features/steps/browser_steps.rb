When("I start a local server") do
  steps %Q{
    When I run the script "features/fixtures/local_server.sh"
    And I wait for 6 seconds
  }
end
When("I navigate to the {string} page in {string}") do |name, browser|
  step("I open the URL \"#{create_url_for_fixture(name)}\" in \"#{browser}\"")
  step("I wait for 2 seconds")
end
When("I navigate to the {string} page") do |name|
  step("I open the URL \"#{create_url_for_fixture(name)}\" in a browser")
  step("I wait for 2 seconds")
end
