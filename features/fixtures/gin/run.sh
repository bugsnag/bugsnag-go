#/bin/sh

set -e

cd $GOPATH/src/github.com/bugsnag/bugsnag-go/features/fixtures/gin
go build
./gin
