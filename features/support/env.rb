Before do
  Maze.config.enforce_bugsnag_integrity = false
  $address = nil
  $api_key = "166f5ad3590596f9aa8d601ea89af845"
  steps %(
    When I configure the base endpoint
  )
end

Maze.config.add_validator('error') do |validator|
  validator.validate_header('bugsnag-api-key') { |value| value.eql?($api_key) }
  validator.validate_header('content-type') { |value| value.eql?('application/json') }
  validator.validate_header('bugsnag-payload-version') { |value| value.eql?('4') }
  validator.validate_header('bugsnag-sent-at') do |value|
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

  ['notifier.url', 'notifier.version', 'events'].each do |element_key|
    element = Maze::Helper.read_key_path(validator.body, element_key)
    if element.nil? || (element.is_a?(Array) && element.empty?)
      validator.success = false
      validator.errors << "Required error element #{element_key} was not present"
    end
  end

  events = Maze::Helper.read_key_path(validator.body, 'events')
  events.each_with_index do |event, index|
    ['severity', 'severityReason.type', 'unhandled', 'exceptions'].each do |element_key|
      element = Maze::Helper.read_key_path(event, element_key)
      if element.nil? || (element.is_a?(Array) && element.empty?)
        validator.success = false
        validator.errors << "Required event element #{element_key} was not present in event #{index}"
      end
    end
  end
end

Maze.config.add_validator('session') do |validator|
  validator.validate_header('bugsnag-api-key') { |value| value.eql?($api_key) }
  validator.validate_header('content-type') { |value| value.eql?('application/json') }
  validator.validate_header('bugsnag-payload-version') { |value| value.eql?('1.0') }
  validator.validate_header('bugsnag-sent-at') do |value|
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

  ['notifier.url', 'notifier.version', 'events'].each do |element_key|
    element = Maze::Helper.read_key_path(validator.body, element_key)
    if element.nil? || (element.is_a?(Array) && element.empty?)
      validator.success = false
      validator.errors << "Required session element #{element_key} was not present"
    end
  end
end
