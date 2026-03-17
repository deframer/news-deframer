import React, { useEffect, useMemo, useState } from 'react';
import { ChevronDown } from 'lucide-react-native';
import { Linking, Pressable, StyleSheet, Text, View } from 'react-native';
import { useTranslation } from 'react-i18next';

import { formatShortDate } from '../../../shared/formatTime';
import { getRatingColors, stripHtml, toPercent } from '../lib/articleFormat';
import { LoadingSpinner } from './LoadingSpinner';
import { AnalyzedArticle, AnalyzedItem, NewsDeframerClient } from '../services/newsDeframerClient';
import { logger } from '../services/logger';
import { Settings } from '../services/settingsService';
import { AppPalette } from '../theme';

const ARTICLES_PER_BATCH = 10;

export const TrendArticleListPanel = ({
  palette,
  term,
  domain,
  settings,
  daysInPast,
  selectedDate,
  headerTitle,
  onOpenArticle,
}: {
  palette: AppPalette;
  term: string;
  domain: string;
  settings: Settings;
  daysInPast?: number;
  selectedDate?: string;
  headerTitle?: string | null;
  onOpenArticle: (item: AnalyzedItem) => void;
}) => {
  const { t, i18n } = useTranslation();
  const [articles, setArticles] = useState<AnalyzedArticle[]>([]);
  const [loading, setLoading] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [offset, setOffset] = useState(0);
  const [hasMore, setHasMore] = useState(true);

  const client = useMemo(() => new NewsDeframerClient(settings), [settings]);

  const isFixedDateMode = Boolean(selectedDate);

  const resolvedHeaderTitle = useMemo(() => {
    if (headerTitle !== undefined) {
      return headerTitle;
    }

    if (selectedDate) {
      return `${t('trends.article_caption', 'Article')} / ${formatShortDate(selectedDate, i18n.language)}`;
    }

    return t('trends.article_caption', 'Article');
  }, [headerTitle, i18n.language, selectedDate, t]);

  useEffect(() => {
    setArticles([]);
    setOffset(0);
    setHasMore(true);
    setError(null);
  }, [daysInPast, domain, selectedDate, term]);

  useEffect(() => {
    let cancelled = false;

    const fetchArticles = async () => {
      if (!term || !domain) {
        setArticles([]);
        setHasMore(false);
        setLoading(false);
        setLoadingMore(false);
        return;
      }

      if (!selectedDate && daysInPast === undefined) {
        setArticles([]);
        setHasMore(false);
        return;
      }

      const firstPage = offset === 0;
      if (firstPage) {
        setLoading(true);
      } else {
        setLoadingMore(true);
      }
      setError(null);

      try {
        const result = await client.getArticlesByTrend(domain, term, selectedDate, selectedDate ? undefined : daysInPast, offset, ARTICLES_PER_BATCH);

        if (cancelled) {
          return;
        }

        setArticles((current) => (firstPage ? result : [...current, ...result]));
        setHasMore(result.length === ARTICLES_PER_BATCH);
      } catch (fetchError) {
        if (cancelled) {
          return;
        }

        logger.error('TrendArticleList fetch failed', { term, domain, selectedDate, daysInPast, offset, error: String(fetchError) });
        setError(fetchError instanceof Error ? fetchError.message : String(fetchError));
        if (firstPage) {
          setArticles([]);
        }
        setHasMore(false);
      } finally {
        if (!cancelled) {
          setLoading(false);
          setLoadingMore(false);
        }
      }
    };

    fetchArticles();

    return () => {
      cancelled = true;
    };
  }, [client, daysInPast, domain, offset, selectedDate, term]);

  const openArticleInBrowser = (url: string) => {
    Linking.openURL(url).catch(() => undefined);
  };

  const handleOpenArticle = async (article: AnalyzedArticle) => {
    if (article.rating === undefined) {
      openArticleInBrowser(article.url);
      return;
    }

    try {
      const item = await client.getItem(article.url);
      if (item) {
        onOpenArticle(item);
        return;
      }
    } catch (err) {
      logger.error('TrendArticleList open article fetch failed', { url: article.url, error: String(err) });
    }

    openArticleInBrowser(article.url);
  };

  return (
    <View style={styles.container}>
      {resolvedHeaderTitle ? <Text style={[styles.headerTitle, { color: palette.text }]}>{resolvedHeaderTitle}</Text> : null}

      {loading ? <LoadingSpinner palette={palette} center label={t('options.loading')} /> : null}

      {!loading && error ? (
        <View style={[styles.stateBox, { borderColor: palette.border, backgroundColor: palette.card }]}>
          <Text style={[styles.stateText, { color: palette.secondaryText }]}>{error}</Text>
        </View>
      ) : null}

      {!loading && !error && articles.length === 0 ? (
        <View style={[styles.stateBox, { borderColor: palette.border, backgroundColor: palette.card }]}>
          <Text style={[styles.stateText, { color: palette.secondaryText }]}>{t('trends.articles_no_data', { term, defaultValue: `No articles found for "${term}".` })}</Text>
        </View>
      ) : null}

      {!loading && !error && articles.length > 0 ? (
        <View style={styles.stack}>
          {articles.map((article, index) => {
            const rating = toPercent(article.rating);
            const ratingColors = getRatingColors(rating, palette);
            const authors = article.authors?.join(', ');
            const articleTitle = stripHtml(article.title) || article.url;

            return (
              <Pressable
                key={`${article.url}-${index}`}
                onPress={() => {
                  handleOpenArticle(article);
                }}
                style={[styles.row, { backgroundColor: palette.card, borderColor: palette.border }]}
              >
                <View style={styles.bodyColumn}>
                  <View style={styles.metaHeaderRow}>
                    <View style={[styles.ratingTrack, styles.ratingTrackCompact, { backgroundColor: palette.ratingBackground }]}>
                      <View style={[styles.ratingFill, { width: `${rating}%`, backgroundColor: ratingColors.backgroundColor }]} />
                      <Text style={[styles.ratingText, { color: ratingColors.textColor }]}>{rating}%</Text>
                    </View>
                    {!isFixedDateMode ? (
                      article.pub_date ? <Text style={[styles.metaText, styles.dateText, { color: palette.secondaryText }]}>{formatShortDate(article.pub_date, i18n.language)}</Text> : <View style={styles.dateSpacer} />
                    ) : null}
                    {authors ? <Text style={[styles.metaText, styles.authorTextInline, { color: palette.secondaryText }]} numberOfLines={1}>{authors}</Text> : null}
                  </View>
                  <Text style={[styles.articleText, { color: palette.accent }]}>{articleTitle}</Text>
                </View>
              </Pressable>
            );
          })}

          {hasMore ? (
            <Pressable
              onPress={() => setOffset((current) => current + ARTICLES_PER_BATCH)}
              disabled={loadingMore}
              style={[styles.loadMoreButton, { borderColor: palette.buttonBorder, backgroundColor: palette.buttonBackground }]}
            >
              {loadingMore ? <LoadingSpinner palette={palette} center label={t('options.loading')} /> : (
                <>
                  <ChevronDown color={palette.buttonText} size={16} strokeWidth={2.2} />
                  <Text style={[styles.loadMoreText, { color: palette.buttonText }]}>{t('mobile.load_more', 'Load more')}</Text>
                </>
              )}
            </Pressable>
          ) : (
            <View style={[styles.infoBox, { borderColor: palette.border, backgroundColor: palette.card }]}> 
              <Text style={[styles.infoText, { color: palette.secondaryText }]}>{t('mobile.article_list_info', 'Titles here are simplified by News Deframer. Open the article to see the original page.')}</Text>
            </View>
          )}
        </View>
      ) : null}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginTop: 14,
    gap: 12,
  },
  headerTitle: {
    fontSize: 18,
    fontWeight: '700',
    lineHeight: 24,
  },
  stateBox: {
    borderWidth: 1,
    borderRadius: 8,
    padding: 14,
  },
  stateText: {
    fontSize: 15,
    lineHeight: 22,
    textAlign: 'center',
    fontStyle: 'italic',
  },
  stack: {
    gap: 12,
  },
  row: {
    borderWidth: 1,
    borderRadius: 12,
    overflow: 'hidden',
    paddingHorizontal: 14,
    paddingVertical: 12,
  },
  ratingTrack: {
    height: 22,
    borderRadius: 999,
    overflow: 'hidden',
    justifyContent: 'center',
  },
  ratingTrackCompact: {
    width: 52,
    flexShrink: 0,
  },
  ratingFill: {
    position: 'absolute',
    left: 0,
    top: 0,
    bottom: 0,
  },
  ratingText: {
    paddingHorizontal: 8,
    fontSize: 11,
    fontWeight: '700',
  },
  bodyColumn: {
    gap: 6,
  },
  metaHeaderRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    minWidth: 0,
  },
  dateSpacer: {
    width: 52,
    flexShrink: 0,
  },
  dateText: {
    width: 52,
    flexShrink: 0,
  },
  metaText: {
    fontSize: 13,
    lineHeight: 18,
    fontWeight: '500',
  },
  authorTextInline: {
    flex: 1,
  },
  articleText: {
    fontSize: 16,
    lineHeight: 22,
    fontWeight: '600',
    textDecorationLine: 'underline',
    width: '100%',
  },
  loadMoreButton: {
    minHeight: 44,
    borderWidth: 1,
    borderRadius: 8,
    alignItems: 'center',
    justifyContent: 'center',
    flexDirection: 'row',
    gap: 8,
  },
  loadMoreText: {
    fontSize: 15,
    fontWeight: '700',
  },
  infoBox: {
    borderWidth: 1,
    borderRadius: 8,
    paddingHorizontal: 12,
    paddingVertical: 10,
  },
  infoText: {
    fontSize: 13,
    lineHeight: 18,
  },
});
