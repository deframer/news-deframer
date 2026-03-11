import AsyncStorage from '@react-native-async-storage/async-storage';

export const storageService = {
  async getJson<T>(key: string): Promise<T | null> {
    const raw = await AsyncStorage.getItem(key);
    return raw ? (JSON.parse(raw) as T) : null;
  },

  async setJson<T>(key: string, value: T): Promise<void> {
    await AsyncStorage.setItem(key, JSON.stringify(value));
  },

  async getString(key: string): Promise<string | null> {
    return AsyncStorage.getItem(key);
  },

  async setString(key: string, value: string): Promise<void> {
    await AsyncStorage.setItem(key, value);
  },
};
