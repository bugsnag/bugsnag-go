run_required_commands([
  ["bundle", "install"],
])

Dir.chdir("features/fixtures") do
  run_required_commands([
    ["npm", "install"],
  ])
end

def create_url_for_fixture name
  fixture_path = "http://localhost:8991/#{name}.html"
  params = {
    "apiKey": @script_env['BUGSNAG_API_KEY'],
    "notifyURL": "http://localhost:#{@script_env['MOCK_API_PORT']}"
  }
  fixture_path + "?#{encode_query_params(params)}"
end

pid = Process.spawn("features/fixtures/local_server.sh",
                    :out => '/dev/null', :err => '/dev/null')
Process.detach pid
sleep(4)

at_exit do
  begin
    Process.kill("HUP", pid)
  rescue
  end
end
