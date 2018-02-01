# coding: utf-8
lib = File.expand_path('../lib', __FILE__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)
require 'version'

Gem::Specification.new do |spec|
  spec.name    = 'bugsnag-maze-runner'
  spec.version = BugsnagMazeRunner::VERSION
  spec.authors = ['Delisa Mason']
  spec.email   = ['iskanamagus@gmail.com']
  spec.required_ruby_version = '>= 2.0.0'
  spec.description =
    %q{
    Automation steps and mock server to validate request payloads
    response.
    }
  spec.summary = 'Bugsnag API request validation harness'
  spec.license = 'MIT'
  spec.require_paths = ["lib"]
  spec.files = [
    'bin/bugsnag-maze-runner',
    'bin/bugsnag-print-load-paths',
    'lib/features/steps/automation_steps.rb',
    'lib/features/steps/error_reporting_steps.rb',
    'lib/features/steps/request_assertion_steps.rb',
    'lib/features/support/compare.rb',
    'lib/features/support/env.rb',
    'lib/version.rb',
  ]
  spec.executables = spec.files.grep(%r{^bin/}) { |f| File.basename(f) }

  spec.add_dependency "cucumber", "~> 3.1.0"
  spec.add_dependency "test-unit", "~> 3.2.0"
  spec.add_dependency "rack", "~> 2.0.0"
  spec.add_dependency "minitest", "~> 5.0"
end

