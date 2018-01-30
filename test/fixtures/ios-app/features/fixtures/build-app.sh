#!/usr/bin/env bash

if [ ! -d "ios-app.xcworkspace" ]; then
    cd "$(dirname "$0")/.."
fi

osascript -e 'launch application "Simulator"'
sleep 4

xcrun xcodebuild \
  -scheme MazeRunner \
  -workspace ios-app.xcworkspace \
  -configuration Debug \
  -destination 'platform=iOS Simulator,name=iPhone 8,OS=11.2' \
  -derivedDataPath build \
  build

xcrun simctl install booted \
  build/Build/Products/Debug-iphonesimulator/MazeRunner.app
