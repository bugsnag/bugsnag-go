# Configure app environment

run_required_commands([
  ["bundle", "install"],
  ["pod", "install"],
  ["features/fixtures/build-app.sh"],
])
