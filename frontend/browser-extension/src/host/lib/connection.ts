import { DomainEntry } from '../../ndf/client';
import { invalidateDomainCache } from '../../shared/domain-cache';
import i18n from '../../shared/i18n';
import { CONNECTION_TIMEOUT_MS, Settings } from '../../shared/settings';
import { ProxyResponse } from '../../shared/types';

export const testConnection = async (settings: Settings): Promise<{ connected: boolean; domains: DomainEntry[]; errorMessage?: string }> => {
  const headers: HeadersInit = {};
  if (settings.username && settings.password) {
    headers.Authorization = 'Basic ' + btoa(`${settings.username}:${settings.password}`);
  }

  const trimmedUrl = settings.backendUrl?.trim() ?? '';
  if (!trimmedUrl) {
    return { connected: false, domains: [], errorMessage: i18n.t('options.invalid_server_url', 'Invalid server URL') };
  }

  try {
    new URL(trimmedUrl);
  } catch {
    return { connected: false, domains: [], errorMessage: i18n.t('options.invalid_server_url', 'Invalid server URL') };
  }
  const url = settings.backendUrl.replace(/\/$/, '') + '/api/domains';
  const response = await new Promise<ProxyResponse>((resolve, reject) => {
    const timeoutId = setTimeout(() => {
      reject(new Error('Connection test timed out'));
    }, CONNECTION_TIMEOUT_MS + 1000);

    chrome.runtime.sendMessage({ type: 'PROXY_REQ', url, headers, timeout: CONNECTION_TIMEOUT_MS }, (res: ProxyResponse) => {
      clearTimeout(timeoutId);
      if (chrome.runtime.lastError) {
        reject(new Error(chrome.runtime.lastError.message));
        return;
      }
      resolve(res);
    });
  });

  if (!response?.ok) {
    const errorMessage = typeof response.status === 'number'
      ? `HTTP ${response.status}`
      : response.error;
    return { connected: false, domains: [], errorMessage };
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
