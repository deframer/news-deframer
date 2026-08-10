import { getDomain } from 'tldts';

import { Theme } from './theme';

export interface Settings {
  backendUrl: string;
  username?: string;
  password?: string;
  enabled: boolean;
  theme: Theme;
  language?: string;
  selectedDomains?: string[];
}

export const CONNECTION_TIMEOUT_MS = 5000;

export const getSettings = (): Promise<Settings> => {
  return new Promise((resolve) => {
    chrome.storage.local.get(
      {
        backendUrl: '',
        username: '',
        password: '',
        enabled: false,
        theme: 'system',
        selectedDomains: [],
      },
      (items: Settings) => {
        resolve(items);
      }
    );
  });
};

const normalizeDomain = (value: string) => getDomain(value) || value;

export const setSelectedDomains = (selectedDomains: string[]): Promise<void> => {
  return new Promise((resolve) => {
    chrome.storage.local.set({ selectedDomains }, () => resolve());
  });
};

export const removeStallDomainsFromSelection = async (knownDomains: string[]): Promise<void> => {
  const settings = await getSettings();
  const selectedDomains = settings.selectedDomains || [];
  const knownDomainSet = new Set(knownDomains.map(normalizeDomain));
  const nextSelectedDomains = selectedDomains.filter((domain) => knownDomainSet.has(domain));

  if (nextSelectedDomains.length === selectedDomains.length) {
    return;
  }

  await setSelectedDomains(nextSelectedDomains);
};
