require 'fileutils'
require 'socket'
require 'timeout'

testBuildFolder = 'features/fixtures/testbuild'

FileUtils.rm_rf(testBuildFolder)
Dir.mkdir testBuildFolder

# Copy the existing dir
`find . -name '*.go' -o -name 'go.sum' -o -name 'go.mod' \
        -not -path "./examples/*" \
        -not -path "./testutil/*" \
        -not -path "./v2/testutil/*" \
        -not -path "./features/*" \
        -not -name '*_test.go' | cpio -pdm #{testBuildFolder}`
