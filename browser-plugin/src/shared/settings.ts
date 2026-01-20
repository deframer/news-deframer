export interface Settings {
  backendUrl: string;
  username?: string;
  password?: string;
  enabled: boolean;
}

export const DEFAULT_BACKEND_URL = 'http://localhost:8080';

export const getSettings = (): Promise<Settings> => {
  return new Promise((resolve) => {
    chrome.storage.local.get(
      {
        backendUrl: DEFAULT_BACKEND_URL,
        username: '',
        password: '',
        enabled: true,
      },
      (items: Settings) => {
        resolve(items);
      }
    );
  });
};
