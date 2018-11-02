Before do
  kill_apps
end

After do
  kill_apps
end

def kill_apps
    begin
      run_command('killall martini || true')
    rescue SignalException
    end
    begin
      run_command('killall negroni || true')
    rescue SignalException
    end
    begin
      run_command('killall revel || true')
    rescue SignalException
    end
end
