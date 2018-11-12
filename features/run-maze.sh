#!/bin/sh

set -e

bundle install

function plain() {
    echo "Running maze tests for plain Go apps"
    bundle exec bugsnag-maze-runner features/plain_features
}

function nethttp() {
    echo "Running maze tests for net/http"
    bundle exec bugsnag-maze-runner features/net_http_features
}

function martini() {
    echo "Running maze tests for Martini"
    bundle exec bugsnag-maze-runner features/martini_features
}

function gin() {
    echo "Running maze tests for Gin version $1"
    GIN_VERSION=$1 bundle exec bugsnag-maze-runner features/gin_features
}

function negroni() {
    echo "Running maze tests for Negroni version $1"
    NEGRONI_VERSION=$1 bundle exec bugsnag-maze-runner features/negroni_features
}

function revel() {
    echo "Running maze tests for Revel version $1 and CMD version $2"
    REVEL_VERSION=$1 REVEL_CMD_VERSION=$2 bundle exec bugsnag-maze-runner features/revel_features --verbose
}

function run_all() {
    export GO_VERSION=$1
    echo "Running maze tests for Go version $GO_VERSION"

    plain
    nethttp
    martini

    gin v1.3.0
    gin v1.2
    gin v1.1
    gin v1.0

    negroni v1.0.0
    negroni v0.3.0
    negroni v0.2.0

    # Revel requires Go 1.8 or higher
    if [ $1 != "1.7" ] ; then
        revel v0.21.0 v0.21.1
        revel v0.20.0 v0.20.2
        revel v0.19.1 v0.19.0
        revel v0.18.0 v0.18.0
    fi
}

run_all 1.11
run_all 1.10
run_all 1.9
run_all 1.8
run_all 1.7
