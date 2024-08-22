Before do
  Maze.config.enforce_bugsnag_integrity = false
  $address = nil
  $api_key = "166f5ad3590596f9aa8d601ea89af845"
  steps %(
    When I configure the base endpoint
  )
end