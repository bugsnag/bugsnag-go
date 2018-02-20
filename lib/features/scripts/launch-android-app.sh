#!/usr/bin/env bash

if [ -z "$APP_BUNDLE" ]; then
    echo APP_BUNDLE environment variable is not set
    exit 1
fi

if [ -z "$APP_ACTIVITY" ]; then
    echo APP_ACTIVITY environment variable is not set
    exit 1
fi

if [ -z "$EVENT_TYPE" ]; then
    echo EVENT_TYPE environment variable is not set
    exit 1
fi

if [ -z "$MOCK_API_PORT" ]; then
    echo MOCK_API_PORT environment variable is not set
    exit 1
fi

if [ -z "$BUGSNAG_API_KEY" ]; then
    echo BUGSNAG_API_KEY environment variable is not set
    exit 1
fi
echo "Launching MainActivity with '$EVENT_TYPE'"
adb shell am start -n "$APP_BUNDLE/$APP_ACTIVITY" \
    --es EVENT_TYPE "$EVENT_TYPE" \
    --es BUGSNAG_PORT "$MOCK_API_PORT" \
    --es BUGSNAG_API_KEY "$BUGSNAG_API_KEY"

echo "Ran Android Test case"
