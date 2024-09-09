Before do
  Maze.config.enforce_bugsnag_integrity = false
  $address = nil
  $api_key = "166f5ad3590596f9aa8d601ea89af845"
  steps %(
    When I configure the base endpoint
  )
end

Maze.config.add_validator('error') do |validator|
  pp validator.headers
  validator.validate_header('Bugsnag-Api-Key') { |value| value.eql?($api_key) }
  validator.validate_header('Content-Type') { |value| value.eql?('application/json') }
  validator.validate_header('Bugsnag-Payload-Version') { |value| value.eql?('4') }
  validator.validate_header('Bugsnag-Sent-At') do |value|
    begin
      Date.iso8601(value)
    rescue Date::Error
      validator.success = false
      validator.errors << "bugsnag-sent-at header was expected to be an ISO 8601 date, but was '#{value}'"
    end
  end

  notifier_name = Maze::Helper.read_key_path(validator.body, 'notifier.name')
  if notifier_name.nil? || !notifier_name.eql?('Bugsnag Go')
    validator.success = false
    validator.errors << "Notifier name in body was expected to be 'Bugsnag Go', but was '#{notifier_name}'"
  end

  error_elements = ['notifier.url', 'notifier.version', 'events']
  error_elements_present = error_elements.all? do |element_key|
    element = Maze::Helper.read_key_path(validator.body, element_key)
    !element.nil? && (!element.is_a?(Array) || !element.empty?)
  end

  unless error_elements_present
    validator.success = false
    validator.errors << "Not all of the error payload elements were present"
  end

  event_elements = ['severity', 'severityReason.type', 'unhandled', 'exceptions']
  events = Maze::Helper.read_key_path(validator.body, 'events')
  event_elements_present = events.all? do |event|
    event_elements.all? do |element_key|
      element = Maze::Helper.read_key_path(event, element_key)
      !element.nil?
    end
  end

  unless event_elements_present
    validator.success = false
    validator.errors << "Not all of the event elements were present"
  end
end

Maze.config.add_validator('session') do |validator|
  validator.validate_header('Bugsnag-Api-Key') { |value| value.eql?($api_key) }
  validator.validate_header('Content-Type') { |value| value.eql?('application/json') }
  validator.validate_header('Bugsnag-Payload-Version') { |value| value.eql?('1.0') }
  validator.validate_header('Bugsnag-Sent-At') do |value|
    begin
      Date.iso8601(value)
    rescue Date::Error
      validator.success = false
      validator.errors << "bugsnag-sent-at header was expected to be an ISO 8601 date, but was '#{value}'"
    end
  end

  notifier_name = Maze::Helper.read_key_path(validator.body, 'notifier.name')
  if notifier_name.nil? || !notifier_name.eql?('Bugsnag Go')
    validator.success = false
    validator.errors << "Notifier name in body was expected to be 'Bugsnag Go', but was '#{notifier_name}'"
  end

  session_elements = ['notifier.url', 'notifier.version', 'events']
  session_elements_present = session_elements.all? do |element_key|
    element = Maze::Helper.read_key_path(validator.body, element_key)
    !element.nil? && (!element.is_a?(Array) || !element.empty?)
  end

  unless session_elements_present
    validator.success = false
    validator.errors << "Not all of the session payload elements were present"
  end
end
