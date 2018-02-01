require 'test_helper'
require 'fileutils'

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

  def test_comparing_requests_to_json_files
    run_scenario("test/fixtures/comparison")
  end

  def test_init_command
    fixture_dir = "test/fixtures/init-test"
    FileUtils.rm_rf fixture_dir
    FileUtils.mkdir_p fixture_dir
    Dir.chdir(fixture_dir) do
      Process.wait Process.spawn("../../../bin/bugsnag-maze-runner", "init", "--skip-install")
      status = $?.exitstatus
      assert_equal(0, status, "Running init failed")
      open("Gemfile", "w") do |file|
        file.puts <<-CONTENTS
source 'https://rubygems.org'
gem 'bugsnag-maze-runner', :path => '../../..'
CONTENTS
      end
      system("bundle", "install")
    end
    run_scenario(fixture_dir)
    FileUtils.rm_rf fixture_dir
  end

  def run_scenario fixture_path
    Dir.chdir(fixture_path) do
        Process.wait Process.spawn("bundle", "exec", "bugsnag-maze-runner")
        status = $?.exitstatus
        assert_equal(0, status, "Scenario failed: #{fixture_path}")
    end
  end
end
