import React, { useEffect, useMemo, useRef, useState } from 'react';
import { Image, Linking, Pressable, ScrollView, StyleSheet, Text, View } from 'react-native';
import { useTranslation } from 'react-i18next';

import { Card } from '../components/Card';
import { AnalyzedItem } from '../services/newsDeframerClient';
import { AppPalette } from '../theme';

type DetailMode = 'closed' | 'original' | 'details';
type ScrollTarget = 'top' | 'overall' | 'original';

const stripHtml = (value?: string): string => {
  if (!value) {
    return '';
  }

  return value.replace(/<[^>]+>/g, ' ').replace(/\s+/g, ' ').trim();
};

const formatPercent = (value?: number): number => Math.round((value || 0) * 100);

const getRatingColors = (percentage: number, palette: AppPalette) => {
  if (percentage < 11) {
    return { backgroundColor: palette.success, textColor: palette.text };
  }

  if (percentage < 34) {
    return { backgroundColor: palette.success, textColor: '#ffffff' };
  }

  if (percentage < 67) {
    return { backgroundColor: palette.warning, textColor: '#000000' };
  }

  return { backgroundColor: palette.danger, textColor: '#ffffff' };
};

const getRelativeTime = (dateStr: string | Date, locale: string): string => {
  const date = new Date(dateStr);
  if (Number.isNaN(date.getTime())) {
    return '';
  }

  const seconds = (Date.now() - date.getTime()) / 1000;
  const absoluteSeconds = Math.abs(seconds);
  const rtf = new Intl.RelativeTimeFormat(locale, { style: 'narrow' });

  if (absoluteSeconds < 60) return '';
  if (absoluteSeconds < 3600) return rtf.format(-Math.round(seconds / 60), 'minute');
  if (absoluteSeconds < 86400) return rtf.format(-Math.round(seconds / 3600), 'hour');
  if (absoluteSeconds < 604800) return rtf.format(-Math.round(seconds / 86400), 'day');
  if (absoluteSeconds < 2592000) return rtf.format(-Math.round(seconds / 604800), 'week');
  if (absoluteSeconds < 31536000) return rtf.format(-Math.round(seconds / 2592000), 'month');

  return rtf.format(-Math.round(seconds / 31536000), 'year');
};

const RatingRow = ({
  palette,
  label,
  value,
  reason,
}: {
  palette: AppPalette;
  label: string;
  value?: number;
  reason?: string;
}) => {
  const { t } = useTranslation();
  const percentage = formatPercent(value);
  const colors = getRatingColors(percentage, palette);

  return (
    <View style={styles.metricBlock}>
      <Text style={[styles.metricTitle, { color: palette.text }]}>{label}</Text>
      <View style={[styles.metricTrack, { backgroundColor: palette.ratingBackground }]}> 
        <View style={[styles.metricFill, { width: `${percentage}%`, backgroundColor: colors.backgroundColor }]} />
        <Text style={[styles.metricPercent, { color: colors.textColor }]}>{percentage}%</Text>
      </View>
      <Text style={[styles.metricReason, { color: palette.secondaryText }]}>{reason || t('rating.no_reason')}</Text>
    </View>
  );
};

