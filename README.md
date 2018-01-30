# Maze Runner 🏃

A test runner for validating requests

## How it works

The test harness launches a mock API which awaits requests from sample
applications. Using the runner, each scenario is executed and the requests are
validates to have to correct fields and values. Uses Gherkin and Cucumber under
the hood to draft semantic tests.

## Setting up a new project

1. Add a Gemfile to the root of your project:

   ```ruby
   source "https://rubygems.org"

   gem "bugsnag-maze-runner", git: "https://github.com/bugsnag/bugsnag-maze-runner"
   ```

2. Add a `features` directory to the root of your project. This is where
   scenarios and helper functions will live.

3. Initializing an application and triggering requests should be done through
   scripts. Add script files to `features/fixtures/` to start each variant of
   your supported applications or endpoint as needed. The script files should be
   marked as executable. The `MOCK_API_PORT` environment variable is provided to
   every script to aid configuration. Additional environment variables can be
   configured within scenarios.

4. Add any setup which should be run once before all of the scenarios to
   `features/support/env.rb`. Any setup which should be run before or after each
   scenario can go into special `Before` and `After` functions respectively.

   ```ruby
   # A helper function included with the harness to run commands and
   # only print output when needed
   run_required_commands([
     ["bundle", "install"]
   ])

   # Maybe shell out to something directly, if necessary
   `echo Peanut Butter Jelly Time`

   # Run before every scenarios
   Before do
     clean_build_artifacts
   end
   ```

5. Write your features. The harness provides a number of reusable step
   definitions for interacting with scripts, setting environment variables, and
   inspecting output. Each new feature should go into a `.feature` file in the
   `features` directory.

   ```gherkin
    Feature: Sinatra support

    Scenario: Sinatra unhandled exception
        When I set environment variable "BUGSNAG_API_KEY" to "a35a2a72bd230ac0aa0f52715bbdc6aa"
        And I start a Sinatra app
        And I navigate to the route "/syntax-error"
        Then I should receive a request
        And the request is a valid for the error reporting API
        And the "Bugsnag-API-Key" header equals "a35a2a72bd230ac0aa0f52715bbdc6aa"
        And the event "unhandled" is true
   ```

   This example includes a mix of the included steps as well as custom ones
   specific to the library being tested. `When I set an environment variable` is
   provided by default while `When I start a Sinatra app` is defined in a custom
   steps file in `features/steps/`, wrapping other included steps:

   ```ruby
   When("I start a Sinatra app") do
     set_script_env "DEMO_APP_PORT", "#{DEMO_APP_PORT}"
     steps %Q{
       When I run the script "features/fixtures/run_sinatra_app.sh"
       And I wait for 8 seconds
     }
   end

   When("I navigate to the route {string}") do |route|
     steps %Q{
       When I open the URL "http://localhost:#{DEMO_APP_PORT}#{route}"
       And I wait for 1 second
     }
   end
   ```

   In addition, any helper functions or instance variables defined in
   `features/support/env.rb` are available to step files. See the included
   `_step.rb` files for examples.

6. Run your tests with `bugsnag-maze-runner`!

## Contributing

If steps would be useful for different projects running the maze, add the to
`lib/features/steps/`. If there are useful helper functions, add them to
`lib/features/support/env.rb`.

### Running the tests

bugsnag-maze-runner uses test-unit and minunit to bootstrap itself and run the
sample app suites in the test fixtures. Run `rake test` to run the suite.
