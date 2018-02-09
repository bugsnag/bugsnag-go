require 'rack'
require 'open3'
require 'webrick'

MOCK_API_PORT = 9291

Before do
  stored_requests.clear
  @script_env = {'MOCK_API_PORT' => "#{MOCK_API_PORT}"}
  @pids = []
end

After do |scenario|
  kill_script
  # TODO: if scenario fails, print script output
end

# Run each command synchronously, printing output only in the event of failure
# and exiting the program
def run_required_commands command_arrays
  command_arrays.each do |args|
    if ENV['VERBOSE']
      command = args.join(' ')
      puts "Running '#{command}'"
      `#{command}`
    else
      out_reader, out_writer = IO.pipe
      err_reader, err_writer = IO.pipe
      pid = Process.spawn(@script_env || {}, args.join(' '),
                          :out => out_writer.fileno,
                          :err => err_writer.fileno)
      Process.waitpid(pid, 0)
      unless $?.exitstatus == 0
        puts "Script failed (#{args}):"
        puts out_reader.gets
        puts err_reader.gets
        exit(1)
      end
    end
  end
end

def encode_query_params hash
  URI.encode_www_form hash
end

def set_script_env key, value
  @script_env[key] = value
end

def run_script script_path
  load_path = File.join(Dir.pwd, script_path)
  if ENV['VERBOSE']
    puts "Running '#{load_path}'"
    pid = Process.spawn(@script_env, load_path)
  else
    pid = Process.spawn(@script_env, load_path, :out => '/dev/null', :err => '/dev/null')
  end
  Process.detach(pid)
  @pids << pid
end

def kill_script
  @pids.each {|p|
    begin
    Process.kill("HUP", p)
    rescue Errno::ESRCH
    end
  }
end

def load_event request_index=0, event_index=0
  stored_requests[request_index][:body]["events"][event_index]
end

def stored_requests
  $requests ||= []
end

def read_key_path hash, key_path
  value = hash
  key_path.split('.').each do |key|
    if key =~ /^(\d+)$/
      key = key.to_i
      if value.length > key
        value = value[key.to_i]
      else
        return nil
      end
    else
      if value.keys.include? key
        value = value[key]
      else
        return nil
      end
    end
  end
  value
end


class Servlet < WEBrick::HTTPServlet::AbstractServlet
  def do_POST request, response
    stored_requests << {body: JSON.load(request.body()), request:request}
    response.header['Access-Control-Allow-Origin'] = '*'
    response.status = 200
  end

  def do_OPTIONS request, response
    response.header['Access-Control-Allow-Origin'] = '*'
    response.header['Access-Control-Allow-Methods'] = 'POST, OPTIONS'
    response.header['Access-Control-Allow-Headers'] = 'Origin,Content-Type,Bugsnag-Sent-At,Bugsnag-Api-Key,Bugsnag-Payload-Version,Accept'
    response.status = 200
  end
end

def start_server
  @thread = Thread.new do
    server = WEBrick::HTTPServer.new(
      :Port => MOCK_API_PORT,
      Logger: WEBrick::Log.new("/dev/null"),
      AccessLog: [],
    )
    server.mount '/', Servlet
    server.start
  end
end

def stop_server
  @thread.kill if @thread and @thread.alive?
  @thread = nil
end

start_server

at_exit do
  stop_server
end

