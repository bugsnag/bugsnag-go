#!/bin/sh

set -e

# Use the [configfile] profile declared in conf/app.conf if the USE_PROPERTIES_FILE_CONFIG is set
if [ -z "${USE_PROPERTIES_FILE_CONFIG}" ]
then
    $GOPATH/bin/revel run test
else
    $GOPATH/bin/revel run test configfile
fi
