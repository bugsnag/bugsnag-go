Given("I am using the docker-compose stack {string}") do |filename|
  set_compose_file filename
end

Given("I have built the service {string}") do |service|
  build_service service
end

Given("I have built the service {string} from the stack {string}") do |service, filename|
  build_service(service, filename)
end

When("I start the service {string}") do |service|
  start_service service
end

When("I start the service {string} from the stack {string}") do |service, filename|
  start_service(service, filename)
end

When("I stop the service {string}") do |service|
  # Warning! By docker-compose design this command will not stop containers started by the Run command
  stop_service service
end

When("I stop the service {string} from the stack {string}") do |service, filename|
  stop_service(service, filename)
end

When("I kill the service {string}") do |service|
  kill_service service
end

When("I kill the service {string} from the stack {string}") do |service, filename|
  kill_service(service, filename)
end

When("I start the compose stack") do
  start_stack
end

When("I start the compose stack {string}") do |filename|
  start_stack filename
end

When("I stop the compose stack") do
  stop_stack
end

When("I stop the compose stack {string}") do |filename|
  stop_stack filename
end

When("I run the command {string} on the service {string}") do |command, service|
  run_command_on_service(command, service)
end

When("I run the command {string} on the service {string} from the stack {string}") do |command, service, filename|
  run_command_on_service(command, service, filename)
end

When("I run the service {string} with the command {string}") do |service, command|
  run_service_with_command(service, command)
end

When("I run the service {string} with the command {string} from the stack {string}") do |service, command, filename|
  run_service_with_command(service, command, filename)
end

Then("the service {string} should be running") do |service|
  test_service_running(service)
end

Then("the service {string} from the stack {string} should be running") do |service, filename|
  test_service_running(service, filename)
end

Then("the service {string} should not be running") do |service|
  test_service_running(service, false)
end

Then("the service {string} from the stack {string} should not be running") do |service, filename|
  test_service_running(service, false, filename)
end