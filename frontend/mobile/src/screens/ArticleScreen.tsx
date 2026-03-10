import React from 'react';
import { ScrollView, StyleSheet, Text } from 'react-native';
import { useTranslation } from 'react-i18next';

import { Card } from '../components/Card';
import { AppPalette } from '../theme';

export const ArticleScreen = ({ palette, url }: { palette: AppPalette; url: string }) => {
  const { t } = useTranslation();

  return (
    <ScrollView style={{ flex: 1, backgroundColor: palette.background }} contentContainerStyle={styles.content}>
      <Card palette={palette}>
        <Text style={[styles.title, { color: palette.text }]}>{t('article.placeholder_title', 'Article screen placeholder')}</Text>
        <Text style={[styles.body, { color: palette.secondaryText }]}>{t('article.selected_url', 'Selected URL')}</Text>
        <Text style={[styles.url, { color: palette.text }]}>{url}</Text>
        <Text style={[styles.body, { color: palette.secondaryText }]}>{t('article.placeholder_body', 'The dedicated mobile article screen will be implemented next.')}</Text>
      </Card>
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  content: { padding: 24 },
  title: { marginBottom: 12, fontSize: 24, fontWeight: '700' },
  body: { fontSize: 15, lineHeight: 22 },
  url: { marginTop: 8, marginBottom: 18, fontSize: 16, fontWeight: '600' },
});
