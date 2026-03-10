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
make ios
```

## Local builds

Android:

```sh
make android-apk
make android-apk-release
make android-bundle-release
```

iOS:

```sh
make ios-build
```

## Android release signing

To sign release Android builds with your own keystore:

- copy `android/keystore.properties.example` to `android/keystore.properties`
- create or place your keystore file inside `frontend/mobile/android/app/`
- update the values in `android/keystore.properties`

If `android/keystore.properties` does not exist, release builds fall back to the debug keystore.
