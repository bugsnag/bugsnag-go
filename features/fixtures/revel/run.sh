#!/bin/sh

set -e

if [ -z "${USE_CODE_CONFIG}" ]
then
    $GOPATH/bin/revel run test configfile
else
    $GOPATH/bin/revel run test
fi
