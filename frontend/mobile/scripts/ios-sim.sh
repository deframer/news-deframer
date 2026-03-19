#!/usr/bin/env bash

set -euo pipefail

simulator="${1:-}"

if [ -z "$simulator" ]; then
  simulator="$(xcrun simctl list devices available | perl -ne 'if (/^\s+(iPhone[^\(]+?) \([0-9A-F-]+\) \((?:Shutdown|Booted)\)\s*$/) { print $1; exit }')"
fi

if [ -z "$simulator" ]; then
  echo "No available iPhone simulator found."
  echo "Use 'make ios-sim-list' to inspect devices or pass SIMULATOR=\"...\"."
  exit 1
fi

echo "Using simulator: $simulator"
bundle install
bundle exec pod install --project-directory=ios
npm run ios:sim -- --simulator="$simulator"
