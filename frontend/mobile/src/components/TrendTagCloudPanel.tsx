import React, { useEffect, useMemo, useState } from 'react';
import { LayoutChangeEvent, Pressable, StyleSheet, Text, View } from 'react-native';
import { useTranslation } from 'react-i18next';

import { buildWordCloudWords, getWordCloudColor, layoutWordCloud, MeasureWordOptions } from '../../../shared/wordcloud';

import { TrendMetric, NewsDeframerClient } from '../services/newsDeframerClient';
import { logger } from '../services/logger';
import { Settings } from '../services/settingsService';
import { AppPalette } from '../theme';
import { Card } from './Card';
import { LoadingSpinner } from './LoadingSpinner';
import { TrendDetailsPanel } from './TrendDetailsPanel';

type TrendDetailTab = 'lifecycle' | 'context' | 'articles';

const DEFAULT_WIDTH = 640;
const DEFAULT_HEIGHT = 400;

const approximateMeasureWord = ({ text, fontSize }: MeasureWordOptions) => ({
  width: Math.max(fontSize, text.length * fontSize * 0.58),
  height: fontSize * 1.2,
});

const createMeasureWord = () => {
  if (typeof document === 'undefined') {
    return approximateMeasureWord;
  }

  const canvas = document.createElement('canvas');
  let context: CanvasRenderingContext2D | null = null;

  try {
    context = canvas.getContext('2d');
  } catch {
    return approximateMeasureWord;
  }

  if (!context) {
    return approximateMeasureWord;
  }

  return ({ text, fontSize, fontFamily, fontWeight = 'normal' }: MeasureWordOptions) => {
    context.font = `${fontWeight} ${fontSize}px ${fontFamily}`;
    const metrics = context.measureText(text);

    return {
      width: Math.max(metrics.width, fontSize),
      height: fontSize * 1.2,
    };
  };
};

