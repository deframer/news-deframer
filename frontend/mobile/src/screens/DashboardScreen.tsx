import React from 'react';
import { Pressable, ScrollView, StyleSheet, Text, View } from 'react-native';
import { useTranslation } from 'react-i18next';

import { Card } from '../components/Card';
import { LoadingSpinner } from '../components/LoadingSpinner';
import { AppPalette } from '../theme';
import { DomainEntry } from '../services/newsDeframerClient';

export const DashboardScreen = ({
  palette,
  domains,
  domainsLoading,
  configured,
  onOpenSession,
}: {
  palette: AppPalette;
  domains: DomainEntry[];
  domainsLoading: boolean;
  configured: boolean;
  onOpenSession: (domain: DomainEntry) => void;
}) => {
  const { t } = useTranslation();

  return (
    <ScrollView style={{ flex: 1, backgroundColor: palette.background }} contentContainerStyle={styles.content}>
      <View style={styles.headerBlock}>
        <Text style={[styles.eyebrow, { color: palette.secondaryText }]}>NEWS DEFRAMER</Text>
        <Text style={[styles.title, { color: palette.text }]}>{t('mobile.dashboard_title')}</Text>
        {!configured ? <Text style={[styles.subtitle, { color: palette.secondaryText }]}>{t('mobile.missing_config')}</Text> : null}
      </View>
      <Card palette={palette}>
        {domainsLoading ? (
          <LoadingSpinner palette={palette} label={t('options.status_loading')} />
        ) : domains.length === 0 ? (
          <Text style={[styles.empty, { color: palette.secondaryText }]}>{t('mobile.no_domains')}</Text>
        ) : (
          domains.map((domain, index) => (
            <Pressable
              key={domain.domain}
              onPress={() => onOpenSession(domain)}
              style={[styles.domainRow, index > 0 ? { borderTopWidth: 1, borderTopColor: palette.border } : null]}
            >
              <Text style={[styles.domainName, { color: palette.text }]}>{domain.domain}</Text>
              <Text style={[styles.domainMeta, { color: palette.secondaryText }]}>{domain.language.toUpperCase()}</Text>
            </Pressable>
          ))
        )}
      </Card>
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  content: { padding: 24, gap: 16 },
  headerBlock: { marginBottom: 4 },
  eyebrow: { marginBottom: 6, fontSize: 12, letterSpacing: 2, textTransform: 'uppercase' },
  title: { marginBottom: 8, fontSize: 30, fontWeight: '700' },
  subtitle: { fontSize: 16, lineHeight: 22 },
  empty: { fontSize: 16 },
  domainRow: { paddingVertical: 14, flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center' },
  domainName: { fontSize: 18, fontWeight: '600' },
  domainMeta: { fontSize: 14, fontWeight: '500' },
});
