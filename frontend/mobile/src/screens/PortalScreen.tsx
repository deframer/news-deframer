import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Animated, FlatList, Modal, NativeScrollEvent, NativeSyntheticEvent, Platform, Pressable, ScrollView, StyleSheet, Text, View } from 'react-native';
import { Info } from 'lucide-react-native';
import { useTranslation } from 'react-i18next';

import { ArticleTile } from '../components/ArticleTile';
import { Card } from '../components/Card';
import { LoadingSpinner } from '../components/LoadingSpinner';
import { SegmentedControl } from '../components/SegmentedControl';
import { TrendComparePanel } from '../components/TrendComparePanel';
import { TrendSearchPanel } from '../components/TrendSearchPanel';
import { TrendTagCloudPanel } from '../components/TrendTagCloudPanel';
import { DomainEntry, NewsDeframerClient, AnalyzedItem } from '../services/newsDeframerClient';
import { Settings } from '../services/settingsService';
import { AppPalette } from '../theme';

type PortalTab = 'articles' | 'trend-mining';
type TrendSubview = 'cloud' | 'compare' | 'search';
type TrendRange = '24h' | '7d' | '30d' | '90d' | '365d';
type DomainOption = { id: string; name: string };

const TIME_RANGES: Array<{ id: TrendRange; label: string; days: number }> = [
  { id: '24h', label: 'trends.time_ranges.last_24h', days: 1 },
  { id: '7d', label: 'trends.time_ranges.last_7d', days: 7 },
  { id: '30d', label: 'trends.time_ranges.last_30d', days: 30 },
  { id: '90d', label: 'trends.time_ranges.last_90d', days: 90 },
  { id: '365d', label: 'trends.time_ranges.last_365d', days: 365 },
];

const ArticleSeparator = () => <View style={styles.articleSeparator} />;

