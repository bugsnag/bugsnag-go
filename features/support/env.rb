require 'fileutils'

testBuildFolder = 'features/fixtures/testbuild'

Dir.mkdir testBuildFolder

# Copy the existing air
`find . -name '*.go' \
        -not -path "./examples/*" \
        -not -path "./testutil/*" \
        -not -path "./features/*" \
        -not -name '*_test.go' | cpio -pdm #{testBuildFolder}`


# Scenario hooks
Before do
# Runs before every Scenario
end

After do
# Runs after every Scenario
end

at_exit do
# Runs when the test run is completed
  FileUtils.rm_rf(testBuildFolder)
end