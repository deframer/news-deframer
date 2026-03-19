#!/usr/bin/env bash

set -euo pipefail

timeout_seconds="${ADB_WAIT_TIMEOUT:-90}"
pair_timeout_seconds="${ADB_PAIR_TIMEOUT:-20}"
connect_timeout_seconds="${ADB_CONNECT_TIMEOUT:-20}"
manual_connect_address="${ADB_CONNECT:-}"

if [ -z "$manual_connect_address" ]; then
  echo "ADB_CONNECT is required."
  echo "Find it on your phone in Developer options -> Wireless debugging -> IP address & Port."
  echo "Example: ADB_CONNECT=192.168.1.23:37139 make android-remote-qr"
  exit 1
fi

generate_name() {
  python -c "import secrets,string; chars=string.ascii_letters + string.digits; print('ADB_WIFI_' + ''.join(secrets.choice(chars) for _ in range(14)) + '-' + ''.join(secrets.choice(chars) for _ in range(6)))"
}

generate_password() {
  python -c "import secrets,string; chars=string.ascii_letters + string.digits; print(''.join(secrets.choice(chars) for _ in range(21)))"
}

find_mdns_address() {
  local service_name="$1"
  local deadline output address

  deadline=$((SECONDS + timeout_seconds))

  while [ "$SECONDS" -lt "$deadline" ]; do
    output="$(adb mdns services 2>/dev/null || true)"
    address="$(printf '%s\n' "$output" | awk -v service="$service_name" '$2 == service || $2 == service "." || $2 "." == service { print $3; exit }')"
    if [ -n "$address" ]; then
      printf '%s\n' "$address"
      return 0
    fi
    sleep 1
  done

  return 1
}

show_mdns_services() {
  echo "Current adb mdns services:"
  adb mdns services || true
}

name="$(generate_name)"
password="$(generate_password)"
payload="WIFI:T:ADB;S:${name};P:${password};;"

echo "ADB wireless QR payload:"
echo "$payload"
echo "ADB_QR_NAME=$name"
echo "ADB_QR_PASSWORD=$password"
echo "Scan from Developer options -> Wireless debugging -> Pair device with QR code"

adb start-server >/dev/null
adb mdns check || true

if command -v qrencode >/dev/null 2>&1; then
  printf '%s' "$payload" | qrencode -t ANSIUTF8
else
  echo "qrencode is not installed; payload printed above."
fi

echo "Waiting for Android wireless pairing service..."
pair_address="$(find_mdns_address "_adb-tls-pairing._tcp")" || {
  echo "Timed out waiting for pairing service after ${timeout_seconds}s."
  show_mdns_services
  exit 1
}

echo "Pairing with ${pair_address}..."
if ! timeout "${pair_timeout_seconds}s" adb pair "$pair_address" "$password"; then
  echo "Pairing did not complete within ${pair_timeout_seconds}s."
  show_mdns_services
  exit 1
fi

echo "Waiting for Android wireless connect step..."
connect_address="$manual_connect_address"
echo "Using connect address ${connect_address}."

echo "Connecting to ${connect_address}..."
if ! timeout "${connect_timeout_seconds}s" adb connect "$connect_address"; then
  echo "Connect did not complete within ${connect_timeout_seconds}s."
  show_mdns_services
  exit 1
fi
adb devices
