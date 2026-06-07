const BYPASS_KEY_PREFIX = 'ndf-bypass';

type BypassRequest =
  | { type: 'SET_BYPASS' }
  | { type: 'CONSUME_BYPASS' };

type BypassResponse = {
  ok: boolean;
  bypass?: boolean;
  error?: string;
};

export const getBypassStorageKey = (tabId: number, url: string) => {
  return `${BYPASS_KEY_PREFIX}:${tabId}:${new URL(url).origin}`;
};

const sendBypassMessage = (message: BypassRequest) => {
  return new Promise<BypassResponse>((resolve, reject) => {
    chrome.runtime.sendMessage(message, (response?: BypassResponse) => {
      if (chrome.runtime.lastError) {
        reject(new Error(chrome.runtime.lastError.message));
        return;
      }

      if (!response || response.ok === false) {
        reject(new Error(response?.error || 'Bypass request failed'));
        return;
      }

      resolve(response);
    });
  });
};

export const setBypassForCurrentTab = async () => {
  await sendBypassMessage({ type: 'SET_BYPASS' });
};

export const consumeBypassForCurrentTab = async () => {
  try {
    const response = await sendBypassMessage({ type: 'CONSUME_BYPASS' });
    return Boolean(response.bypass);
  } catch {
    return false;
  }
};
