import { ThemeMode } from '../theme';
import { storageService } from './storageService';

export interface Settings {
  backendUrl: string;
  username: string;
  password: string;
  theme: ThemeMode;
  language: string;
}

const SETTINGS_KEY = 'news-deframer.mobile.settings';
export const DEFAULT_BACKEND_URL = 'https://api.deframer.org';
export const CONNECTION_TIMEOUT_MS = 5000;

export const DEFAULT_SETTINGS: Settings = {
  backendUrl: DEFAULT_BACKEND_URL,
  username: '',
  password: '',
  theme: 'system',
  language: 'default',
};

export const settingsService = {
  async loadSettings(): Promise<{ settings: Settings; configured: boolean }> {
    const stored = await storageService.getJson<Settings>(SETTINGS_KEY);
    return {
      settings: { ...DEFAULT_SETTINGS, ...(stored ?? {}) },
      configured: stored !== null,
    };
  },

  async saveSettings(settings: Settings): Promise<void> {
    await storageService.setJson(SETTINGS_KEY, settings);
  },

  hasRequiredConfiguration(settings: Settings): boolean {
    return (settings.backendUrl || '').trim().length > 0;
  },
};
