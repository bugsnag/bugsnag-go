require 'fileutils'
require 'socket'
require 'timeout'

testBuildFolder = 'features/fixtures/testbuild'

FileUtils.rm_rf(testBuildFolder)
Dir.mkdir testBuildFolder

# Copy the existing dir
`find . -name '*.go' \
        -not -path "./examples/*" \
        -not -path "./testutil/*" \
        -not -path "./features/*" \
        -not -name '*_test.go' | cpio -pdm #{testBuildFolder}`