export const PortalScreen = ({
  palette,
  domain,
  settings,
  onOpenArticle,
}: {
  palette: AppPalette;
  domain: DomainEntry;
  settings: Settings;
  onOpenArticle: (item: AnalyzedItem) => void;
}) => {
  const { t } = useTranslation();
  const [activeTab, setActiveTab] = useState<PortalTab>('articles');
  const [items, setItems] = useState<AnalyzedItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [reasonText, setReasonText] = useState<string | null>(null);
  const [trendSubview, setTrendSubview] = useState<TrendSubview>('cloud');
  const [trendRange, setTrendRange] = useState<TrendRange>('7d');
  const [availableDomains, setAvailableDomains] = useState<DomainEntry[]>([]);
  const portalScrollRef = useRef<ScrollView>(null);
  const scrollOffsetRef = useRef(0);
  const fadeAnim = useRef(new Animated.Value(0)).current;
  const scaleAnim = useRef(new Animated.Value(0.96)).current;
  const useNativeDriver = Platform.OS !== 'web';

  const client = useMemo(() => new NewsDeframerClient(settings), [settings]);
  const domainOptions = useMemo<DomainOption[]>(() => availableDomains.filter((d) => d.domain !== domain.domain && d.language === domain.language).map((d) => ({ id: d.domain, name: d.domain })), [availableDomains, domain.domain, domain.language]);
  const [compareDomain, setCompareDomain] = useState<string | null>(null);

  const loadItems = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      const result = await client.getSite(domain.domain);
      setItems(result);
    } catch {
      setItems([]);
      setError(t('mobile.portal_load_error'));
    } finally {
      setIsLoading(false);
    }
  }, [client, domain.domain, t]);

  useEffect(() => {
    loadItems();
  }, [loadItems]);

  useEffect(() => {
    let cancelled = false;

    const loadDomains = async () => {
      try {
        const domains = await client.getDomains();
        if (!cancelled) {
          setAvailableDomains(domains);
        }
      } catch {
        if (!cancelled) {
          setAvailableDomains([]);
        }
      }
    };

    loadDomains();

    return () => {
      cancelled = true;
    };
  }, [client]);

  useEffect(() => {
    if (domainOptions.length === 0) {
      setCompareDomain(null);
      return;
    }

    if (!compareDomain || !domainOptions.some((option) => option.id === compareDomain)) {
      setCompareDomain(domainOptions[0].id);
    }
  }, [compareDomain, domainOptions]);

  useEffect(() => {
    if (!reasonText) {
      fadeAnim.setValue(0);
      scaleAnim.setValue(0.96);
      return;
    }

    Animated.parallel([
      Animated.timing(fadeAnim, {
        toValue: 1,
        duration: 180,
        useNativeDriver,
      }),
      Animated.timing(scaleAnim, {
        toValue: 1,
        duration: 180,
        useNativeDriver,
      }),
    ]).start();
  }, [fadeAnim, reasonText, scaleAnim, useNativeDriver]);

  useEffect(() => {
    if (typeof document === 'undefined') {
      return;
    }

    const html = document.documentElement;
    const body = document.body;
    const root = document.getElementById('root');
    const previousHtmlOverflow = html.style.overflow;
    const previousBodyOverflow = body.style.overflow;
    const previousRootOverflow = root?.style.overflow ?? '';

    if (reasonText) {
      html.style.overflow = 'hidden';
      body.style.overflow = 'hidden';
      if (root) {
        root.style.overflow = 'hidden';
      }
    }

    return () => {
      html.style.overflow = previousHtmlOverflow;
      body.style.overflow = previousBodyOverflow;
      if (root) {
        root.style.overflow = previousRootOverflow;
      }
    };
  }, [reasonText]);

  const handlePortalScroll = useCallback((event: NativeSyntheticEvent<NativeScrollEvent>) => {
    scrollOffsetRef.current = event.nativeEvent.contentOffset.y;
  }, []);

  const restorePortalScroll = useCallback((offset: number) => {
    requestAnimationFrame(() => {
      portalScrollRef.current?.scrollTo({ y: offset, animated: false });
    });
  }, []);

  const renderArticleItem = useCallback(
    ({ item }: { item: AnalyzedItem }) => <ArticleTile item={item} palette={palette} onOpenArticle={onOpenArticle} onShowReason={setReasonText} />,
    [onOpenArticle, palette],
  );

  const renderArticleListEmpty = useCallback(() => {
    if (isLoading) {
      return <LoadingSpinner palette={palette} label={t('options.loading')} />;
    }

    if (error) {
      return (
        <Card palette={palette}>
          <Text style={[styles.stateTitle, { color: palette.text }]}>{t('mobile.portal_load_error_title')}</Text>
          <Text style={[styles.stateBody, { color: palette.secondaryText }]}>{error}</Text>
          <Pressable onPress={loadItems} style={[styles.actionButton, { backgroundColor: palette.buttonBackground, borderColor: palette.buttonBorder }]}> 
            <Text style={[styles.actionButtonText, { color: palette.buttonText }]}>{t('mobile.retry')}</Text>
          </Pressable>
        </Card>
      );
    }

    return (
      <Card palette={palette}>
        <Text style={[styles.stateTitle, { color: palette.text }]}>{t('mobile.portal_empty_title')}</Text>
        <Text style={[styles.stateBody, { color: palette.secondaryText }]}>{t('mobile.portal_empty_body')}</Text>
      </Card>
    );
  }, [error, isLoading, loadItems, palette, t]);

  return (
    <View style={[styles.container, { backgroundColor: palette.background }]}> 
      <Modal animationType="none" transparent visible={Boolean(reasonText)} onRequestClose={() => setReasonText(null)}>
        <Animated.View style={[styles.modalOverlay, { opacity: fadeAnim }]}> 
          <Pressable style={StyleSheet.absoluteFill} onPress={() => setReasonText(null)} />
          <Animated.View
            accessibilityViewIsModal
            aria-modal={true}
            style={[styles.modalCard, { backgroundColor: palette.card, borderColor: palette.border, transform: [{ scale: scaleAnim }] }]}
          >
            <View style={styles.modalHeader}>
              <Info color={palette.text} size={18} strokeWidth={2.2} />
              <Text style={[styles.modalTitle, { color: palette.text }]}>{t('rating.reason_title', 'Rating reason')}</Text>
            </View>
            <Text style={[styles.modalBody, { color: palette.secondaryText }]}>{reasonText}</Text>
            <Pressable onPress={() => setReasonText(null)} style={[styles.modalButton, { backgroundColor: palette.accent }]}> 
              <Text style={[styles.modalButtonText, { color: palette.accentText }]}>{t('mobile.ok')}</Text>
            </Pressable>
          </Animated.View>
        </Animated.View>
      </Modal>

      {activeTab === 'articles' ? (
        <FlatList
          data={!isLoading && !error ? items : []}
          keyExtractor={(item) => item.url}
          renderItem={renderArticleItem}
          ListHeaderComponent={
            <View style={[styles.stickyTabs, { backgroundColor: palette.background }]}> 
              <SegmentedControl
                palette={palette}
                value={activeTab}
                onChange={(value) => setActiveTab(value as PortalTab)}
                options={[{ label: t('trends.articles'), value: 'articles' }, { label: t('mobile.trends'), value: 'trend-mining' }]}
              />
            </View>
          }
          ListEmptyComponent={renderArticleListEmpty}
          contentContainerStyle={styles.articleListContent}
          style={styles.scrollView}
          scrollEnabled={!reasonText}
          ItemSeparatorComponent={ArticleSeparator}
          removeClippedSubviews={Platform.OS === 'android'}
          initialNumToRender={6}
          maxToRenderPerBatch={6}
          windowSize={5}
        />
      ) : (
        <ScrollView
          ref={portalScrollRef}
          style={styles.scrollView}
          contentContainerStyle={[styles.content, styles.contentWithBottomTabs]}
          scrollEnabled={!reasonText}
          onScroll={handlePortalScroll}
          scrollEventThrottle={16}
        >
          <View style={[styles.stickyTabs, { backgroundColor: palette.background }]}> 
            <SegmentedControl
              palette={palette}
              value={activeTab}
              onChange={(value) => setActiveTab(value as PortalTab)}
              options={[{ label: t('trends.articles'), value: 'articles' }, { label: t('mobile.trends'), value: 'trend-mining' }]}
            />
          </View>

          <View style={styles.stack}>
            <View style={styles.timeRow}>
              {TIME_RANGES.map((range) => {
                const selected = trendRange === range.id;
                return (
                  <Pressable
                    key={range.id}
                    onPress={() => setTrendRange(range.id)}
                    style={[
                      styles.timeChip,
                      {
                        borderColor: palette.buttonBorder,
                        backgroundColor: selected ? palette.accent : palette.buttonBackground,
                      },
                    ]}
                  >
                    <Text style={[styles.timeChipText, { color: selected ? palette.accentText : palette.buttonText }]}>{t(range.label)}</Text>
                  </Pressable>
                );
              })}
            </View>

            {trendSubview === 'cloud' ? (
              <TrendTagCloudPanel
                palette={palette}
                domain={domain.domain}
                language={domain.language}
                daysInPast={TIME_RANGES.find((range) => range.id === trendRange)?.days || 7}
                settings={settings}
                onOpenArticle={onOpenArticle}
              />
            ) : null}
            {trendSubview === 'compare' ? (
              <TrendComparePanel
                palette={palette}
                domain={domain.domain}
                language={domain.language}
                daysInPast={TIME_RANGES.find((range) => range.id === trendRange)?.days || 7}
                settings={settings}
                availableDomains={domainOptions}
                compareDomain={compareDomain}
                onSelectDomain={setCompareDomain}
                getScrollOffset={() => scrollOffsetRef.current}
                onRestoreScrollOffset={restorePortalScroll}
                onOpenArticle={onOpenArticle}
              />
            ) : null}
            {trendSubview === 'search' ? (
              <TrendSearchPanel
                palette={palette}
                domain={domain.domain}
                language={domain.language}
                daysInPast={TIME_RANGES.find((range) => range.id === trendRange)?.days || 7}
                settings={settings}
                onOpenArticle={onOpenArticle}
              />
            ) : null}
          </View>
        </ScrollView>
      )}

      {activeTab === 'trend-mining' ? (
        <View style={[styles.trendSubviewDock, { backgroundColor: palette.background, borderTopColor: palette.border }]}> 
          <SegmentedControl
            palette={palette}
            value={trendSubview}
            onChange={(value) => setTrendSubview(value as TrendSubview)}
            options={[
              { label: t('trends.cloud'), value: 'cloud' },
              { label: t('trends.compare_view'), value: 'compare' },
              { label: t('trends.search'), value: 'search' },
            ]}
          />
        </View>
      ) : null}
    </View>
  );
};

