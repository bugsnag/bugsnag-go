require 'test_helper'

class SampleTest < Test::Unit::TestCase

  def test_running_an_ios_simulator_harness
    run_scenario("test/fixtures/ios-app")
  end

  def test_running_a_ruby_web_server_harness
    run_scenario("test/fixtures/sinatra-app")
  end

  def test_running_a_browser_harness
    run_scenario("test/fixtures/js-app")
  end

  def run_scenario fixture_path
    Dir.chdir(fixture_path) do
        Process.wait Process.spawn("bundle", "exec", "bugsnag-maze-runner")
        status = $?.exitstatus
        assert_equal(0, status, "Scenario failed: #{fixture_path}")
    end
  end
end
