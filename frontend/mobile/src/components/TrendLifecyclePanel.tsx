import React, { useEffect, useMemo, useState } from 'react';
import { ArrowLeftRight } from 'lucide-react-native';
import { LayoutChangeEvent, Pressable, ScrollView, StyleSheet, Text, View } from 'react-native';
import { useTranslation } from 'react-i18next';

import { Lifecycle, NewsDeframerClient } from '../services/newsDeframerClient';
import { AnalyzedItem } from '../services/newsDeframerClient';
import { logger } from '../services/logger';
import { Settings } from '../services/settingsService';
import { AppPalette } from '../theme';
import { LoadingSpinner } from './LoadingSpinner';
import { TrendArticleListPanel } from './TrendArticleListPanel';

const toIsoDay = (value: string) => new Date(value).toISOString().split('T')[0];

const formatVelocity = (value: number) => `${value > 0 ? '+' : ''}${value}`;

const getTextColorForHex = (hex: string) => {
  const normalized = hex.replace('#', '');
  if (normalized.length !== 6) {
    return '#ffffff';
  }

  const r = parseInt(normalized.slice(0, 2), 16);
  const g = parseInt(normalized.slice(2, 4), 16);
  const b = parseInt(normalized.slice(4, 6), 16);
  const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255;

  return luminance > 0.62 ? '#000000' : '#ffffff';
};

