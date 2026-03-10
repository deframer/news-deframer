const getStorage = (): Storage | null => {
  if (typeof window === 'undefined') {
    return null;
  }

  return window.localStorage;
};

export const storageService = {
  async getJson<T>(key: string): Promise<T | null> {
    const storage = getStorage();
    const raw = storage?.getItem(key) ?? null;
    return raw ? (JSON.parse(raw) as T) : null;
  },

  async setJson<T>(key: string, value: T): Promise<void> {
    getStorage()?.setItem(key, JSON.stringify(value));
  },

  async getString(key: string): Promise<string | null> {
    return getStorage()?.getItem(key) ?? null;
  },

  async setString(key: string, value: string): Promise<void> {
    getStorage()?.setItem(key, value);
  },
};
