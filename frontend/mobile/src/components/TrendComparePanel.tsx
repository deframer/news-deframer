import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { ExternalLink } from 'lucide-react-native';
import { Pressable, StyleSheet, Text, View } from 'react-native';
import { Dropdown } from 'react-native-element-dropdown';
import { useTranslation } from 'react-i18next';

import { colorValues } from '../../../shared/theme';
import { Card } from './Card';
import { LoadingSpinner } from './LoadingSpinner';
import { TrendArticleListPanel } from './TrendArticleListPanel';
import { AnalyzedItem, DomainComparison, NewsDeframerClient } from '../services/newsDeframerClient';
import { logger } from '../services/logger';
import { Settings } from '../services/settingsService';
import { AppPalette } from '../theme';

type DomainOption = { id: string; name: string };
type SelectedCompareItem = { term: string; domain: string };

const BULLET_DELIMITER = '•';

const getTrendActions = (classification: DomainComparison['classification']) => {
  switch (classification) {
    case 'BLINDSPOT_A':
      return ['current'] as const;
    case 'BLINDSPOT_B':
      return ['compare'] as const;
    case 'INTERSECT':
      return ['current', 'compare'] as const;
  }
};

const renderSectionRows = ({
  items,
  palette,
  mode,
  currentDomain,
  compareDomain,
  onSelectItem,
}: {
  items: DomainComparison[];
  palette: AppPalette;
  mode: 'current' | 'compare' | 'shared';
  currentDomain: string;
  compareDomain: string;
  onSelectItem: (selection: SelectedCompareItem) => void;
}) =>
  items.map((item, index) => {
    const actions = getTrendActions(item.classification);
    const rowSelection = mode === 'current' ? { term: item.trend_topic, domain: currentDomain } : mode === 'compare' ? { term: item.trend_topic, domain: compareDomain } : null;

    return (
      <Pressable
        key={`${item.classification}-${item.trend_topic}-${item.rank_group}`}
        disabled={mode === 'shared'}
        onPress={rowSelection ? () => onSelectItem(rowSelection) : undefined}
        style={[
          styles.row,
          { borderBottomColor: palette.border },
          index === items.length - 1 ? styles.rowLast : null,
        ]}
      >
        <Text style={[styles.topic, { color: palette.text }]} numberOfLines={2}>
          {item.trend_topic}
        </Text>
        <View style={styles.actions}>
          {actions.map((action) => {
            const isCompare = action === 'compare';
            const compareColor = palette.accent;
            const buttonColor = isCompare ? colorValues.white : palette.background;
            const actionButtonStyle = {
              borderColor: isCompare ? compareColor : palette.text,
              backgroundColor: isCompare ? compareColor : palette.text,
            };
            const actionSelection = {
              term: item.trend_topic,
              domain: isCompare ? compareDomain : currentDomain,
            };

            return (
              <Pressable
                key={`${item.trend_topic}-${action}`}
                onPress={() => onSelectItem(actionSelection)}
                style={[styles.actionButton, actionButtonStyle]}
              >
                <ExternalLink color={buttonColor} size={14} strokeWidth={2.2} />
              </Pressable>
            );
          })}
        </View>
      </Pressable>
    );
  });

