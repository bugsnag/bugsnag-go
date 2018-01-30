require 'rack'
require 'open3'
require 'webrick'

MOCK_API_PORT = 9291

Before do
  stop_server
  start_server
  @requests = []
  @script_env = {'MOCK_API_PORT' => "#{MOCK_API_PORT}"}
  @pids = []
end

After do |scenario|
  stop_server
  kill_script
  # TODO: if scenario fails, print script output
end

# Run each command synchronously, printing output only in the event of failure
# and exiting the program
def run_required_commands command_arrays
  command_arrays.each do |args|
    out, err, ps = Open3.capture3(*args)
    unless ps.exitstatus == 0
      puts out.read
      puts err.read
      exit(1)
    end
  end
end

def set_script_env key, value
  @script_env[key] = value
end

def run_script script_path
  load_path = File.join(Dir.pwd, script_path)
  pid = Process.spawn(@script_env, load_path,
                     :err => '/dev/null', :out => '/dev/null')
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
  @requests[request_index][:body]["events"][event_index]
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

def start_server
  @thread = Thread.new do
    server = WEBrick::HTTPServer.new(
      :Port => MOCK_API_PORT,
      Logger: WEBrick::Log.new("/dev/null"),
      AccessLog: [],
    )
    server.mount_proc '/' do |req, res|
      @requests << {body: JSON.load(req.body()), request:req}
      res.status = 200
      res.body = 'OK'
    end
    server.start
  end
end

def stop_server
  @thread.kill if @thread and @thread.alive?
  @thread = nil
end
