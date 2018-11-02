#/bin/sh

set -e

cd $GOPATH/src/github.com/bugsnag/bugsnag-go/features/fixtures/negroni
go build
./negroni
