# Android over WLAN (ADB QR)

This is a WLAN-only guide.

## Prerequisites

- Android device and development machine on the same network
- ADB installed (`adb --version`)

## Wireless debugging pairing (Android 11+)

1. On the phone, open **Developer options -> Wireless debugging**.
2. Tap **Pair device with QR code**.
3. Note the **IP address & Port** shown on the same screen, e.g. `192.168.1.42:40991`.

From the `mobile` directory:

```bash
ADB_CONNECT=192.168.1.42:40991 make android-remote-qr
```

This flow:

- shows a valid Android wireless debugging QR code
- waits for the phone pairing service
- runs `adb pair`
- runs `adb connect <ip:port>` using `ADB_CONNECT`
- runs `adb devices`

After pairing, the first plausible device from `adb devices` is used for follow-up commands. You can override it with `ADB_DEVICE=<ip:port>`.

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

Optional: set `ADB_DEVICE=<ip:port>` if you want to force a specific device. Otherwise the first plausible entry from `adb devices` is used.

```bash
adb devices
npm run start
make android-metro-reverse
make android-apk-install
```

This rebuilds debug APK and installs it on the WLAN-connected device.

## Troubleshooting

- **`unable to connect`**: host and device are not on the same network, or firewall blocks the port.
- **Device is `offline`**: run `adb disconnect` and connect again.
- **Multiple devices connected**: pass `--deviceId <ip>:<port>`.
- **Connection lost after reboot**: reconnect with `ADB_CONNECT=<ip:port> make android-remote-qr`.
- **`Unable to load script from localhost:8081`**: start Metro (`npm run start`) and run `make android-metro-reverse ADB_DEVICE=<ip:port>`.
