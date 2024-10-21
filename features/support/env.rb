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
  validator.validate_header('bugsnag-sent-at') { |value| Date.iso8601(value) }

  validator.element_has_value('notifier.name', 'Bugsnag Go')
  validator.each_element_exists(['notifier.url', 'notifier.version', 'events'])
  validator.each_event_contains_each(['severity', 'severityReason.type', 'unhandled', 'exceptions'])
end

Maze.config.add_validator('session') do |validator|
  validator.validate_header('bugsnag-api-key') { |value| value.eql?($api_key) }
  validator.validate_header('content-type') { |value| value.eql?('application/json') }
  validator.validate_header('bugsnag-payload-version') { |value| value.eql?('1.0') }
  validator.validate_header('bugsnag-sent-at') { |value| Date.iso8601(value) }

  validator.element_has_value('notifier.name', 'Bugsnag Go')
  validator.each_element_exists(['notifier.url', 'notifier.version', 'app', 'device'])
end
