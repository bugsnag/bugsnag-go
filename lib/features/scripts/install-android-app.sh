#!/usr/bin/env bash

if [ -z "$APP_BUNDLE" ]; then
    echo APP_BUNDLE environment variable is not set
    exit 1
fi
if [ -z "$APK_PATH" ]; then
    echo APK_PATH environment variable is not set
    exit 1
fi

echo "Uninstalling '$APP_BUNDLE'"
adb uninstall "$APP_BUNDLE"
echo "Installing '$APK_PATH'"
adb install "$APK_PATH"
