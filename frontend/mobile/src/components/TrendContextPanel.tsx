import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { useTranslation } from 'react-i18next';

import { AppPalette } from '../theme';

export const TrendContextPanel = ({ palette, term }: { palette: AppPalette; term: string }) => {
  const { t } = useTranslation();

  return (
    <View style={styles.container}> 
      <Text style={[styles.title, { color: palette.text }]}>{t('trends.context')}</Text>
      <Text style={[styles.body, { color: palette.secondaryText }]}>{t('mobile.trends_context_placeholder', { term })}</Text>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginTop: 14,
  },
  title: { fontSize: 15, fontWeight: '700' },
  body: { marginTop: 6, fontSize: 14, lineHeight: 20 },
});
