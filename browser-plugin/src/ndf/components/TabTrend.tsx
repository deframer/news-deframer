import { useTranslation } from 'react-i18next';
import { getDomain } from 'tldts';

import log from '../../shared/logger';

export const TabTrend = () => {
  const { t } = useTranslation();
  const rootDomain = getDomain(window.location.hostname) || window.location.hostname;

  log.info(`trend analysis - current domain is ${rootDomain}`);

  return (
    <div style={{ textAlign: 'center', padding: '3em 1em', color: 'var(--secondary-text)' }}>
      <h2>{t('trends.title') || 'Trends'}</h2>
      <p>{t('trends.coming_soon') || 'Trend analysis is coming soon.'}</p>
      <p>
        <a
          href="https://github.com/deframer/news-deframer-mining"
          target="_blank"
          rel="noopener noreferrer"
          style={{ color: 'inherit', textDecoration: 'underline' }}
        >
          https://github.com/deframer/news-deframer-mining
        </a>
      </p>
    </div>
  );
};