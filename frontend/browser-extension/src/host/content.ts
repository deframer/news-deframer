import { getDomain } from 'tldts';

import * as ndf from '../ndf/index';
import log from '../shared/logger';
import { consumeBypassForCurrentTab } from '../shared/session-bypass';
import { getSettings, Settings } from '../shared/settings';

const getCurrentRootDomain = () => getDomain(window.location.hostname.replace(/:\d+$/, '')) || window.location.hostname;

const isCurrentDomainSelected = (settings: Settings) => {
  const selectedDomains = settings.selectedDomains || [];
  const currentDomain = getCurrentRootDomain();

  if (selectedDomains.length === 0) {
    log.debug(`NDF is not enabled for ${currentDomain}: no selected domains configured.`);
    return false;
  }

  const selected = selectedDomains.includes(currentDomain);
  if (!selected) {
    log.debug(`NDF is not enabled for ${currentDomain}: domain not selected.`);
  }

  return selected;
};

async function startNdf() {
  if (await consumeBypassForCurrentTab()) {
    log.debug('Bypass flag detected. Skipping NDF for this load.');
    return;
  }

  const settings = await getSettings();
  if (!settings.enabled) {
    log.debug('NDF is disabled by user settings.');
    return;
  }

  if (!isCurrentDomainSelected(settings)) {
    log.debug('Current domain is not selected.');
    return;
  }

  log.debug('Starting NDF...');
  await ndf.start(settings);
}

// Early exit if the content is not HTML.
// This is to prevent the script from running on non-HTML content like XML, JSON, etc.
if (document.contentType === 'text/html') {
  startNdf();
}
