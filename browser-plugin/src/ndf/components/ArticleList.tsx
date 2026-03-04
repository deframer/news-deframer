import { useTranslation } from 'react-i18next';

import { DomainEntry } from '../client';

interface ArticleListProps {
  term: string;
  domain: DomainEntry;
  date?: string;
  days?: number;
  titleOverride?: string;
}

export const ArticleList = ({ term, domain, date, days, titleOverride }: ArticleListProps) => {
  const { t } = useTranslation();

  const formattedDate = date ? new Date(date).toLocaleDateString(undefined, { 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric' 
  }) : '';

  const renderTitle = () => {
    if (titleOverride) return titleOverride;
    if (date) return t('trends.articles_for_date', { date: formattedDate });
    return t('trends.articles', 'Articles');
  };

  return (
    <div className="article-list-container" style={{ padding: '16px', border: '1px solid var(--border-color)', marginTop: '16px', borderRadius: '8px', backgroundColor: 'var(--card-bg)' }}>
      <h4 style={{ margin: '0 0 12px 0', color: 'var(--text-color)' }}>
        {renderTitle()}
      </h4>
      <div style={{ fontSize: '0.85em', color: 'var(--secondary-text)', fontFamily: 'monospace', backgroundColor: 'var(--bg-color-secondary)', padding: '10px', borderRadius: '4px' }}>
        <div style={{ fontWeight: 'bold', marginBottom: '4px', borderBottom: '1px solid var(--border-color)' }}>DEBUG INFO:</div>
        <div>- Term: {term}</div>
        <div>- Domain: {domain.domain}</div>
        <div>- Specific Date: {date || 'None (Rolling)'}</div>
        <div>- Window (Days): {days || 'N/A'}</div>
      </div>
      <p style={{ marginTop: '16px', textAlign: 'center', fontStyle: 'italic', color: 'var(--secondary-text)' }}>
        {t('trends.articles_placeholder', 'Articles for "{{term}}" will be listed here.', { term })}
      </p>
    </div>
  );
};
