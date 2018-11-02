Before do
  kill_apps
end

After do
  kill_apps
end

def kill_apps
  %w[gin martini revel negroni].each do |framework|
    begin
      run_command("killall #{framework} || true")
    rescue SignalException
      # This can be raised in cases where the app isn't running
    end
  end
end
