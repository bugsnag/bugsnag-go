require 'net/http'
require 'os'

When('I run {string}') do |scenario_name|
  execute_command 'run-scenario', scenario_name
end

When('I configure the base endpoint') do
  steps %(
    When I set environment variable "DEFAULT_MAZE_ADDRESS" to "http://#{local_ip}:9339"
  )
end

def execute_command(action, scenario_name = '')
    address = $address ? $address : "#{local_ip}:9339"

  command = {
    action: action,
    scenario_name: scenario_name,
    notify_endpoint: "http://#{address}/notify",
    sessions_endpoint: "http://#{address}/sessions",
    api_key: $api_key,
  }

  $logger.debug("Queuing command: #{command}")
  Maze::Server.commands.add command

  # Ensure fixture has read the command
  count = 900
  sleep 0.1 until Maze::Server.commands.remaining.empty? || (count -= 1) < 1
  raise 'Test fixture did not GET /command' unless Maze::Server.commands.remaining.empty?
end

def local_ip
  if OS.mac?
    'host.docker.internal'
  else
    ip_addr = `ifconfig | grep -Eo 'inet (addr:)?([0-9]*\\\.){3}[0-9]*' | grep -v '127.0.0.1'`
    ip_list = /((?:[0-9]*\.){3}[0-9]*)/.match(ip_addr)
    ip_list.captures.first
  end
end
