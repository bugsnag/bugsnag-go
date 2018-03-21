Feature: Running a docker stack

    Background:
        Given I start the compose stack

    Scenario: The docker stack can be brought up
        Then the service "test_1" should be running
        And the service "test_2" should be running
        And the service "dep" should be running

    Scenario: The docker stack can be taken down
        When I stop the compose stack
        Then the service "test_1" should not be running
        And the service "test_2" should not be running
        And the service "dep" should not be running

    Scenario: A different docker stack can be brought up
        When I stop the compose stack
        And I start the compose stack "features/fixtures/other-compose.yml"
        Then the service "test_1" should be running
        And the service "test_2" should not be running
        And the service "dep" should not be running

    Scenario: A different docker stack can be taken down
        When I stop the compose stack
        And I start the compose stack "features/fixtures/other-compose.yml"
        And I stop the compose stack "features/fixtures/other-compose.yml"
        Then the service "test_1" should not be running

    Scenario: The default docker stack can be changed
        When I stop the compose stack
        And I am using the docker-compose stack "features/fixtures/other-compose.yml"
        And I start the compose stack
        Then the service "test_1" should be running
