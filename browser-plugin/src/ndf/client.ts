import { Settings } from '../shared/settings';

export class NewsDeframerClient {
  constructor(private config: Settings) {}

  async getDomains(): Promise<string[]> {
    const headers: Record<string, string> = {};
    if (this.config.username && this.config.password) {
      headers['Authorization'] = 'Basic ' + btoa(`${this.config.username}:${this.config.password}`);
    }

    const url = this.config.backendUrl.replace(/\/$/, '') + '/api/domains';

    // Use chrome.runtime.sendMessage to proxy via background script to avoid CORS/Mixed-Content
    if (typeof chrome !== 'undefined' && chrome.runtime && chrome.runtime.sendMessage) {
      return new Promise((resolve, reject) => {
        chrome.runtime.sendMessage({ type: 'PROXY_REQ', url, headers }, (response) => {
          if (chrome.runtime.lastError) {
            reject(new Error(chrome.runtime.lastError.message));
          } else if (response && response.ok) {
            resolve(response.data);
          } else {
            reject(new Error(response?.error || `API Error: ${response?.status}`));
          }
        });
      });
    }

    throw new Error('Extension context missing: Cannot proxy request');
  }
}