export const TrendLifecyclePanel = ({
  palette,
  term,
  domain,
  language,
  daysInPast,
  settings,
  onOpenArticle,
}: {
  palette: AppPalette;
  term: string;
  domain: string;
  language: string;
  daysInPast: number;
  settings: Settings;
  onOpenArticle: (item: AnalyzedItem) => void;
}) => {
  const { t, i18n } = useTranslation();
  const [data, setData] = useState<Lifecycle[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedEntryKey, setSelectedEntryKey] = useState<string | null>(null);
  const [compactMode, setCompactMode] = useState(false);
  const [chartWidth, setChartWidth] = useState(320);

  const client = useMemo(() => new NewsDeframerClient(settings), [settings]);

  useEffect(() => {
    let cancelled = false;

    const fetchData = async () => {
      if (!term) {
        setData([]);
        return;
      }

      setLoading(true);
      try {
        const result = await client.getLifecycleByDomain(term, domain, language, daysInPast);
        const sorted = [...result].sort((a, b) => new Date(a.time_slice).getTime() - new Date(b.time_slice).getTime());

        if (!cancelled) {
          setData(sorted);
        }
      } catch {
        if (!cancelled) {
          logger.error('TrendLifecycle fetch failed', { term, domain, language, daysInPast });
          setData([]);
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    };

    fetchData();

    return () => {
      cancelled = true;
    };
  }, [client, daysInPast, domain, language, term]);

  useEffect(() => {
    setSelectedEntryKey(null);
  }, [term, domain]);

  useEffect(() => {
    if (!selectedEntryKey || loading) {
      return;
    }

    const isStillInRange = data.some((entry, index) => `${entry.time_slice}-${index}` === selectedEntryKey);

    if (!isStillInRange) {
      setSelectedEntryKey(null);
    }
  }, [data, loading, selectedEntryKey]);

  const maxFreq = data.length > 0 ? Math.max(...data.map((entry) => entry.frequency)) : 0;
  const barCount = Math.max(data.length, 1);
  const wideBarWidth = 64;
  const wideChartWidth = Math.max(320, barCount * wideBarWidth);
  const compactGap = 2;
  const compactBarWidth = Math.max(8, Math.floor((chartWidth - compactGap * Math.max(barCount - 1, 0)) / barCount));

  const selectedEntry = useMemo(() => {
    if (!selectedEntryKey) {
      return null;
    }

    return data.find((entry, index) => `${entry.time_slice}-${index}` === selectedEntryKey) ?? null;
  }, [data, selectedEntryKey]);

  if (loading) {
    return (
      <View style={styles.stateWrap}>
        <LoadingSpinner palette={palette} center label={t('options.loading')} />
      </View>
    );
  }

  if (data.length === 0) {
    return (
      <View style={styles.stateWrap}>
        <Text style={[styles.stateText, { color: palette.secondaryText }]}>{t('trends.lifecycle_no_data', 'No lifecycle data available for this topic.')}</Text>
      </View>
    );
  }

  const handleChartLayout = (event: LayoutChangeEvent) => {
    const width = Math.round(event.nativeEvent.layout.width);
    if (width > 0 && width !== chartWidth) {
      setChartWidth(width);
    }
  };

  return (
    <View style={styles.container}>
      <View style={styles.modeRow}>
        <Pressable
          onPress={() => setCompactMode((current) => !current)}
          style={[styles.modeButton, { borderColor: palette.buttonBorder, backgroundColor: palette.buttonBackground }]}
          accessibilityRole="button"
          accessibilityLabel={compactMode ? t('mobile.trends_chart_mode_wide') : t('mobile.trends_chart_mode_compact')}
        >
          <ArrowLeftRight color={palette.text} size={16} strokeWidth={2.2} />
        </Pressable>
      </View>

      {compactMode ? (
        <View style={styles.chartFrame} onLayout={handleChartLayout}>
          <View style={styles.compactRow}>
            {data.map((item, index) => {
              const dateLabel = new Date(item.time_slice).toLocaleDateString(i18n.language, { month: 'short', day: 'numeric' });
              const heightPercent = maxFreq > 0 ? (item.frequency / maxFreq) * 100 : 0;
              const barHeight = Math.max(8, Math.round((heightPercent / 100) * 180));
              const tone = item.velocity > 0 ? palette.trendUp : item.velocity < 0 ? palette.trendDown : palette.trendSteady;
              const icon = item.velocity > 0 ? '▲' : item.velocity < 0 ? '▼' : '▶';
              const entryKey = `${item.time_slice}-${index}`;
              const selected = selectedEntryKey === entryKey;

              return (
                <Pressable
                  key={entryKey}
                  onPress={() => setSelectedEntryKey((current) => (current === entryKey ? null : entryKey))}
                  style={[
                    styles.compactBarWrap,
                    { width: compactBarWidth },
                    selected ? [styles.selectedWrap, { borderColor: palette.accent }] : null,
                  ]}
                  accessibilityRole="button"
                  accessibilityState={{ selected }}
                  accessibilityLabel={t('trends.search_aria_label', '{{date}}: Frequency {{frequency}}, Velocity {{velocity}}', {
                    date: dateLabel,
                    frequency: item.frequency,
                    velocity: formatVelocity(item.velocity),
                  })}
                >
                  <Text style={[styles.barIcon, { color: tone }]}>{icon}</Text>
                  <View style={[styles.barBody, { height: barHeight, backgroundColor: tone }]} />
                </Pressable>
              );
            })}
          </View>
        </View>
      ) : (
        <View style={styles.chartFrame}>
          <ScrollView horizontal showsHorizontalScrollIndicator style={styles.chartScroll} contentContainerStyle={{ minWidth: wideChartWidth }}>
            <View style={styles.wideRow}>
              {data.map((item, idx) => {
                const dateLabel = new Date(item.time_slice).toLocaleDateString(i18n.language, { month: 'short', day: 'numeric' });
                const heightPercent = maxFreq > 0 ? (item.frequency / maxFreq) * 100 : 0;
                const barHeight = Math.max(8, Math.round((heightPercent / 100) * 180));
                const tone = item.velocity > 0 ? palette.trendUp : item.velocity < 0 ? palette.trendDown : palette.trendSteady;
                const icon = item.velocity > 0 ? '▲' : item.velocity < 0 ? '▼' : '▶';
                const barTextTone = getTextColorForHex(tone);
                const entryKey = `${item.time_slice}-${idx}`;
                const selected = selectedEntryKey === entryKey;
                const showLabel = data.length < 15 || idx % Math.ceil(data.length / 10) === 0;
                const showDate = showLabel && wideBarWidth >= 44;

                return (
                  <Pressable
                    key={entryKey}
                    onPress={() => setSelectedEntryKey((current) => (current === entryKey ? null : entryKey))}
                    style={[styles.wideBarWrap, selected ? [styles.selectedWrap, { borderColor: palette.accent }] : null]}
                    accessibilityRole="button"
                    accessibilityState={{ selected }}
                    accessibilityLabel={t('trends.search_aria_label', '{{date}}: Frequency {{frequency}}, Velocity {{velocity}}', {
                      date: dateLabel,
                      frequency: item.frequency,
                      velocity: formatVelocity(item.velocity),
                    })}
                  >
                    <View style={styles.barTopMeta}>
                      <Text style={[styles.barTopMetaText, { color: palette.secondaryText }]}>{t('trends.freq', 'Freq')}: {item.frequency}</Text>
                      <Text style={[styles.barTopMetaText, { color: palette.secondaryText }]}>{t('trends.vel', 'Vel')}: {formatVelocity(item.velocity)}</Text>
                    </View>
                    <Text style={[styles.barIcon, { color: tone }]}>{icon}</Text>
                    <View style={[styles.barBody, styles.barBodyContent, { height: barHeight, backgroundColor: tone }]}> 
                      {showDate ? <Text style={[styles.barDateText, { color: barTextTone }]} numberOfLines={1}>{dateLabel}</Text> : null}
                    </View>
                  </Pressable>
                );
              })}
            </View>
          </ScrollView>
        </View>
      )}

      {selectedEntry ? <TrendArticleListPanel palette={palette} term={term} domain={domain} settings={settings} selectedDate={toIsoDay(selectedEntry.time_slice)} onOpenArticle={onOpenArticle} /> : null}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginTop: 14,
  },
  modeRow: {
    marginBottom: 8,
    flexDirection: 'row',
    justifyContent: 'flex-start',
  },
  modeButton: {
    width: 32,
    height: 32,
    borderWidth: 1,
    borderRadius: 8,
    alignItems: 'center',
    justifyContent: 'center',
  },
  stateWrap: {
    marginTop: 14,
    minHeight: 200,
    justifyContent: 'center',
    alignItems: 'center',
  },
  stateText: {
    fontSize: 15,
    textAlign: 'center',
  },
  chartFrame: {
    minHeight: 220,
  },
  chartScroll: {
    width: '100%',
  },
  wideRow: {
    minHeight: 220,
    flexDirection: 'row',
    alignItems: 'flex-end',
    gap: 4,
    paddingBottom: 8,
  },
  compactRow: {
    minHeight: 220,
    flexDirection: 'row',
    alignItems: 'flex-end',
    gap: 2,
    paddingBottom: 8,
  },
  wideBarWrap: {
    width: 60,
    minHeight: 212,
    alignItems: 'center',
    justifyContent: 'flex-end',
    borderRadius: 6,
    paddingTop: 4,
    paddingBottom: 4,
  },
  compactBarWrap: {
    minHeight: 212,
    alignItems: 'center',
    justifyContent: 'flex-end',
    borderRadius: 6,
    paddingTop: 4,
    paddingBottom: 4,
  },
  selectedWrap: {
    borderWidth: 1,
  },
  barIcon: {
    marginBottom: 4,
    fontSize: 12,
    fontWeight: '700',
  },
  barBody: {
    width: '100%',
    borderRadius: 4,
    minHeight: 8,
  },
  barBodyContent: {
    justifyContent: 'flex-end',
    alignItems: 'center',
    paddingHorizontal: 4,
    paddingBottom: 6,
  },
  barDateText: {
    fontSize: 11,
    fontWeight: '700',
    lineHeight: 13,
    textAlign: 'center',
  },
  barTopMeta: {
    minHeight: 28,
    width: '100%',
    alignItems: 'center',
    justifyContent: 'flex-start',
  },
  barTopMetaText: {
    fontSize: 10,
    fontWeight: '700',
    lineHeight: 11,
    textAlign: 'center',
  },
});
