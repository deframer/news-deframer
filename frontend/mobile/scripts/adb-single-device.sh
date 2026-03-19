#!/usr/bin/env bash

set -euo pipefail

devices="$(adb devices | awk 'NR>1 && $2=="device" {print $1}')"
count="$(printf '%s\n' "$devices" | awk 'NF' | wc -l)"

if [ "$count" -eq 0 ]; then
  echo "No adb device connected. Connect one device or pass ADB_DEVICE=<ip:port>."
  exit 1
fi

if [ "$count" -ne 1 ]; then
  echo "Multiple adb devices detected. Pass ADB_DEVICE=<ip:port>."
  adb devices
  exit 1
fi

device="$(printf '%s\n' "$devices" | awk 'NF {print; exit}')"
adb -s "$device" "$@"
