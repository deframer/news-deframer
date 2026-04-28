import React from 'react';
import { Image, Pressable, StyleSheet, Text, View } from 'react-native';
import { Info } from './icons';
import { useTranslation } from 'react-i18next';

import { formatTime } from '../../../shared/formatTime';
import { getRatingColors, stripHtml, toPercent } from '../lib/articleFormat';
import { AnalyzedItem } from '../services/newsDeframerClient';
import { logger } from '../services/logger';
import { AppPalette } from '../theme';

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
  const rating = toPercent(item.rating);
  const ratingColors = getRatingColors(rating, palette);
  const ratingReason = item.overall_reason || t('rating.no_reason', 'No reason provided');
  const timeAgo = item.pubDate ? formatTime(item.pubDate, i18n.language) : '';
  const author = item.authors?.join(', ');
  const openArticle = () => {
    logger.info('ArticleTile open article pressed', { url: item.url, title });
    onOpenArticle(item);
  };
  const showReason = () => {
    logger.info('ArticleTile reason pressed', { url: item.url, title, hasReason: Boolean(item.overall_reason) });
    onShowReason(ratingReason);
  };

  return (
    <View style={[styles.tile, { backgroundColor: palette.card, borderColor: palette.border }]}> 
      <View style={styles.ratingRow}>
        <Pressable onPress={openArticle} style={styles.ratingTrackWrap}>
          <View style={[styles.ratingTrack, { backgroundColor: palette.ratingBackground }]}> 
            <View style={[styles.ratingFill, { width: `${rating}%`, backgroundColor: ratingColors.backgroundColor }]} />
            <Text style={[styles.ratingText, { color: ratingColors.textColor }]}>{rating}%</Text>
          </View>
        </Pressable>
        <Pressable
          accessibilityLabel={t('rating.info_aria_label', 'Show rating reason')}
          onPress={showReason}
          style={[styles.infoButton, { borderColor: palette.border, backgroundColor: palette.buttonBackground }]}
        >
          <Info color={palette.text} size={18} strokeWidth={2.2} />
        </Pressable>
      </View>

      <Pressable onPress={openArticle}>
        {imageUrl ? <Image source={{ uri: imageUrl }} style={styles.image} resizeMode="cover" /> : null}

        <View style={styles.content}>
          <Text style={[styles.title, { color: palette.text }]}>{title}</Text>
          <Text style={[styles.description, { color: palette.secondaryText }]} numberOfLines={4}>{description}</Text>
        </View>

        {timeAgo || author ? (
          <View style={styles.metaRow}>
            {timeAgo ? <Text style={[styles.metaText, { color: palette.secondaryText }]}>{timeAgo}</Text> : null}
            {timeAgo && author ? <Text style={[styles.metaSeparator, { color: palette.secondaryText }]}>|</Text> : null}
            {author ? <Text style={[styles.metaText, styles.metaAuthorText, { color: palette.secondaryText }]} numberOfLines={1}>{author}</Text> : null}
          </View>
        ) : null}
      </Pressable>
    </View>
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
  ratingTrackWrap: {
    flex: 1,
  },
  ratingTrack: {
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
  metaAuthorText: {
    flexShrink: 1,
  },
});
