import React, { useEffect, useMemo, useState } from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { useTranslation } from 'react-i18next';

import { NewsDeframerClient, TrendContext } from '../services/newsDeframerClient';
import { logger } from '../services/logger';
import { Settings } from '../services/settingsService';
import { AppPalette } from '../theme';
import { LoadingSpinner } from './LoadingSpinner';

export const TrendContextPanel = ({
  palette,
  term,
  domain,
  language,
  daysInPast,
  settings,
}: {
  palette: AppPalette;
  term: string;
  domain: string;
  language: string;
  daysInPast: number;
  settings: Settings;
}) => {
  const { t } = useTranslation();
  const [items, setItems] = useState<TrendContext[]>([]);
  const [loading, setLoading] = useState(true);

  const client = useMemo(() => new NewsDeframerClient(settings), [settings]);

  useEffect(() => {
    let mounted = true;

    const fetchData = async () => {
      if (!term) {
        setItems([]);
        setLoading(false);
        return;
      }

      setLoading(true);
      try {
        const data = await client.getContextByDomain(term, domain, language, daysInPast);
        if (mounted) {
          setItems(data);
        }
      } catch {
        if (mounted) {
          logger.error('TrendContext fetch failed', { term, domain, language, daysInPast });
          setItems([]);
        }
      } finally {
        if (mounted) {
          setLoading(false);
        }
      }
    };

    fetchData();

    return () => {
      mounted = false;
    };
  }, [client, daysInPast, domain, language, term]);

  if (loading) {
    return (
      <View style={styles.stateWrap}>
        <LoadingSpinner palette={palette} center label={t('options.loading')} />
      </View>
    );
  }

  if (items.length === 0) {
    return (
      <View style={styles.stateWrap}>
        <Text style={[styles.stateText, { color: palette.secondaryText }]}>{t('trends.context_no_data', 'No context data available for this topic.')}</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <View style={[styles.tableWrap, { borderColor: palette.border }]}> 
        <View style={[styles.headerRow, { borderBottomColor: palette.border, backgroundColor: palette.secondaryBackground }]}>
          <Text style={[styles.headerCellFrequency, styles.headerCellBase, { color: palette.secondaryText }]}>{t('trends.frequency_label', 'Frequency')}</Text>
          <Text style={[styles.headerCellVerb, styles.headerCellBase, { color: palette.secondaryText }]}>{t('trends.verb_label', 'Verb')}</Text>
        </View>
        {items.map((item) => (
          <View key={item.context} style={[styles.dataRow, { borderTopColor: palette.border }]}> 
            <Text style={[styles.frequencyCell, { color: palette.secondaryText }]}>{item.frequency}</Text>
            <Text style={[styles.verbCell, { color: palette.text }]}>{item.context}</Text>
          </View>
        ))}
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginTop: 14,
  },
  stateWrap: {
    minHeight: 80,
    justifyContent: 'center',
    alignItems: 'center',
  },
  stateText: {
    fontSize: 15,
    textAlign: 'center',
  },
  tableWrap: {
    borderWidth: 1,
    borderRadius: 8,
    overflow: 'hidden',
  },
  headerRow: {
    flexDirection: 'row',
    alignItems: 'center',
    borderBottomWidth: 1,
    minHeight: 38,
    paddingHorizontal: 10,
  },
  headerCellBase: {
    fontSize: 12,
    fontWeight: '700',
    paddingVertical: 8,
  },
  headerCellFrequency: {
    width: 86,
  },
  headerCellVerb: {
    flex: 1,
  },
  dataRow: {
    flexDirection: 'row',
    alignItems: 'center',
    minHeight: 36,
    paddingHorizontal: 10,
    borderTopWidth: 1,
  },
  frequencyCell: {
    width: 86,
    fontSize: 13,
    fontWeight: '600',
  },
  verbCell: {
    flex: 1,
    fontSize: 14,
    fontWeight: '600',
  },
});
