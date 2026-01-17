import log from '../shared/logger';

export const handlePortal = () => {
  log.info('Portal page detected.');
  document.body.style.border = '15px solid green'; // Visual proof for portal
};
