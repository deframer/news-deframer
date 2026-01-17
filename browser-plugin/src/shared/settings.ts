export interface Settings {
  backendUrl: string;
  username?: string;
  password?: string;
}

export const DEFAULT_BACKEND_URL = 'http://localhost:8080';

export const getSettings = (): Promise<Settings> => {
  return new Promise((resolve) => {
    chrome.storage.local.get(
      {
        backendUrl: DEFAULT_BACKEND_URL,
        username: '',
        password: '',
      },
      (items: Settings) => {
        resolve(items);
      }
    );
  });
};
