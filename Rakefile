require "bundler/gem_tasks"
require "rake/testtask"

namespace :test do
  Rake::TestTask.new(:integration) do |t|
    t.libs << "test"
    t.libs << "lib"
    t.test_files = FileList['test/integration/{**/,}*_test.rb']
  end
  Rake::TestTask.new(:unit) do |t|
    t.libs << "test"
    t.libs << "lib"
    t.test_files = FileList['test/*_test.rb']
  end
  desc "Run all tests"
  task :all => [:unit, :integration]
end

task :default => :test
