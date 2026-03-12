import { BackHandler, Keyboard, Platform } from 'react-native';
import { useEffect } from 'react';

export const useAndroidBackHandler = (handler: () => boolean) => {
  useEffect(() => {
    if (Platform.OS !== 'android') {
      return;
    }

    const subscription = BackHandler.addEventListener('hardwareBackPress', () => {
      const keyboardVisible = typeof Keyboard.isVisible === 'function' ? Keyboard.isVisible() : false;

      if (keyboardVisible) {
        return false;
      }

      return handler();
    });

    return () => subscription.remove();
  }, [handler]);
};
