#!/usr/bin/env bash

set -euo pipefail

project="ios/NewsDeframerMobile.xcodeproj"
mkdir -p ios/build
rm -f ios/build/xcode-preflight.log

if xcodebuild -list -project "$project" >/dev/null 2>ios/build/xcode-preflight.log; then
  exit 0
fi

if grep -q "runFirstLaunch" ios/build/xcode-preflight.log; then
  echo "Xcode is not fully initialized for simulator builds."
  echo "Run this once with admin privileges, then retry:"
  echo "  sudo xcodebuild -runFirstLaunch"
  echo
  echo "Raw xcodebuild error:"
fi

cat ios/build/xcode-preflight.log
exit 1
