#!/bin/sh

set -e

$GOPATH/bin/revel run github.com/bugsnag/bugsnag-go/features/fixtures/revel
