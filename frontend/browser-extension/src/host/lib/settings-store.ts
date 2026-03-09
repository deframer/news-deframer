import { DEFAULT_SETTINGS, Settings } from '@frontend-shared/settings';

export type { Settings } from '@frontend-shared/settings';
export { DEFAULT_BACKEND_URL } from '@frontend-shared/settings';

export const getSettings = (): Promise<Settings> => {
  return new Promise((resolve) => {
    chrome.storage.local.get(DEFAULT_SETTINGS, (items: Settings) => {
      resolve(items);
    });
  });
};
