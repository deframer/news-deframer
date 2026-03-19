#!/usr/bin/env bash

set -euo pipefail

workdir="$(cd "$(dirname "$0")/.." && pwd)"
cd "$workdir"

output=$(make -n android-remote-qr)
if [[ "$output" != *"scripts/android-remote-qr.sh"* || "$output" != *"ADB_WAIT_TIMEOUT="* || "$output" != *"ADB_CONNECT="* ]]; then
  printf 'unexpected output for android-remote-qr target\n%s\n' "$output" >&2
  exit 1
fi

output=$(make -n ios-preflight)
if [[ "$output" != *"scripts/ios-preflight.sh"* ]]; then
  printf 'unexpected output for ios-preflight target\n%s\n' "$output" >&2
  exit 1
fi

output=$(make -n ios-sim SIMULATOR="iPhone 16")
if [[ "$output" != *"scripts/ios-sim.sh \"iPhone 16\""* ]]; then
  printf 'unexpected output for ios-sim target\n%s\n' "$output" >&2
  exit 1
fi

output=$(make -n ios-sim-build SIMULATOR="iPhone 16")
if [[ "$output" != *"scripts/ios-sim-build.sh \"iPhone 16\" \"ios/build\""* ]]; then
  printf 'unexpected output for ios-sim-build target\n%s\n' "$output" >&2
  exit 1
fi

tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

cat >"$tmpdir/adb" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

count_file="${ADB_TEST_COUNT_FILE:?}"

case "$1" in
  start-server)
    printf 'started\n'
    ;;
  mdns)
    if [ "${2:-}" = "check" ]; then
      printf 'mdns daemon version [adb discovery 0.0.0]\n'
      exit 0
    fi
    if [ "${2:-}" != "services" ]; then
      printf 'unexpected adb mdns args: %s\n' "$*" >&2
      exit 1
    fi
    count=0
    if [ -f "$count_file" ]; then
      count=$(cat "$count_file")
    fi
    count=$((count + 1))
    printf '%s' "$count" >"$count_file"
    if [ "$count" -eq 1 ]; then
      printf 'List of discovered mdns services\n'
      printf 'adb-test _adb-tls-pairing._tcp. 192.168.0.2:12345\n'
    else
      printf 'List of discovered mdns services\n'
      printf 'adb-test _adb-tls-connect._tcp. 192.168.0.2:23456\n'
    fi
    ;;
  pair)
    printf 'PAIRED %s %s\n' "$2" "$3"
    ;;
  connect)
    printf 'CONNECTED %s\n' "$2"
    ;;
  devices)
    printf 'List of devices attached\n192.168.0.2:23456 device\n'
    ;;
  *)
    printf 'unexpected adb args: %s\n' "$*" >&2
    exit 1
    ;;
esac
EOF
chmod +x "$tmpdir/adb"

cat >"$tmpdir/qrencode" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
cat >/dev/null
printf 'QR\n'
EOF
chmod +x "$tmpdir/qrencode"

count_file="$tmpdir/count"
set +e
output=$(PATH="$tmpdir:$PATH" ADB_TEST_COUNT_FILE="$count_file" ADB_WAIT_TIMEOUT=2 bash ./scripts/android-remote-qr.sh 2>&1)
status=$?
set -e

if [[ $status -eq 0 || "$output" != *"ADB_CONNECT is required."* || "$output" != *"Wireless debugging -> IP address & Port"* ]]; then
  printf 'unexpected output for missing ADB_CONNECT\n%s\n' "$output" >&2
  exit 1
fi

output=$(PATH="$tmpdir:$PATH" ADB_TEST_COUNT_FILE="$count_file" ADB_WAIT_TIMEOUT=2 ADB_CONNECT=192.168.0.2:23456 bash ./scripts/android-remote-qr.sh)

if [[ "$output" != *"ADB wireless QR payload:"* || "$output" != *"WIFI:T:ADB;S:ADB_WIFI_"* || "$output" != *"ADB_QR_PASSWORD="* || "$output" != *"QR"* || "$output" != *"Pairing with 192.168.0.2:12345..."* || "$output" != *"PAIRED 192.168.0.2:12345"* || "$output" != *"Using connect address 192.168.0.2:23456."* || "$output" != *"Connecting to 192.168.0.2:23456..."* || "$output" != *"CONNECTED 192.168.0.2:23456"* ]]; then
  printf 'unexpected output for android-remote-qr script\n%s\n' "$output" >&2
  exit 1
fi

printf 'Makefile QR tests passed\n'
