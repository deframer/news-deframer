import { DomainEntry } from '@frontend-shared/ndf-api-interfaces';
import type { Settings } from '@frontend-shared/settings';

import { invalidateDomainCache } from '../../domain-cache';
import { ProxyResponse } from '../../types';

export const testConnection = async (settings: Settings): Promise<{ connected: boolean; domains: DomainEntry[] }> => {
  const headers: HeadersInit = {};
  if (settings.username && settings.password) {
    headers.Authorization = 'Basic ' + btoa(`${settings.username}:${settings.password}`);
  }

  const url = settings.backendUrl.replace(/\/$/, '') + '/api/domains';
  const response = await new Promise<ProxyResponse>((resolve, reject) => {
    chrome.runtime.sendMessage({ type: 'PROXY_REQ', url, headers, timeout: 5000 }, (res: ProxyResponse) => {
      if (chrome.runtime.lastError) {
        reject(new Error(chrome.runtime.lastError.message));
        return;
      }
      resolve(res);
    });
  });

  if (!response?.ok) {
    return { connected: false, domains: [] };
  }

  await invalidateDomainCache();

  let data: unknown = response.data;
  if (typeof data === 'string') {
    try {
      data = JSON.parse(data);
    } catch {
      data = undefined;
    }
  }

  const domainList = Array.isArray(data)
    ? (data as DomainEntry[])
    : ((data as Record<string, unknown> | undefined)?.domains as DomainEntry[] | undefined) || ((data as Record<string, unknown> | undefined)?.data as DomainEntry[] | undefined) || [];

  return { connected: true, domains: domainList };
};
