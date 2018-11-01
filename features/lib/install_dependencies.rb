unless ENV['MAZE_SKIP_INSTALL']
  Dir.chdir('features/fixtures') do
    run_command('./build.sh')
  end
else
  puts 'SKIPPING DEPENDENCY INSTALLATION'
end
