# News Deframer Mobile

Local-build React Native app for Android and iOS.

## Development

From `frontend/mobile/`:

```sh
make install
make start
make android
make web
```

For iOS on a Mac:

```sh
make pods
make ios-sim
```

If Xcode reports that a required plug-in failed to load or suggests `xcodebuild -runFirstLaunch`, finish Xcode setup once with admin privileges:

```sh
sudo xcodebuild -runFirstLaunch
```

Choose a specific simulator:

```sh
make ios-sim SIMULATOR="iPhone SE (3rd generation)"
```

If `SIMULATOR` is not set, the Make target picks the first available iPhone simulator automatically.

## Local builds

Android:

```sh
make android-apk
make android-apk-release
make android-bundle-release
```

iOS:

```sh
make ios-sim-build
```

List available iOS simulators:

```sh
make ios-sim-list
```

`make ios` and `make ios-build` still work as compatibility aliases for the simulator run/build flow.

## Android release signing

To sign release Android builds with your own keystore:

- copy `android/keystore.properties.example` to `android/keystore.properties`
- create or place your keystore file inside `frontend/mobile/android/app/`
- update the values in `android/keystore.properties`

If `android/keystore.properties` does not exist, release builds fall back to the debug keystore.
