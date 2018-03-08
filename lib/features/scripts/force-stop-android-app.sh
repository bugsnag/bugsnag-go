#!/usr/bin/env bash

if [ -z "$APP_BUNDLE" ]; then
    echo APP_BUNDLE environment variable is not set
    exit 1
fi

adb shell am force-stop $APP_BUNDLE
