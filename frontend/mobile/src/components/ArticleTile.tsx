import React from 'react';
import { Image, Pressable, StyleSheet, Text, View } from 'react-native';
import { Info } from 'lucide-react-native';
import { useTranslation } from 'react-i18next';

import { AnalyzedItem } from '../services/newsDeframerClient';
import { AppPalette } from '../theme';

const stripHtml = (value?: string): string => {
  if (!value) {
    return '';
  }

  return value.replace(/<[^>]+>/g, ' ').replace(/\s+/g, ' ').trim();
};

const formatRatingPercent = (rating?: number): number => Math.round((rating || 0) * 100);

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

export const ArticleTile = ({
  item,
  palette,
  onOpenArticle,
  onShowReason,
}: {
  item: AnalyzedItem;
  palette: AppPalette;
  onOpenArticle: (item: AnalyzedItem) => void;
  onShowReason: (reason: string) => void;
}) => {
  const { i18n, t } = useTranslation();
  const title = item.title_corrected || stripHtml(item.title_original) || t('article.no_title', 'No Title');
  const description = item.description_corrected || stripHtml(item.description_original) || t('article.no_description', 'No Description');
  const imageUrl = item.media?.medium === 'image' ? item.media.url : undefined;
  const rating = formatRatingPercent(item.rating);
  const ratingColors = getRatingColors(rating, palette);
  const ratingReason = item.overall_reason || t('rating.no_reason', 'No reason provided');
  const timeAgo = item.pubDate ? getRelativeTime(item.pubDate, i18n.language) || t('metadata.just_now', 'a moment ago') : '';
  const author = item.authors?.join(', ');

  return (
    <Pressable onPress={() => onOpenArticle(item)} style={[styles.tile, { backgroundColor: palette.card, borderColor: palette.border }]}>
      <View style={styles.ratingRow}>
        <View style={[styles.ratingTrack, { backgroundColor: palette.ratingBackground }]}>
          <View style={[styles.ratingFill, { width: `${rating}%`, backgroundColor: ratingColors.backgroundColor }]} />
          <Text style={[styles.ratingText, { color: ratingColors.textColor }]}>{rating}%</Text>
        </View>
        <Pressable
          accessibilityLabel={t('rating.info_aria_label', 'Show rating reason')}
          onPress={(event) => {
            event.stopPropagation();
            onShowReason(ratingReason);
          }}
          style={[styles.infoButton, { borderColor: palette.border, backgroundColor: palette.buttonBackground }]}
        >
          <Info color={palette.text} size={18} strokeWidth={2.2} />
        </Pressable>
      </View>

      {imageUrl ? <Image source={{ uri: imageUrl }} style={styles.image} resizeMode="cover" /> : null}

      <View style={styles.content}>
        <Text style={[styles.title, { color: palette.text }]}>{title}</Text>
        <Text style={[styles.description, { color: palette.secondaryText }]} numberOfLines={4}>{description}</Text>
      </View>

      {timeAgo || author ? (
        <View style={styles.metaRow}>
          {timeAgo ? <Text style={[styles.metaText, { color: palette.secondaryText }]}>{timeAgo}</Text> : null}
          {timeAgo && author ? <Text style={[styles.metaSeparator, { color: palette.secondaryText }]}>|</Text> : null}
          {author ? <Text style={[styles.metaText, { color: palette.secondaryText }]} numberOfLines={1}>{author}</Text> : null}
        </View>
      ) : null}
    </Pressable>
  );
};

const styles = StyleSheet.create({
  tile: {
    borderWidth: 1,
    borderRadius: 12,
    overflow: 'hidden',
  },
  ratingRow: {
    paddingHorizontal: 14,
    paddingTop: 14,
    flexDirection: 'row',
    alignItems: 'center',
    gap: 10,
  },
  ratingTrack: {
    flex: 1,
    height: 22,
    borderRadius: 999,
    overflow: 'hidden',
    justifyContent: 'center',
  },
  infoButton: {
    width: 36,
    height: 36,
    borderWidth: 1,
    borderRadius: 999,
    alignItems: 'center',
    justifyContent: 'center',
    flexShrink: 0,
  },
  ratingFill: {
    position: 'absolute',
    left: 0,
    top: 0,
    bottom: 0,
  },
  ratingText: {
    paddingHorizontal: 10,
    fontSize: 12,
    fontWeight: '700',
  },
  image: {
    width: '100%',
    aspectRatio: 16 / 9,
    marginTop: 14,
  },
  content: {
    paddingHorizontal: 14,
    paddingTop: 14,
    gap: 8,
  },
  title: {
    fontSize: 20,
    fontWeight: '700',
    lineHeight: 26,
  },
  description: {
    fontSize: 15,
    lineHeight: 22,
  },
  metaRow: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 14,
    paddingTop: 12,
    paddingBottom: 14,
    gap: 8,
  },
  metaText: {
    flexShrink: 1,
    fontSize: 13,
    fontWeight: '500',
  },
  metaSeparator: {
    fontSize: 13,
    fontWeight: '600',
  },
});
