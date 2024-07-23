require 'os'

When('I run {string}') do |scenario_name|
  execute_command 'run-scenario', scenario_name
end

When('I configure the maze endpoint') do
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

# Then(/^the exception is a PathError for request (\d+)$/) do |request_index|
#   body = find_request(request_index)[:body]
#   error_class = body["events"][0]["exceptions"][0]["errorClass"]
#   if ['1.11', '1.12', '1.13', '1.14', '1.15'].include? ENV['GO_VERSION']
#     assert_equal(error_class, '*os.PathError')
#   else
#     assert_equal(error_class, '*fs.PathError')
#   end
# end
# # 
# Then(/^the event unhandled sessions count equals (\d+) for request (\d+)$/) do |count, request_index|
#   step "the payload field \"events.0.session.events.unhandled\" equals #{count} for request #{request_index}"
# end
# 
# Then(/^the event handled sessions count equals (\d+) for request (\d+)$/) do |count, request_index|
#   step "the payload field \"events.0.session.events.handled\" equals #{count} for request #{request_index}"
# end
# 
# Then(/^the number of sessions started equals (\d+) for request (\d+)$/) do |count, request_index|
#   step "the payload field \"sessionCounts.0.sessionsStarted\" equals #{count} for request #{request_index}"
# end
#  
# Then("the in-project frames of the stacktrace are:") do |table|
#   body = find_request(0)[:body]
#   stacktrace = body["events"][0]["exceptions"][0]["stacktrace"]
#   found = 0 # counts matching frames and ensures ordering is correct
#   expected = table.hashes.length
#   stacktrace.each do |frame|
#     if found < expected and frame["inProject"] and
#         frame["file"] == table.hashes[found]["file"] and
#         frame["method"] == table.hashes[found]["method"]
#       found = found + 1
#     end
#   end
#   assert_equal(found, expected, "expected #{expected} matching frames but found #{found}. stacktrace:\n#{stacktrace}")
# end
# 
# 
# Then("the exception {string} is one of:") do |key, table|
#   body = find_request(0)[:body]
#   exception = body["events"][0]["exceptions"][0]
#   options = table.raw.flatten
#   assert(options.include?(exception[key]), "expected '#{key}' to be one of #{options}")
# end
# 