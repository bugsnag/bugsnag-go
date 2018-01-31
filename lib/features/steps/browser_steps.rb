When("I navigate to the {string} page in {string}") do |name, browser|
  step("I open the URL \"#{create_url_for_fixture(name)}\" in \"#{browser}\"")
  step("I wait for 2 seconds")
end
When("I navigate to the {string} page") do |name|
  step("I open the URL \"#{create_url_for_fixture(name)}\" in a browser")
  step("I wait for 2 seconds")
end
Then("the script content is in event metadata") do
  step("the event \"metaData.script.content\" is not null")
end
