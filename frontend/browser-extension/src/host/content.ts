import { getDomain } from 'tldts';

import { installPreemptiveDomGuard } from '../ndf/domGuard';
import * as ndf from '../ndf/index';
import log from '../shared/logger';
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
  if (sessionStorage.getItem('__ndf-bypass')) {
    log.debug('Bypass flag detected. Skipping NDF for this load.');
    sessionStorage.removeItem('__ndf-bypass');
    return;
  }

  // Install guard immediately to block popups, even before we know if we are enabled.
  // If we are disabled, we will reload the page without the guard.
  const guard = installPreemptiveDomGuard({ allowedIds: ['ndf-root', 'ndf-global-styles'] });

  const settings = await getSettings();
  if (!settings.enabled) {
    log.debug('NDF is disabled by user settings. Reloading without guard.');
    guard?.release();
    sessionStorage.setItem('__ndf-bypass', 'true');
    window.location.reload();
    return;
  }

  if (!isCurrentDomainSelected(settings)) {
    log.debug('Current domain is not selected. Reloading without guard.');
    guard?.release();
    sessionStorage.setItem('__ndf-bypass', 'true');
    window.location.reload();
    return;
  }

  log.debug('Starting NDF...');
  try {
    await ndf.start(settings);
  } finally {
    guard?.release();
  }
}

// Early exit if the content is not HTML.
// This is to prevent the script from running on non-HTML content like XML, JSON, etc.
if (document.contentType === 'text/html') {
  startNdf();
}
