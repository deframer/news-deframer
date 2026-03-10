import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Animated, Modal, Platform, Pressable, ScrollView, StyleSheet, Text, View } from 'react-native';
import { Info } from 'lucide-react-native';
import { useTranslation } from 'react-i18next';

import { ArticleTile } from '../components/ArticleTile';
import { Card } from '../components/Card';
import { LoadingSpinner } from '../components/LoadingSpinner';
import { SegmentedControl } from '../components/SegmentedControl';
import { DomainEntry, NewsDeframerClient, AnalyzedItem } from '../services/newsDeframerClient';
import { Settings } from '../services/settingsService';
import { AppPalette } from '../theme';

type PortalTab = 'articles' | 'trend-mining';

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
  const fadeAnim = useRef(new Animated.Value(0)).current;
  const scaleAnim = useRef(new Animated.Value(0.96)).current;
  const useNativeDriver = Platform.OS !== 'web';

  const client = useMemo(() => new NewsDeframerClient(settings), [settings]);

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

      <ScrollView
        style={styles.scrollView}
        contentContainerStyle={styles.content}
        scrollEnabled={!reasonText}
        stickyHeaderIndices={[0]}
      >
        <View style={[styles.stickyTabs, { backgroundColor: palette.background }]}> 
          <SegmentedControl
            palette={palette}
            value={activeTab}
            onChange={(value) => setActiveTab(value as PortalTab)}
            options={[{ label: t('trends.articles'), value: 'articles' }, { label: t('mobile.trend_mining'), value: 'trend-mining' }]}
          />
        </View>

        {activeTab === 'articles' ? (
          <View style={styles.stack}>
            {isLoading ? <LoadingSpinner palette={palette} label={t('options.loading')} /> : null}

            {!isLoading && error ? (
              <Card palette={palette}>
                <Text style={[styles.stateTitle, { color: palette.text }]}>{t('mobile.portal_load_error_title')}</Text>
                <Text style={[styles.stateBody, { color: palette.secondaryText }]}>{error}</Text>
                <Pressable onPress={loadItems} style={[styles.actionButton, { backgroundColor: palette.buttonBackground, borderColor: palette.buttonBorder }]}> 
                  <Text style={[styles.actionButtonText, { color: palette.buttonText }]}>{t('mobile.retry')}</Text>
                </Pressable>
              </Card>
            ) : null}

            {!isLoading && !error && items.length === 0 ? (
              <Card palette={palette}>
                <Text style={[styles.stateTitle, { color: palette.text }]}>{t('mobile.portal_empty_title')}</Text>
                <Text style={[styles.stateBody, { color: palette.secondaryText }]}>{t('mobile.portal_empty_body')}</Text>
              </Card>
            ) : null}

            {!isLoading && !error ? items.map((item) => <ArticleTile key={item.url} item={item} palette={palette} onOpenArticle={onOpenArticle} onShowReason={setReasonText} />) : null}
          </View>
        ) : (
          <View style={styles.stack}>
            <Card palette={palette}>
              <Text style={[styles.stateTitle, { color: palette.text }]}>{t('mobile.trend_mining')}</Text>
              <Text style={[styles.stateBody, { color: palette.secondaryText }]}>{t('mobile.trend_mining_placeholder')}</Text>
            </Card>
          </View>
        )}
      </ScrollView>
    </View>
  );
};

const styles = StyleSheet.create({
  container: { flex: 1, overflow: 'hidden' },
  scrollView: { flex: 1, overflow: 'scroll' },
  content: { paddingBottom: 24 },
  stickyTabs: {
    paddingHorizontal: 24,
    paddingBottom: 16,
    paddingTop: 16,
    borderBottomWidth: 1,
  },
  spacer: { height: 16 },
  stack: { paddingHorizontal: 24, gap: 16 },
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
