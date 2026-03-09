# Mobile TODO

## Upgrade back to Expo SDK 55

The app is currently pinned to Expo SDK 54 so it can open in the Android `Expo Go` app that was available during setup.

This is probably not because the phone is too old.

Android 13 is not the problem here. The more likely issue is that the Play Store `Expo Go` release on the device did not yet match Expo SDK 55 when we tried it.

### When to upgrade

Try upgrading back to SDK 55 later when:

- the Play Store `Expo Go` app has updated again
- Expo Go on the phone can open SDK 55 projects
- we no longer need the temporary SDK 54 fallback

### How to upgrade

From `frontend/mobile/`:

```bash
npx expo install expo@~55.0.0 react react-dom react-native react-native-web expo-status-bar @types/react
rm -rf node_modules package-lock.json
npm install
npx expo-doctor
npx expo start -c
```

### Expected SDK 55 package set

- `expo@~55.0.0`
- `expo-status-bar@~55.0.0`
- `react@19.2.0`
- `react-dom@19.2.0`
- `react-native@0.83.x`
- `react-native-web@^0.21.0`
- `@types/react@~19.2.x`

### Check before upgrading

- confirm the phone has the newest `Expo Go`
- uninstall and reinstall `Expo Go` if SDK mismatch errors continue
- restart the local dev server with `npx expo start -c`
