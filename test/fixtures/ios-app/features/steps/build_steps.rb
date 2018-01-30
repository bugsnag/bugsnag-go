When("I build the app") do
  steps %Q{
    When I run the script "features/fixtures/build-app.sh"
    And I wait for 8 seconds
  }
end
When("I launch the app") do
  steps %Q{
    When I run the script "features/fixtures/launch-app.sh"
    And I wait for 5 seconds
  }
end
When("I configure the app to trigger {string}") do |event_type|
end