export const ArticleScreen = ({ palette, item }: { palette: AppPalette; item: AnalyzedItem }) => {
  const { i18n, t } = useTranslation();
  const [mode, setMode] = useState<DetailMode>('closed');
  const [overallY, setOverallY] = useState<number | null>(null);
  const [pendingScrollTarget, setPendingScrollTarget] = useState<ScrollTarget | null>(null);
  const scrollRef = useRef<ScrollView>(null);

  const title = item.title_corrected || stripHtml(item.title_original) || t('article.no_title', 'No Title');
  const description = item.description_corrected || stripHtml(item.description_original) || t('article.no_description', 'No Description');
  const imageUrl = item.media?.medium === 'image' ? item.media.url : undefined;
  const author = item.authors?.join(', ');
  const timeAgo = item.pubDate ? getRelativeTime(item.pubDate, i18n.language) || t('metadata.just_now') : '';

  const metrics = useMemo(
    () => [
      { id: 'framing', label: t('metrics.framing'), value: item.framing, reason: item.framing_reason },
      { id: 'clickbait', label: t('metrics.clickbait'), value: item.clickbait, reason: item.clickbait_reason },
      { id: 'persuasive', label: t('metrics.persuasive'), value: item.persuasive, reason: item.persuasive_reason },
      { id: 'hyper_stimulus', label: t('metrics.hyper_stimulus'), value: item.hyper_stimulus, reason: item.hyper_stimulus_reason },
      { id: 'speculative', label: t('metrics.speculative'), value: item.speculative, reason: item.speculative_reason },
    ],
    [
      item.clickbait,
      item.clickbait_reason,
      item.framing,
      item.framing_reason,
      item.hyper_stimulus,
      item.hyper_stimulus_reason,
      item.persuasive,
      item.persuasive_reason,
      item.speculative,
      item.speculative_reason,
      t,
    ],
  );

  const toggleMode = (next: Exclude<DetailMode, 'closed'>) => {
    setMode((current) => {
      if (current === next) {
        setPendingScrollTarget('top');
        return 'closed';
      }

      setPendingScrollTarget(next === 'details' ? 'overall' : 'original');
      return next;
    });
  };

  useEffect(() => {
    if (!pendingScrollTarget) {
      return;
    }

    const frame = requestAnimationFrame(() => {
      if (pendingScrollTarget === 'top') {
        scrollRef.current?.scrollTo({ y: 0, animated: true });
      } else if (pendingScrollTarget === 'overall') {
        if (overallY === null) {
          return;
        }
        scrollRef.current?.scrollTo({ y: Math.max(0, overallY - 8), animated: true });
      } else {
        scrollRef.current?.scrollToEnd({ animated: true });
      }
      setPendingScrollTarget(null);
    });

    return () => cancelAnimationFrame(frame);
  }, [overallY, pendingScrollTarget]);

  return (
    <View style={[styles.container, { backgroundColor: palette.background }]}> 
      <ScrollView ref={scrollRef} style={styles.screen} contentContainerStyle={styles.content}>
        <Card palette={palette}>
        <Pressable
          onPress={() => {
            Linking.openURL(item.url).catch(() => undefined);
          }}
          style={styles.heroPressable}
        >
          {imageUrl ? <Image source={{ uri: imageUrl }} style={styles.image} resizeMode="cover" /> : null}
          <View style={styles.heroTextBlock}>
            <Text style={[styles.title, { color: palette.text }]}>{title}</Text>
            <Text style={[styles.description, { color: palette.secondaryText }]}>{description}</Text>
          </View>
        </Pressable>

        {timeAgo || author ? (
          <View style={styles.metaRow}>
            {timeAgo ? <Text style={[styles.metaText, { color: palette.secondaryText }]}>{timeAgo}</Text> : null}
            {timeAgo && author ? <Text style={[styles.metaDivider, { color: palette.secondaryText }]}>|</Text> : null}
            {author ? <Text style={[styles.metaText, { color: palette.secondaryText }]}>{author}</Text> : null}
          </View>
        ) : null}

        <View style={[styles.divider, { borderBottomColor: palette.border }]} />

        <View onLayout={(event) => setOverallY(event.nativeEvent.layout.y)}>
          <RatingRow palette={palette} label={t('metrics.overall_rating')} value={item.rating} reason={item.overall_reason} />
        </View>

        {mode === 'details' ? (
          <View style={styles.metricsStack}>
            {metrics.map((metric) => (
              <RatingRow key={metric.id} palette={palette} label={metric.label} value={metric.value} reason={metric.reason} />
            ))}
          </View>
        ) : null}

        {mode === 'original' || mode === 'details' ? (
          <View style={styles.originalSection}>
            <Text style={[styles.originalHeading, { color: palette.text }]}>{t('article.original_section')}</Text>
            <Text style={[styles.originalTitle, { color: palette.text }]}>{stripHtml(item.title_original || '')}</Text>
            <Text style={[styles.originalBody, { color: palette.secondaryText }]}>{stripHtml(item.description_original || '')}</Text>
          </View>
        ) : null}
        </Card>
      </ScrollView>

      <View style={[styles.actionsDock, { backgroundColor: palette.background, borderTopColor: palette.border }]}> 
        <View style={styles.actionsRow}>
          <Pressable
            onPress={() => toggleMode('original')}
            style={[
              styles.actionButton,
              { borderColor: palette.buttonBorder, backgroundColor: mode === 'original' ? palette.accent : palette.buttonBackground },
            ]}
          >
            <Text style={[styles.actionButtonText, { color: mode === 'original' ? palette.accentText : palette.buttonText }]}>{t('article.btn_original_title')}</Text>
          </Pressable>
          <Pressable
            onPress={() => toggleMode('details')}
            style={[
              styles.actionButton,
              { borderColor: palette.buttonBorder, backgroundColor: mode === 'details' ? palette.accent : palette.buttonBackground },
            ]}
          >
            <Text style={[styles.actionButtonText, { color: mode === 'details' ? palette.accentText : palette.buttonText }]}>{t('article.btn_details')}</Text>
          </Pressable>
        </View>
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: { flex: 1 },
  screen: { flex: 1 },
  content: { padding: 24, paddingBottom: 110 },
  heroPressable: { gap: 14 },
  image: { width: '100%', aspectRatio: 16 / 9, borderRadius: 10 },
  heroTextBlock: { gap: 10 },
  title: { fontSize: 34, fontWeight: '800', lineHeight: 40 },
  description: { fontSize: 18, lineHeight: 28 },
  metaRow: { marginTop: 10, flexDirection: 'row', alignItems: 'center', gap: 8 },
  metaText: { fontSize: 15, fontWeight: '500' },
  metaDivider: { fontSize: 15, fontWeight: '700' },
  divider: { marginVertical: 18, borderBottomWidth: 1 },
  metricBlock: { gap: 8 },
  metricTitle: { fontSize: 20, fontWeight: '700' },
  metricTrack: { height: 32, borderRadius: 999, overflow: 'hidden', justifyContent: 'center' },
  metricFill: { position: 'absolute', left: 0, top: 0, bottom: 0 },
  metricPercent: { paddingHorizontal: 12, fontSize: 18, fontWeight: '800' },
  metricReason: { fontSize: 17, lineHeight: 24 },
  metricsStack: { marginTop: 18, gap: 14 },
  originalSection: { marginTop: 20, gap: 10 },
  originalHeading: { fontSize: 30, fontWeight: '800' },
  originalTitle: { fontSize: 26, fontWeight: '700', lineHeight: 34 },
  originalBody: { fontSize: 20, lineHeight: 30 },
  actionsDock: {
    position: 'absolute',
    left: 0,
    right: 0,
    bottom: 0,
    borderTopWidth: 1,
    paddingHorizontal: 24,
    paddingTop: 12,
    paddingBottom: 12,
  },
  actionsRow: { flexDirection: 'row', gap: 12 },
  actionButton: {
    flex: 1,
    borderWidth: 1,
    borderRadius: 10,
    minHeight: 52,
    alignItems: 'center',
    justifyContent: 'center',
    paddingHorizontal: 12,
  },
  actionButtonText: { fontSize: 18, fontWeight: '700' },
});
