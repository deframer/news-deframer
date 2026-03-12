import React, { useEffect, useMemo, useState } from 'react';
import { ExternalLink } from 'lucide-react-native';
import { Pressable, StyleSheet, Text, View } from 'react-native';
import { Dropdown } from 'react-native-element-dropdown';
import { useTranslation } from 'react-i18next';

import { colorValues } from '../../../shared/theme';
import { Card } from './Card';
import { LoadingSpinner } from './LoadingSpinner';
import { DomainComparison, NewsDeframerClient } from '../services/newsDeframerClient';
import { logger } from '../services/logger';
import { Settings } from '../services/settingsService';
import { AppPalette } from '../theme';

type DomainOption = { id: string; name: string };

const COMPARE_DOMAIN_COLOR = colorValues.compareDomain;

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
}: {
  items: DomainComparison[];
  palette: AppPalette;
}) =>
  items.map((item, index) => {
    const actions = getTrendActions(item.classification);

    return (
      <View
        key={`${item.classification}-${item.trend_topic}-${item.rank_group}`}
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
            const buttonColor = isCompare ? colorValues.white : palette.background;
            const actionButtonStyle = {
              borderColor: isCompare ? COMPARE_DOMAIN_COLOR : palette.text,
              backgroundColor: isCompare ? COMPARE_DOMAIN_COLOR : palette.text,
            };

            return (
              <Pressable key={`${item.trend_topic}-${action}`} disabled style={[styles.actionButton, actionButtonStyle]}>
                <ExternalLink color={buttonColor} size={14} strokeWidth={2.2} />
              </Pressable>
            );
          })}
        </View>
      </View>
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
}: {
  palette: AppPalette;
  domain: string;
  language: string;
  daysInPast: number;
  settings: Settings;
  availableDomains: DomainOption[];
  compareDomain: string | null;
  onSelectDomain: (domain: string) => void;
}) => {
  const { t } = useTranslation();
  const [items, setItems] = useState<DomainComparison[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [hasError, setHasError] = useState(false);

  const client = useMemo(() => new NewsDeframerClient(settings), [settings]);
  const data = useMemo(() => availableDomains.map((option) => ({ label: option.name, value: option.id })), [availableDomains]);
  const sortedItems = useMemo(
    () => [...items].sort((left, right) => left.rank_group - right.rank_group || left.trend_topic.localeCompare(right.trend_topic)),
    [items],
  );
  const currentOnlyItems = useMemo(() => sortedItems.filter((item) => item.classification === 'BLINDSPOT_A'), [sortedItems]);
  const compareOnlyItems = useMemo(() => sortedItems.filter((item) => item.classification === 'BLINDSPOT_B'), [sortedItems]);
  const sharedItems = useMemo(() => sortedItems.filter((item) => item.classification === 'INTERSECT'), [sortedItems]);

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

  return (
    <View style={styles.stack}>
      <View>
        <Text style={[styles.label, { color: palette.secondaryText }]}>Compare to</Text>
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
          selectedTextStyle={[styles.selectedText, { color: palette.text }]}
          itemTextStyle={[styles.itemText, { color: palette.text }]}
          inputSearchStyle={[styles.inputSearch, { color: palette.text, borderColor: palette.buttonBorder, backgroundColor: palette.background }]}
          activeColor={palette.accent}
          iconColor={palette.text}
          backgroundColor="rgba(0,0,0,0.35)"
          maxHeight={320}
        />
      </View>

      {!compareDomain ? null : (
        <Card palette={palette}>
          {isLoading ? <LoadingSpinner palette={palette} center label={t('options.loading')} /> : null}
          {!isLoading && hasError ? <Text style={[styles.stateText, { color: palette.secondaryText }]}>{t('mobile.portal_load_error')}</Text> : null}
          {!isLoading && !hasError && sortedItems.length === 0 ? <Text style={[styles.stateText, { color: palette.secondaryText }]}>{t('trends.no_data')}</Text> : null}

          {!isLoading && !hasError && sortedItems.length > 0 ? (
            <View style={styles.list}>
              <View style={styles.section}>
                <View style={[styles.sectionHeader, { backgroundColor: palette.text }]}> 
                  <Text style={[styles.sectionTitle, { color: palette.background }]}>{domain}</Text>
                </View>
                {currentOnlyItems.length > 0 ? renderSectionRows({ items: currentOnlyItems, palette }) : <Text style={[styles.emptySectionText, { color: palette.secondaryText }]}>{t('trends.no_data')}</Text>}
              </View>

              <View style={styles.section}>
                <View style={[styles.sectionHeader, styles.compareSectionHeader, { backgroundColor: COMPARE_DOMAIN_COLOR }]}>
                  <Text style={[styles.sectionTitle, styles.compareSectionTitle, { color: colorValues.white }]}>{compareDomain}</Text>
                </View>
                {compareOnlyItems.length > 0 ? renderSectionRows({ items: compareOnlyItems, palette }) : <Text style={[styles.emptySectionText, { color: palette.secondaryText }]}>{t('trends.no_data')}</Text>}
              </View>

              <View style={styles.section}>
                <View style={[styles.sectionHeader, { backgroundColor: palette.text }]}> 
                  <Text style={[styles.sectionTitle, { color: palette.background }]}>{t('trends.compare.shared', 'Shared Trends')}</Text>
                </View>
                {sharedItems.length > 0 ? renderSectionRows({ items: sharedItems, palette }) : <Text style={[styles.emptySectionText, { color: palette.secondaryText }]}>{t('trends.no_data')}</Text>}
              </View>
            </View>
          ) : null}
        </Card>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  stack: { gap: 16 },
  label: {
    marginBottom: 8,
    fontSize: 14,
    letterSpacing: 1,
    textTransform: 'uppercase',
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
  list: {
    margin: -16,
    gap: 16,
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
