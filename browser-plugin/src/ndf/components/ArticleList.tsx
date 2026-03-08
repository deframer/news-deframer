import { Fragment, useEffect, useId, useRef, useState } from 'react';
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
  const articleTooltipId = useId();
  const [articles, setArticles] = useState<AnalyzedArticle[]>([]);
  const [loading, setLoading] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [isArticleTooltipOpen, setIsArticleTooltipOpen] = useState(false);
  const [articleTooltipPos, setArticleTooltipPos] = useState<{ x: number; y: number }>({ x: 0, y: 0 });
  const containerRef = useRef<HTMLDivElement>(null);
  const didInitPageRef = useRef(false);
  const hasLoadedOnceRef = useRef(false);

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
    didInitPageRef.current = false;
    hasLoadedOnceRef.current = false;
  }, [term, domain, date, days]);

  const scrollListIntoView = () => {
    if (typeof window === 'undefined') return;
    if (window.matchMedia('(max-width: 799px)').matches) return;
    containerRef.current?.scrollIntoView({ behavior: 'smooth', block: 'start' });
  };

  useEffect(() => {
    if (!didInitPageRef.current) {
      didInitPageRef.current = true;
      return;
    }
    scrollListIntoView();
  }, [currentPage]);

  useEffect(() => {
    if (loading) return;
    if (articles.length === 0) return;
    if (currentPage !== 1) return;
    if (hasLoadedOnceRef.current) return;
    hasLoadedOnceRef.current = true;
    scrollListIntoView();
  }, [loading, articles.length, currentPage]);

  const handleNextPage = () => {
    setCurrentPage(prev => prev + 1);
  };

  const handlePrevPage = () => {
    setCurrentPage(prev => Math.max(1, prev - 1));
  };

  const formatPubDate = (pubDate: string) => {
    return new Date(pubDate).toLocaleDateString(i18n.language, { month: 'short', day: 'numeric' });
  };

  const formatAuthors = (authors?: string[]) => {
    if (!authors || authors.length === 0) return '';
    return authors.join(', ');
  };

  const showArticleTooltip = (target: HTMLElement) => {
    const rect = target.getBoundingClientRect();
    setArticleTooltipPos({ x: rect.left + rect.width / 2, y: rect.top - 10 });
    setIsArticleTooltipOpen(true);
  };

  const hideArticleTooltip = () => {
    setIsArticleTooltipOpen(false);
  };

  return (
    <div ref={containerRef} className="article-list-container" style={{ padding: '16px', border: '1px solid var(--border-color)', marginTop: '16px', borderRadius: '8px', backgroundColor: 'var(--card-bg)' }}>
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
                  <th style={{ width: '180px', textAlign: 'left' }}>{t('trends.author_caption', 'Author')}</th>
                  <th style={{ width: 'auto', textAlign: 'left' }}>
                    <div className="article-header-info">
                      <span>{t('trends.article_caption', 'Article')}</span>
                      <div className="article-header-info-trigger-wrap">
                        <button
                          type="button"
                          className="article-header-info-trigger"
                          aria-label={t('trends.article_column_tooltip_content')}
                          aria-describedby={articleTooltipId}
                          onMouseEnter={(e) => showArticleTooltip(e.currentTarget)}
                          onMouseLeave={hideArticleTooltip}
                          onFocus={(e) => showArticleTooltip(e.currentTarget)}
                          onBlur={hideArticleTooltip}
                        >
                          i
                        </button>
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
                    <td style={{ padding: '8px 5px', verticalAlign: 'middle', fontSize: '0.9em', color: 'var(--secondary-text)' }}>
                      {formatAuthors(article.authors)}
                    </td>
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

          <div className="pagination-controls" style={{ display: 'flex', alignItems: 'center', marginTop: '15px' }}>
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
          {isArticleTooltipOpen && (
            <div
              id={articleTooltipId}
              role="tooltip"
              className="article-header-floating-tooltip"
              style={{ left: articleTooltipPos.x, top: articleTooltipPos.y }}
            >
              {t('trends.article_column_tooltip_content')}
            </div>
          )}
        </Fragment>
      )}
    </div>
  );
};
