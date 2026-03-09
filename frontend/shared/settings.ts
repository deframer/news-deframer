import { Theme } from './theme';

export interface Settings {
  backendUrl: string;
  username?: string;
  password?: string;
  enabled: boolean;
  theme: Theme;
  searchEngineUrl?: string;
  language?: string;
}

export const DEFAULT_BACKEND_URL = 'http://localhost:8080';

export const DEFAULT_SETTINGS: Settings = {
  backendUrl: DEFAULT_BACKEND_URL,
  username: '',
  password: '',
  enabled: false,
  theme: 'system',
  searchEngineUrl: 'https://search.brave.com',
};
