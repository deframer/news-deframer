import React from 'react';
import { Pressable, ScrollView, StyleSheet, Text } from 'react-native';
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
  onOpenPortal,
}: {
  palette: AppPalette;
  domains: DomainEntry[];
  domainsLoading: boolean;
  configured: boolean;
  onOpenPortal: (domain: DomainEntry) => void;
}) => {
  const { t } = useTranslation();

  return (
    <ScrollView style={[styles.screen, { backgroundColor: palette.background }]} contentContainerStyle={styles.content}>
      {!configured ? <Text style={[styles.subtitle, { color: palette.secondaryText }]}>{t('mobile.missing_config')}</Text> : null}
      <Card palette={palette}>
        {domainsLoading ? (
          <LoadingSpinner palette={palette} label={t('options.status_loading')} />
        ) : domains.length === 0 ? (
          <Text style={[styles.empty, { color: palette.secondaryText }]}>{t('mobile.no_domains')}</Text>
        ) : (
          domains.map((domain, index) => (
            <Pressable
              key={domain.domain}
              onPress={() => onOpenPortal(domain)}
              style={[styles.domainRow, index > 0 ? [styles.domainRowDivider, { borderTopColor: palette.border }] : null]}
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
  screen: { flex: 1 },
  content: { padding: 24, gap: 16 },
  subtitle: { fontSize: 16, lineHeight: 22 },
  empty: { fontSize: 16 },
  domainRow: { paddingVertical: 14, flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center' },
  domainRowDivider: { borderTopWidth: 1 },
  domainName: { fontSize: 18, fontWeight: '600' },
  domainMeta: { fontSize: 14, fontWeight: '500' },
});