export const TrendTagCloudPanel = ({
  palette,
  domain,
  language,
  daysInPast,
  settings,
}: {
  palette: AppPalette;
  domain: string;
  language: string;
  daysInPast: number;
  settings: Settings;
}) => {
  const { t } = useTranslation();
  const [items, setItems] = useState<TrendMetric[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [hasError, setHasError] = useState(false);
  const [selectedTerm, setSelectedTerm] = useState<string | null>(null);
  const [activeDetailTab, setActiveDetailTab] = useState<TrendDetailTab>('lifecycle');
  const [cloudWidth, setCloudWidth] = useState<number>(DEFAULT_WIDTH);

  const client = useMemo(() => new NewsDeframerClient(settings), [settings]);
  const measureWord = useMemo(() => createMeasureWord(), []);
  const cloudHeight = DEFAULT_HEIGHT;
  const trendColors = useMemo(() => [palette.trendUp, palette.trendDown, palette.trendSteady] as const, [palette.trendDown, palette.trendSteady, palette.trendUp]);

  useEffect(() => {
    let cancelled = false;

    const fetchTopTrends = async () => {
      setIsLoading(true);
      setHasError(false);
      logger.info('TrendTagCloud fetch start', { domain, language, daysInPast });

      try {
        const data = await client.getTopTrendByDomain(domain, language, daysInPast);
        if (!cancelled) {
          setItems(data);
          logger.info('TrendTagCloud fetch success', { domain, language, daysInPast, itemCount: data.length });
        }
      } catch {
        if (!cancelled) {
          setItems([]);
          setHasError(true);
          logger.error('TrendTagCloud fetch failed', { domain, language, daysInPast });
        }
      } finally {
        if (!cancelled) {
          setIsLoading(false);
        }
      }
    };

    fetchTopTrends();

    return () => {
      cancelled = true;
    };
  }, [client, daysInPast, domain, language]);

  const words = useMemo(
    () =>
      buildWordCloudWords(
        items,
        (item) => item.trend_topic,
        (item) => item.outlier_ratio,
        (item, index) => ({ ...item, rank: index + 1 }),
      ),
    [items],
  );

  const positionedWords = useMemo(
    () =>
      layoutWordCloud({
        words,
        width: cloudWidth,
        height: cloudHeight,
        measureWord,
        getColor: (text) => getWordCloudColor(text, trendColors),
      }),
    [cloudHeight, cloudWidth, measureWord, trendColors, words],
  );

  const cloudBounds = useMemo(() => {
    if (positionedWords.length === 0) {
      return null;
    }

    return positionedWords.reduce(
      (acc, word) => ({
        left: Math.min(acc.left, word.x - word.width / 2),
        top: Math.min(acc.top, word.y - word.height / 2),
        right: Math.max(acc.right, word.x + word.width / 2),
        bottom: Math.max(acc.bottom, word.y + word.height / 2),
      }),
      { left: Infinity, top: Infinity, right: -Infinity, bottom: -Infinity },
    );
  }, [positionedWords]);

  const isSelectedTermVisible = useMemo(() => {
    if (!selectedTerm) {
      return false;
    }

    return positionedWords.some((word) => word.text === selectedTerm);
  }, [positionedWords, selectedTerm]);

  useEffect(() => {
    setSelectedTerm(null);
    setActiveDetailTab('lifecycle');
  }, [daysInPast, domain, language]);

  const handleCloudLayout = (event: LayoutChangeEvent) => {
    const width = Math.round(event.nativeEvent.layout.width);
    if (width > 0 && Math.abs(width - cloudWidth) > 1) {
      logger.info('TrendTagCloud layout width update', { prevWidth: cloudWidth, nextWidth: width });
      setCloudWidth(width);
    }
  };

  const handleSelectTerm = (term: string) => {
    setSelectedTerm((current) => {
      const next = current === term ? null : term;
      logger.info('TrendTagCloud term selection changed', { previous: current, next });
      return next;
    });
  };

  return (
    <View style={styles.stack}>
      <Card palette={palette}>
        <View style={[styles.cloudFrame, { height: cloudHeight }]} onLayout={handleCloudLayout}>
          {isLoading ? <LoadingSpinner palette={palette} center label={t('options.loading')} /> : null}
          {!isLoading && hasError ? <Text style={[styles.stateText, { color: palette.secondaryText }]}>{t('mobile.portal_load_error')}</Text> : null}
          {!isLoading && !hasError && items.length === 0 ? <Text style={[styles.stateText, { color: palette.secondaryText }]}>{t('trends.no_data')}</Text> : null}

          {!isLoading && !hasError && cloudBounds ? (
            <View
              style={[
                styles.cloudWrap,
                {
                  width: cloudBounds.right - cloudBounds.left,
                  height: cloudBounds.bottom - cloudBounds.top,
                },
              ]}
            >
              {positionedWords.map((word) => {
                const selected = selectedTerm === word.text;
                const toneStyle = selected ? styles.termSelected : styles.termDefault;
                const placementStyle = {
                  left: word.x - cloudBounds.left - Math.ceil(word.width) / 2,
                  top: word.y - cloudBounds.top - Math.ceil(word.height) / 2,
                };
                const colorStyle = {
                  color: selected ? palette.accent : word.color,
                  fontSize: word.fontSize,
                  lineHeight: Math.ceil(word.fontSize * 1.12),
                };
                const tooltipLabel = `${t('trends.rank')}: ${word.rank}. ${t('trends.trend')}: ${word.data.outlier_ratio.toFixed(2)}x. ${t('trends.vol')}: ${word.data.frequency}.`;

                return (
                  <Pressable
                    key={word.text}
                    onPress={() => handleSelectTerm(word.text)}
                    style={[styles.termHitBox, styles.termAbsolute, placementStyle]}
                    accessibilityRole="button"
                    accessibilityState={{ selected }}
                    accessibilityLabel={`${word.text}. ${tooltipLabel}`}
                  >
                    <Text style={[styles.term, toneStyle, colorStyle]}>{word.text}</Text>
                  </Pressable>
                );
              })}
            </View>
          ) : null}
        </View>
      </Card>

      {selectedTerm && isSelectedTermVisible ? (
        <TrendDetailsPanel
          palette={palette}
          term={selectedTerm}
          activeTab={activeDetailTab}
          setActiveTab={setActiveDetailTab}
        />
      ) : null}
    </View>
  );
};

const styles = StyleSheet.create({
  stack: { gap: 16 },
  cloudFrame: {
    marginTop: 8,
    width: '100%',
    position: 'relative',
    alignItems: 'center',
    justifyContent: 'center',
    overflow: 'hidden',
  },
  cloudWrap: {
    position: 'relative',
  },
  stateText: {
    fontSize: 15,
    textAlign: 'center',
  },
  termHitBox: { alignItems: 'flex-start', justifyContent: 'flex-start', padding: 0, margin: 0 },
  termAbsolute: {
    position: 'absolute',
  },
  term: { lineHeight: 24 },
  termDefault: { fontWeight: '500' },
  termSelected: { textDecorationLine: 'underline', fontWeight: '700' },
});
