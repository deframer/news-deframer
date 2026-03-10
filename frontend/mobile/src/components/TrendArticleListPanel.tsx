import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { useTranslation } from 'react-i18next';

import { AppPalette } from '../theme';

export const TrendArticleListPanel = ({ palette, term }: { palette: AppPalette; term: string }) => {
  const { t } = useTranslation();

  return (
    <View style={[styles.container, { borderColor: palette.border }]}> 
      <View style={styles.headerRow}>
        <Text style={[styles.headerCellRating, styles.headerCellBase, { color: palette.secondaryText }]}>{t('trends.rating_caption', 'Rating')}</Text>
        <Text style={[styles.headerCellDate, styles.headerCellBase, { color: palette.secondaryText }]}>{t('trends.date_caption', 'Date')}</Text>
        <Text style={[styles.headerCellAuthor, styles.headerCellBase, { color: palette.secondaryText }]}>{t('trends.author_caption', 'Author')}</Text>
        <Text style={[styles.headerCellArticle, styles.headerCellBase, { color: palette.secondaryText }]}>{t('trends.article_caption', 'Article')}</Text>
      </View>

      <View style={[styles.placeholderBody, { backgroundColor: palette.secondaryBackground }]}> 
        <Text style={[styles.placeholderText, { color: palette.secondaryText }]}>{t('mobile.trends_article_list_placeholder', { term })}</Text>
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginTop: 14,
    borderWidth: 1,
    borderRadius: 8,
    overflow: 'hidden',
  },
  headerRow: {
    flexDirection: 'row',
    alignItems: 'center',
    borderBottomWidth: 1,
    borderBottomColor: 'rgba(127,127,127,0.28)',
    minHeight: 38,
    paddingHorizontal: 10,
  },
  headerCellBase: {
    fontSize: 12,
    fontWeight: '700',
    paddingVertical: 8,
  },
  headerCellRating: { width: 52 },
  headerCellDate: { width: 58 },
  headerCellAuthor: { width: 80 },
  headerCellArticle: { flex: 1 },
  placeholderBody: {
    minHeight: 112,
    padding: 12,
    justifyContent: 'center',
  },
  placeholderText: {
    fontSize: 14,
    lineHeight: 20,
    fontStyle: 'italic',
  },
});