export const TrendComparePanel = ({
  palette,
  domain,
  language,
  daysInPast,
  settings,
  availableDomains,
  compareDomain,
  onSelectDomain,
  getScrollOffset,
  onRestoreScrollOffset,
  onOpenArticle,
  onBackRequestChange,
}: {
  palette: AppPalette;
  domain: string;
  language: string;
  daysInPast: number;
  settings: Settings;
  availableDomains: DomainOption[];
  compareDomain: string | null;
  onSelectDomain: (domain: string) => void;
  getScrollOffset: () => number;
  onRestoreScrollOffset: (offset: number) => void;
  onOpenArticle: (item: AnalyzedItem) => void;
  onBackRequestChange: (action: (() => void) | null) => void;
}) => {
  const { t } = useTranslation();
  const [items, setItems] = useState<DomainComparison[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [hasError, setHasError] = useState(false);
  const [selectedItem, setSelectedItem] = useState<SelectedCompareItem | null>(null);
  const [savedScrollOffset, setSavedScrollOffset] = useState<number | null>(null);

  const client = useMemo(() => new NewsDeframerClient(settings), [settings]);
  const data = useMemo(() => availableDomains.map((option) => ({ label: option.name, value: option.id })), [availableDomains]);
  const sortedItems = useMemo(
    () => [...items].sort((left, right) => left.rank_group - right.rank_group || left.trend_topic.localeCompare(right.trend_topic)),
    [items],
  );
  const currentOnlyItems = useMemo(() => sortedItems.filter((item) => item.classification === 'BLINDSPOT_A'), [sortedItems]);
  const compareOnlyItems = useMemo(() => sortedItems.filter((item) => item.classification === 'BLINDSPOT_B'), [sortedItems]);
  const sharedItems = useMemo(() => sortedItems.filter((item) => item.classification === 'INTERSECT'), [sortedItems]);

  const selectedCompareDomain = compareDomain ?? '';
  const compareColor = palette.accent;

  const isSelectedItemPresent = useMemo(() => {
    if (!selectedItem) {
      return false;
    }

    return items.some((item) => {
      if (item.trend_topic !== selectedItem.term) {
        return false;
      }

      if (selectedItem.domain === domain) {
        return item.classification === 'BLINDSPOT_A' || item.classification === 'INTERSECT';
      }

      if (selectedItem.domain === selectedCompareDomain) {
        return item.classification === 'BLINDSPOT_B' || item.classification === 'INTERSECT';
      }

      return false;
    });
  }, [domain, items, selectedCompareDomain, selectedItem]);

  const handleSelectItem = (selection: SelectedCompareItem) => {
    setSavedScrollOffset(getScrollOffset());
    setSelectedItem(selection);
  };

  const handleCloseSelection = useCallback(() => {
    setSelectedItem(null);
    if (savedScrollOffset !== null) {
      onRestoreScrollOffset(savedScrollOffset);
    }
  }, [onRestoreScrollOffset, savedScrollOffset]);

  useEffect(() => {
    let cancelled = false;

    const fetchCompare = async () => {
      if (!compareDomain) {
        setItems([]);
        setHasError(false);
        setIsLoading(false);
        return;
      }

      setIsLoading(true);
      setHasError(false);

      try {
        const response = await client.getDomainComparison(domain, compareDomain, language, daysInPast);
        if (!cancelled) {
          setItems(response);
          logger.info('TrendCompare fetch success', { domain, compareDomain, language, daysInPast, itemCount: response.length });
        }
      } catch {
        if (!cancelled) {
          setItems([]);
          setHasError(true);
          logger.error('TrendCompare fetch failed', { domain, compareDomain, language, daysInPast });
        }
      } finally {
        if (!cancelled) {
          setIsLoading(false);
        }
      }
    };

    fetchCompare();

    return () => {
      cancelled = true;
    };
  }, [client, compareDomain, daysInPast, domain, language]);

  useEffect(() => {
    if (!selectedItem) {
      return;
    }

    if (!isSelectedItemPresent) {
      setSelectedItem(null);
      setSavedScrollOffset(null);
    }
  }, [isSelectedItemPresent, selectedItem]);

  useEffect(() => {
    setSelectedItem(null);
    setSavedScrollOffset(null);
  }, [domain, selectedCompareDomain, daysInPast]);

  useEffect(() => {
    onBackRequestChange(selectedItem ? handleCloseSelection : null);

    return () => onBackRequestChange(null);
  }, [handleCloseSelection, onBackRequestChange, selectedItem]);

  return (
    <View style={styles.stack}>
      {!selectedItem ? (
        <View>
          <Text style={[styles.label, { color: palette.secondaryText }]}>{t('mobile.compare_with_label', 'Compare with')}</Text>
          <Dropdown
            mode="modal"
            data={data}
            search
            disable={data.length === 0}
            labelField="label"
            valueField="value"
            value={compareDomain}
            placeholder={data.length === 0 ? 'No domains available.' : 'Select domain'}
            searchPlaceholder="Search..."
            onChange={(item) => onSelectDomain(item.value)}
            style={[styles.dropdown, { borderColor: palette.buttonBorder, backgroundColor: palette.background }]}
            containerStyle={[styles.dropdownContainer, { backgroundColor: palette.card, borderColor: palette.border }]}
            placeholderStyle={[styles.placeholder, { color: palette.secondaryText }]}
            selectedTextStyle={[styles.selectedText, { color: palette.accent }]}
            itemTextStyle={[styles.itemText, { color: palette.text }]}
            inputSearchStyle={[styles.inputSearch, { color: palette.text, borderColor: palette.buttonBorder, backgroundColor: palette.background }]}
            activeColor={palette.accent}
            iconColor={palette.text}
            backgroundColor="rgba(0,0,0,0.35)"
            maxHeight={320}
          />
        </View>
      ) : null}

      {!compareDomain ? null : (
        <View style={styles.stack}>
          {!selectedItem ? <View style={[styles.divider, { borderBottomColor: palette.border }]} /> : null}

          {selectedItem ? (
            <Pressable
              onPress={handleCloseSelection}
              style={[styles.selectedTrendPill, { backgroundColor: palette.card, borderColor: palette.border }]}
              accessibilityRole="button"
              accessibilityLabel={`${t('mobile.change_term')}: ${selectedItem.term} ${BULLET_DELIMITER} ${selectedItem.domain}`}
            >
              <Text style={[styles.selectedTrendText, { color: colorValues.white }]}>
                {selectedItem.term} {BULLET_DELIMITER} {selectedItem.domain}
              </Text>
            </Pressable>
          ) : null}

          <Card palette={palette}>
          {isLoading ? <LoadingSpinner palette={palette} center label={t('options.loading')} /> : null}
          {!isLoading && hasError ? <Text style={[styles.stateText, { color: palette.secondaryText }]}>{t('mobile.portal_load_error')}</Text> : null}
          {!isLoading && !hasError && sortedItems.length === 0 ? <Text style={[styles.stateText, { color: palette.secondaryText }]}>{t('trends.no_data')}</Text> : null}

          {!selectedItem && !isLoading && !hasError && sortedItems.length > 0 ? (
            <View style={styles.list}>
              <View style={styles.section}>
                <View style={[styles.sectionHeader, { backgroundColor: palette.text }]}> 
                  <Text style={[styles.sectionTitle, { color: palette.background }]}>{domain}</Text>
                </View>
                {currentOnlyItems.length > 0 ? renderSectionRows({ items: currentOnlyItems, palette, mode: 'current', currentDomain: domain, compareDomain: selectedCompareDomain, onSelectItem: handleSelectItem }) : <Text style={[styles.emptySectionText, { color: palette.secondaryText }]}>{t('trends.no_data')}</Text>}
              </View>

              <View style={styles.section}>
                <View style={[styles.sectionHeader, styles.compareSectionHeader, { backgroundColor: compareColor }]}> 
                  <Text style={[styles.sectionTitle, styles.compareSectionTitle, { color: colorValues.white }]}>{compareDomain}</Text>
                </View>
                {compareOnlyItems.length > 0 ? renderSectionRows({ items: compareOnlyItems, palette, mode: 'compare', currentDomain: domain, compareDomain: selectedCompareDomain, onSelectItem: handleSelectItem }) : <Text style={[styles.emptySectionText, { color: palette.secondaryText }]}>{t('trends.no_data')}</Text>}
              </View>

              <View style={styles.section}>
                <View style={[styles.sectionHeader, { backgroundColor: palette.text }]}> 
                  <Text style={[styles.sectionTitle, { color: palette.background }]}>{t('trends.compare.shared', 'Shared Trends')}</Text>
                </View>
                {sharedItems.length > 0 ? renderSectionRows({ items: sharedItems, palette, mode: 'shared', currentDomain: domain, compareDomain: selectedCompareDomain, onSelectItem: handleSelectItem }) : <Text style={[styles.emptySectionText, { color: palette.secondaryText }]}>{t('trends.no_data')}</Text>}
              </View>
            </View>
          ) : null}

          {selectedItem && !isLoading && !hasError ? <TrendArticleListPanel palette={palette} term={selectedItem.term} domain={selectedItem.domain} settings={settings} daysInPast={daysInPast} headerTitle={t('trends.article_caption', 'Article')} onOpenArticle={onOpenArticle} /> : null}
          </Card>
        </View>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  stack: { gap: 16 },
  label: {
    marginBottom: 8,
    fontSize: 14,
    fontWeight: '700',
  },
  dropdown: {
    minHeight: 54,
    borderWidth: 1,
    borderRadius: 8,
    paddingHorizontal: 16,
  },
  dropdownContainer: {
    borderWidth: 1,
    borderRadius: 12,
    overflow: 'hidden',
  },
  placeholder: {
    fontSize: 18,
    fontWeight: '500',
  },
  selectedText: {
    fontSize: 18,
    fontWeight: '500',
  },
  itemText: {
    fontSize: 18,
    fontWeight: '500',
  },
  inputSearch: {
    minHeight: 46,
    borderWidth: 1,
    borderRadius: 8,
    margin: 12,
    paddingHorizontal: 14,
    fontSize: 16,
  },
  stateText: {
    minHeight: 80,
    textAlign: 'center',
    textAlignVertical: 'center',
    fontSize: 15,
    lineHeight: 22,
  },
  divider: {
    borderBottomWidth: 1,
  },
  list: {
    margin: -16,
    gap: 16,
  },
  selectedTrendPill: {
    borderWidth: 1,
    borderRadius: 8,
    minHeight: 44,
    paddingHorizontal: 16,
    paddingVertical: 10,
    justifyContent: 'center',
    alignItems: 'center',
  },
  selectedTrendText: {
    fontSize: 20,
    lineHeight: 26,
    fontWeight: '700',
    textAlign: 'center',
  },
  section: {
    gap: 0,
  },
  sectionHeader: {
    minHeight: 60,
    paddingHorizontal: 16,
    justifyContent: 'center',
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '700',
  },
  compareSectionTitle: {
    letterSpacing: 0.2,
  },
  compareSectionHeader: {
  },
  emptySectionText: {
    paddingHorizontal: 16,
    paddingVertical: 18,
    fontSize: 14,
    fontStyle: 'italic',
  },
  row: {
    minHeight: 60,
    paddingHorizontal: 16,
    paddingVertical: 12,
    borderTopWidth: 1,
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
  },
  rowLast: {
    borderBottomWidth: 0,
  },
  topic: {
    flex: 1,
    fontSize: 18,
    lineHeight: 24,
    fontWeight: '600',
  },
  actions: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  actionButton: {
    width: 44,
    height: 44,
    borderWidth: 1,
    borderRadius: 999,
    alignItems: 'center',
    justifyContent: 'center',
  },
});
