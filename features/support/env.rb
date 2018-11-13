require 'fileutils'
require 'socket'
require 'timeout'

testBuildFolder = 'features/fixtures/testbuild'

FileUtils.rm_rf(testBuildFolder)
Dir.mkdir testBuildFolder

# Copy the existing dir
`find . -name '*.go' \
        -not -path "./examples/*" \
        -not -path "./testutil/*" \
        -not -path "./features/*" \
        -not -name '*_test.go' | cpio -pdm #{testBuildFolder}`

def port_open?(ip, port, seconds=1)
  Timeout::timeout(seconds) do
    begin
      TCPSocket.new(ip, port).close
      true
    rescue Errno::ECONNREFUSED, Errno::EHOSTUNREACH
      false
    end
  end
rescue Timeout::Error
  false
end

def wait_for_port(port)
  max_attempts = 10
  attempts = 0
  up = false
  until (attempts >= max_attempts) || up
    attempts += 1
    up = port_open?("localhost", port)
    sleep 1
  end
  raise "Port not ready in time!" unless up
end
