import React from 'react';
import { ScrollView, StyleSheet, Text, View } from 'react-native';
import { useTranslation } from 'react-i18next';

import { Card } from '../components/Card';
import { AppPalette } from '../theme';
import { DomainEntry } from '../services/newsDeframerClient';

export const SessionScreen = ({ palette, domain }: { palette: AppPalette; domain: DomainEntry }) => {
  const { t } = useTranslation();

  return (
    <ScrollView style={{ flex: 1, backgroundColor: palette.background }} contentContainerStyle={styles.content}>
      <View style={styles.topTabs}>
        {[t('trends.cloud'), t('trends.compare_view'), t('trends.search')].map((label, index) => (
          <View key={label} style={[styles.tab, index === 0 ? { borderBottomColor: palette.accent } : null]}>
            <Text style={[styles.tabText, { color: index === 0 ? palette.accent : palette.text }]}>{label}</Text>
          </View>
        ))}
      </View>

      <Card palette={palette}>
        <Text style={[styles.domainTitle, { color: palette.text }]}>{domain.domain}</Text>
        <Text style={[styles.domainMeta, { color: palette.secondaryText }]}>{t('mobile.session_title')}</Text>
        <Text style={[styles.helper, { color: palette.secondaryText }]}>
          {t('portal.trends')} - {domain.language.toUpperCase()}
        </Text>
      </Card>
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  content: { padding: 24, gap: 16 },
  topTabs: { flexDirection: 'row', borderBottomWidth: 1, marginBottom: 16 },
  tab: { flex: 1, alignItems: 'center', paddingVertical: 12, borderBottomWidth: 3, borderBottomColor: 'transparent' },
  tabText: { fontSize: 18, fontWeight: '600' },
  domainTitle: { marginBottom: 8, fontSize: 28, fontWeight: '700' },
  domainMeta: { marginBottom: 8, fontSize: 14, fontWeight: '600', textTransform: 'uppercase' },
  helper: { fontSize: 16 },
});
