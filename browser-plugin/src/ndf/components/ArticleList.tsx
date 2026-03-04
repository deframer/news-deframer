import { Fragment, useEffect, useState, useRef } from 'react';
import { useTranslation } from 'react-i18next';

import { getSettings } from '../../shared/settings';
import { AnalyzedArticle, DomainEntry, NewsDeframerClient } from '../client';
import { RatingBar } from './RatingBar';
import { Spinner } from './Spinner';

interface ArticleListProps {
  term: string;
  domain: DomainEntry;
  date?: string; // Specific date, if provided, means it's a fixed window
  days?: number; // Rolling window size (ignored if date is provided)
  titleOverride?: string;
  hideTitle?: boolean;
}

const ARTICLES_PER_PAGE = 10;

export const ArticleList = ({ term, domain, date, days, titleOverride, hideTitle = false }: ArticleListProps) => {
  const { t, i18n } = useTranslation();
  const [articles, setArticles] = useState<AnalyzedArticle[]>([]);
  const [loading, setLoading] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const paginationRef = useRef<HTMLDivElement>(null); // Ref for the pagination controls

  const isFixedDateMode = !!date;

  const renderTitle = () => {
    if (titleOverride) return titleOverride;
    return t('trends.articles', 'Articles');
  };

  useEffect(() => {
    const fetchArticles = async () => {
      if (!term || !domain.domain) {
        setArticles([]);
        setLoading(false);
        return;
      }
      setLoading(true);
      try {
        const settings = await getSettings();
        const client = new NewsDeframerClient(settings);
        const offset = (currentPage - 1) * ARTICLES_PER_PAGE;
        const limit = ARTICLES_PER_PAGE;

        const result = await client.getArticlesByTrend(domain.domain, term, date, days, offset, limit);

        if (currentPage > 1 && (!result || result.length === 0)) {
          setHasMore(false);
        } else {
          setArticles(result || []);
          setHasMore((result?.length || 0) === ARTICLES_PER_PAGE);
        }
      } catch (error) {
        console.error('Failed to fetch articles', error);
        setArticles([]);
        setHasMore(false);
      } finally {
        setLoading(false);
      }
    };
    fetchArticles();
  }, [term, domain, date, days, currentPage]);

  // Reset pagination when the core context changes
  useEffect(() => {
    setCurrentPage(1);
    setHasMore(true);
  }, [term, domain, date, days]);

  // Scroll to pagination controls when page changes
  useEffect(() => {
    if (paginationRef.current) {
      paginationRef.current.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
    }
  }, [currentPage]);

  const handleNextPage = () => {
    setCurrentPage(prev => prev + 1);
  };

  const handlePrevPage = () => {
    setCurrentPage(prev => Math.max(1, prev - 1));
  };

  const formatPubDate = (pubDate: string) => {
    return new Date(pubDate).toLocaleDateString(i18n.language, { month: 'short', day: 'numeric' });
  };

  return (
    <div className="article-list-container" style={{ padding: '16px', border: '1px solid var(--border-color)', marginTop: '16px', borderRadius: '8px', backgroundColor: 'var(--card-bg)' }}>
      {!hideTitle && (
        <h4 style={{ margin: '0 0 12px 0', color: 'var(--text-color)' }}>
          {renderTitle()}
        </h4>
      )}

      {loading && (
        <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '150px' }}>
          <Spinner />
        </div>
      )}

      {!loading && articles.length === 0 && (
        <p style={{ marginTop: '16px', textAlign: 'center', fontStyle: 'italic', color: 'var(--secondary-text)' }}>
          {t('trends.articles_no_data', { term })}
        </p>
      )}

      {!loading && articles.length > 0 && (
        <Fragment>
          <div className="article-table-wrapper" style={{ overflowX: 'auto' }}>
            <table className="article-table" style={{ width: '100%', borderCollapse: 'collapse', tableLayout: 'fixed' }}>
              <thead>
                <tr>
                  <th style={{ width: '70px', textAlign: 'left' }}>{t('trends.rating_caption', 'Rating')}</th>
                  {!isFixedDateMode && <th style={{ width: '80px', textAlign: 'left' }}>{t('trends.date_caption', 'Date')}</th>}
                  <th style={{ width: 'auto', textAlign: 'left' }}>
                    <div style={{ position: 'relative', display: 'inline-flex', alignItems: 'center' }}>
                      <span>{t('trends.article_caption', 'Article')}</span>
                      <span
                        title={t('trends.article_column_tooltip_content')}
                        style={{ marginLeft: '5px', cursor: 'help', fontSize: '0.8em', verticalAlign: 'super', color: 'var(--secondary-text)' }}
                      >
                        (i)
                      </span>
                      <div className="article-header-tooltip" style={{ position: 'absolute', bottom: '100%', left: '50%', transform: 'translateX(-50%)', marginBottom: '8px' }}>
                        {t('trends.article_column_tooltip_content')}
                      </div>
                    </div>
                  </th>
                </tr>
              </thead>
              <tbody>
                {articles.map((article, index) => (
                  <tr key={article.url + index} style={{ borderTop: '1px solid var(--border-color)' }}>
                    <td style={{ padding: '8px 5px', verticalAlign: 'middle' }}>
                      {article.rating !== undefined && (
                        <RatingBar
                          value={article.rating}
                          label={undefined}
                          id={`article-rating-${index}`}
                        />
                      )}
                    </td>
                    {!isFixedDateMode && <td style={{ padding: '8px 5px', verticalAlign: 'middle', fontSize: '0.9em', color: 'var(--secondary-text)' }}>{formatPubDate(article.pub_date)}</td>}
                    <td style={{ padding: '8px 5px', verticalAlign: 'middle' }}>
                      <a
                        href={article.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        style={{ color: 'var(--accent-color)', textDecoration: 'underline', fontSize: '0.95em', lineHeight: '1.3' }}
                      >
                        {article.title || article.url}
                      </a>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <div ref={paginationRef} className="pagination-controls" style={{ display: 'flex', alignItems: 'center', marginTop: '15px' }}>
            <button
              className="pagination-btn"
              onClick={handlePrevPage}
              disabled={currentPage === 1}
            >
              {t('trends.pagination_prev', 'Previous')}
            </button>
            <span style={{ fontSize: '0.9em', color: 'var(--secondary-text)' }}>
              Page {currentPage}
            </span>
            <button
              className="pagination-btn"
              onClick={handleNextPage}
              disabled={!hasMore}
            >
              {t('trends.pagination_next', 'Next')}
            </button>
          </div>
        </Fragment>
      )}
    </div>
  );
};