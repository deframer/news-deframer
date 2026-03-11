# Android over WLAN (ADB Remote)

This is a WLAN-only guide.

## Prerequisites

- Android device and development machine on the same network
- ADB installed (`adb --version`)

## Method 1: Wireless debugging pairing (Android 11+)

1. On the phone, open **Developer options -> Wireless debugging**.
2. Tap **Pair device with pairing code**.
3. Note these values shown by Android:
   - Pair endpoint, e.g. `192.168.1.42:37123`
   - Pairing code, e.g. `123456`
   - Connect endpoint, e.g. `192.168.1.42:40991`

From the `mobile` directory:

```bash
make android-remote-pair ADB_PAIR=192.168.1.42:37123 ADB_PAIR_CODE=123456 ADB_HOST=192.168.1.42 ADB_PORT=40991
```

This runs:

- `adb pair <ip:pair-port> <pairing-code>`
- `adb connect <ip>:<connect-port>`
- `adb devices`

## Method 2: Direct connect (when device already listens on TCP)

If your device already has ADB over TCP enabled, run:

```bash
make android-remote ADB_HOST=192.168.1.42 ADB_PORT=5555
```

This runs:

- `adb connect <ip>:<port>`
- `adb devices`

`adb connect` opens a debug session from your computer to the phone over WLAN.

## Run the app over WLAN

Before launching a debug build, set Metro reverse tunnel:

```bash
make android-metro-reverse ADB_DEVICE=192.168.1.42:5555
```

Then run:

```bash
npm run android -- --deviceId 192.168.1.42:5555
```

If only one Android device is connected, this often works too:

```bash
npm run android
```

Recommended debug flow (prevents `localhost:8081` issues):

```bash
npm run start
make android-metro-reverse ADB_DEVICE=192.168.1.42:5555
npm run android -- --deviceId 192.168.1.42:5555
```

## After code changes

If the app gets stuck, keeps old behavior, or hot reload is unreliable, do a full rebuild + reinstall:

```bash
npm run start
make android-metro-reverse ADB_DEVICE=192.168.1.42:5555
make android-apk-install ADB_DEVICE=192.168.1.42:5555
```

This rebuilds debug APK and installs it on the WLAN-connected device.

## Troubleshooting

- **`unable to connect`**: host and device are not on the same network, or firewall blocks the port.
- **Device is `offline`**: run `adb disconnect` and connect again.
- **Multiple devices connected**: pass `--deviceId <ip>:<port>`.
- **Connection lost after reboot**: reconnect with the same WLAN method.
- **`Unable to load script from localhost:8081`**: start Metro (`npm run start`) and run `make android-metro-reverse ADB_DEVICE=<ip:port>`.
