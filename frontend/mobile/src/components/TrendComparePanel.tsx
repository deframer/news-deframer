import React from 'react';
import { StyleSheet, Text } from 'react-native';
import { useTranslation } from 'react-i18next';

import { AppPalette } from '../theme';
import { Card } from './Card';

export const TrendComparePanel = ({ palette }: { palette: AppPalette }) => {
  const { t } = useTranslation();

  return (
    <Card palette={palette}>
      <Text style={[styles.title, { color: palette.text }]}>{t('mobile.trends_compare_title')}</Text>
      <Text style={[styles.body, { color: palette.secondaryText }]}>{t('mobile.trends_compare_body')}</Text>
    </Card>
  );
};

const styles = StyleSheet.create({
  title: { marginBottom: 8, fontSize: 20, fontWeight: '700' },
  body: { fontSize: 15, lineHeight: 22 },
});
