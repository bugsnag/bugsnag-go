Feature: Running docker services

    Background:
        Given I stop the compose stack

    Scenario: A service can be built and run
        Given I have built the service "test_1"
        When I start the service "test_1"
        Then the service "test_1" should be running

    Scenario: A service can be stopped
        Given I have built the service "test_1"
        When I start the service "test_1"
        Then the service "test_1" should be running
        When I stop the service "test_1"
        Then the service "test_1" should not be running

    Scenario: A service with dependencies can be built and run
        Given I have built the service "test_2"
        When I start the service "test_2"
        Then the service "test_2" should be running
        And the service "dep" should be running

    Scenario: A service can be run with a different command
        Given I have built the service "test_1"
        When I run the service "test_1" with the command "bundle exec ruby server.rb"
        Then the service "test_1" should be running
        And I kill the service "test_1"

    Scenario: A service can be started from a different stack
        Given I have built the service "test_1" from the stack "features/fixtures/other-compose.yml"
        When I start the service "test_1" from the stack "features/fixtures/other-compose.yml"
        Then the service "test_1" should be running