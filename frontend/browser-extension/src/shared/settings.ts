import { Theme } from './theme';

export interface Settings {
  backendUrl: string;
  username?: string;
  password?: string;
  enabled: boolean;
  theme: Theme;
  searchEngineUrl?: string;
  language?: string;
  selectedDomains?: string[];
}

export const DEFAULT_BACKEND_URL = 'https://api.deframer.org';
export const CONNECTION_TIMEOUT_MS = 5000;

export const getSettings = (): Promise<Settings> => {
  return new Promise((resolve) => {
    chrome.storage.local.get(
      {
        backendUrl: DEFAULT_BACKEND_URL,
        username: '',
        password: '',
        enabled: false,
        theme: 'system',
        searchEngineUrl: 'https://search.brave.com',
        selectedDomains: [],
      },
      (items: Settings) => {
        resolve(items);
      }
    );
  });
};
