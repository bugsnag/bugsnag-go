REVEL_PORT = 9020

# Replaces the revel configuration property with the given value.
# In the file config files all these properties have the format have the format
# bugsnag.propertyname=
# This method uses sed to replace this line with the given key-value pair for
# the given fixture.
def replace_revel_conf(fixture:, property_name:, property_value:)
  old = property_name + '='
  new = old + property_value
  full_path = "features/fixtures/#{fixture}/conf/app.conf"
  # 'sed' requires an extra flag for it to work properly on mac
  if (/darwin/ =~ RUBY_PLATFORM).nil?
    run_command("sed -i 's/\##{old}/#{new}/g' #{full_path}")
  else
    run_command("sed -i \"\" 's/\##{old}/#{new}/g' #{full_path}")
  end
end

def go_version_is_unsupported
  /go1.7/ =~ `go version`
end

Before do
  kill_apps
end

After do
  kill_apps
end

def kill_apps
  %w[gin martini revel negroni].each do |framework|
    begin
      run_command("killall -9 #{framework} || true")
    rescue SignalException
      # This can be raised in cases where the app isn't running
    end
  end
end