const styles = StyleSheet.create({
  container: { flex: 1, overflow: 'hidden' },
  scrollView: { flex: 1, overflow: 'scroll' },
  content: { paddingBottom: 24 },
  contentWithBottomTabs: { paddingBottom: 100 },
  articleListContent: { paddingBottom: 24, paddingHorizontal: 24 },
  articleSeparator: { height: 16 },
  stickyTabs: {
    paddingBottom: 16,
    paddingTop: 16,
    paddingHorizontal: 24,
    borderBottomWidth: 1,
  },
  spacer: { height: 16 },
  stack: { paddingHorizontal: 24, gap: 16 },
  timeRow: { flexDirection: 'row', gap: 8, flexWrap: 'wrap', justifyContent: 'flex-end' },
  timeChip: {
    borderWidth: 1,
    borderRadius: 8,
    minHeight: 34,
    paddingHorizontal: 10,
    alignItems: 'center',
    justifyContent: 'center',
  },
  timeChipText: { fontSize: 14, fontWeight: '700' },
  trendSubviewDock: {
    position: 'absolute',
    left: 0,
    right: 0,
    bottom: 0,
    borderTopWidth: 1,
    paddingHorizontal: 24,
    paddingTop: 12,
    paddingBottom: 12,
  },
  stateTitle: { marginBottom: 8, fontSize: 20, fontWeight: '700' },
  stateBody: { fontSize: 15, lineHeight: 22 },
  actionButton: {
    marginTop: 16,
    borderWidth: 1,
    borderRadius: 8,
    minHeight: 46,
    paddingHorizontal: 16,
    alignItems: 'center',
    justifyContent: 'center',
    alignSelf: 'flex-start',
  },
  actionButtonText: { fontSize: 16, fontWeight: '600' },
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0,0,0,0.7)',
    alignItems: 'center',
    justifyContent: 'center',
    padding: 24,
  },
  modalCard: {
    width: '100%',
    maxWidth: 420,
    borderWidth: 1,
    borderRadius: 16,
    padding: 20,
    shadowColor: '#000000',
    shadowOpacity: 0.35,
    shadowRadius: 24,
    shadowOffset: { width: 0, height: 10 },
    elevation: 10,
  },
  modalHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 10,
    marginBottom: 12,
  },
  modalTitle: {
    fontSize: 20,
    fontWeight: '700',
  },
  modalBody: {
    fontSize: 15,
    lineHeight: 22,
  },
  modalButton: {
    marginTop: 18,
    alignSelf: 'flex-end',
    borderRadius: 10,
    paddingHorizontal: 18,
    paddingVertical: 10,
  },
  modalButtonText: {
    fontSize: 15,
    fontWeight: '700',
  },
});
