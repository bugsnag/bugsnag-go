#!/usr/bin/env bash

set -e

curentDir=$(pwd)

# Create an empty test build folder
rm -rf testbuild
mkdir testbuild

# Move to projects root folder
cd ../..

# Find all go files that are not testing files and copy into the test build folder
find . -name '*.go' -not -path "./examples/*" -not -path "./testutil/*" -not -path "./features/*"  -not -name '*_test.go' | cpio -pdm $curentDir/testbuild
cd $curentDir

docker build --build-arg GO_VERSION=1.11 -t bugsnag_go_test:1.11 -f base-image/Dockerfile .
docker build --build-arg GO_VERSION=1.10 -t bugsnag_go_test:1.10 -f base-image/Dockerfile .
docker build --build-arg GO_VERSION=1.9  -t bugsnag_go_test:1.9  -f base-image/Dockerfile .
docker build --build-arg GO_VERSION=1.8  -t bugsnag_go_test:1.8  -f base-image/Dockerfile .
docker build --build-arg GO_VERSION=1.7  -t bugsnag_go_test:1.7  -f base-image/Dockerfile .