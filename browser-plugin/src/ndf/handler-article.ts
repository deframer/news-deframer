import log from '../shared/logger';

export const handleArticle = () => {
  log.info('Article page detected.');
  document.body.style.border = '15px solid blue'; // Visual proof for article
};
