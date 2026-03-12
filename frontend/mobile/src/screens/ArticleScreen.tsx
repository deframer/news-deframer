import React, { useEffect, useMemo, useRef, useState } from 'react';
import { Image, Linking, Pressable, ScrollView, StyleSheet, Text, View } from 'react-native';
import { ExternalLink } from 'lucide-react-native';
import { useTranslation } from 'react-i18next';

import { Card } from '../components/Card';
import { getRatingColors, stripHtml, toPercent } from '../lib/articleFormat';
import { getRelativeTime } from '../lib/getRelativeTime';
import { AnalyzedItem } from '../services/newsDeframerClient';
import { AppPalette } from '../theme';

type DetailMode = 'closed' | 'original' | 'details';
type ScrollTarget = 'top' | 'overall' | 'original';

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
  const percentage = toPercent(value);
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

  const openOriginal = () => {
    Linking.openURL(item.url).catch(() => undefined);
  };

  return (
    <View style={[styles.container, { backgroundColor: palette.background }]}> 
      <ScrollView ref={scrollRef} style={styles.screen} contentContainerStyle={styles.content}>
        <Card palette={palette}>
        <Pressable
          onPress={openOriginal}
          style={styles.heroPressable}
        >
          {imageUrl ? (
            <View style={styles.imageWrap}>
              <Image source={{ uri: imageUrl }} style={styles.image} resizeMode="cover" />
              <Pressable
                accessibilityRole="button"
                accessibilityLabel={t('article.btn_open_article', 'Open article')}
                onPress={(event) => {
                  event?.stopPropagation?.();
                  openOriginal();
                }}
                style={[styles.openButton, { backgroundColor: palette.buttonBackground, borderColor: palette.buttonBorder }]}
              >
                <ExternalLink color={palette.buttonText} size={14} strokeWidth={2.2} />
                <Text style={[styles.openButtonText, { color: palette.buttonText }]}>{t('article.btn_open_article', 'Open article')}</Text>
              </Pressable>
            </View>
          ) : null}
          <View style={[styles.heroTextBlock, !imageUrl ? styles.heroTextBlockNoImage : null]}>
            {!imageUrl ? (
              <Pressable
                accessibilityRole="button"
                accessibilityLabel={t('article.btn_open_article', 'Open article')}
                onPress={(event) => {
                  event?.stopPropagation?.();
                  openOriginal();
                }}
                style={[styles.openButtonNoImage, { backgroundColor: palette.buttonBackground, borderColor: palette.buttonBorder }]}
              >
                <ExternalLink color={palette.buttonText} size={14} strokeWidth={2.2} />
                <Text style={[styles.openButtonText, { color: palette.buttonText }]}>{t('article.btn_open_article', 'Open article')}</Text>
              </Pressable>
            ) : null}
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
          <View style={[styles.originalSection, { borderTopColor: palette.border }]}>
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
  heroPressable: { gap: 0 },
  imageWrap: { position: 'relative', marginHorizontal: -16, marginTop: -16, marginBottom: 14 },
  image: { width: '100%', aspectRatio: 16 / 9 },
  openButton: {
    position: 'absolute',
    top: 10,
    right: 10,
    borderWidth: 1,
    borderRadius: 999,
    minHeight: 30,
    paddingHorizontal: 10,
    alignItems: 'center',
    justifyContent: 'center',
    flexDirection: 'row',
    gap: 6,
  },
  openButtonText: { fontSize: 12, fontWeight: '700' },
  openButtonNoImage: {
    position: 'absolute',
    top: 0,
    right: 0,
    borderWidth: 1,
    borderRadius: 999,
    minHeight: 30,
    paddingHorizontal: 10,
    alignItems: 'center',
    justifyContent: 'center',
    flexDirection: 'row',
    gap: 6,
  },
  heroTextBlock: { gap: 10 },
  heroTextBlockNoImage: {
    position: 'relative',
    paddingTop: 40,
  },
  title: { fontSize: 20, fontWeight: '700', lineHeight: 26 },
  description: { fontSize: 15, lineHeight: 22 },
  metaRow: { marginTop: 10, flexDirection: 'row', alignItems: 'center', gap: 8 },
  metaText: { fontSize: 13, fontWeight: '500' },
  metaDivider: { fontSize: 13, fontWeight: '600' },
  divider: { marginVertical: 18, borderBottomWidth: 1 },
  metricBlock: { gap: 8 },
  metricTitle: { fontSize: 18, fontWeight: '700', lineHeight: 24 },
  metricTrack: { height: 32, borderRadius: 999, overflow: 'hidden', justifyContent: 'center' },
  metricFill: { position: 'absolute', left: 0, top: 0, bottom: 0 },
  metricPercent: { paddingHorizontal: 12, fontSize: 14, fontWeight: '800' },
  metricReason: { fontSize: 15, lineHeight: 22 },
  metricsStack: { marginTop: 18, gap: 14 },
  originalSection: { marginTop: 20, gap: 10, borderTopWidth: 1, paddingTop: 14 },
  originalHeading: { fontSize: 18, fontWeight: '700', lineHeight: 24 },
  originalTitle: { fontSize: 20, fontWeight: '700', lineHeight: 26 },
  originalBody: { fontSize: 15, lineHeight: 22 },
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
  actionButtonText: { fontSize: 16, fontWeight: '700' },
});
