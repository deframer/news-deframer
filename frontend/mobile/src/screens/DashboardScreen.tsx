import React from 'react';
import { Pressable, ScrollView, StyleSheet, Text, View } from 'react-native';
import { useTranslation } from 'react-i18next';

import { Card } from '../components/Card';
import { LoadingSpinner } from '../components/LoadingSpinner';
import { Radio } from '../components/icons';
import { AppPalette } from '../theme';
import { DomainEntry } from '../services/newsDeframerClient';

const formatDomainLocale = (country: string, language: string) => `${country.toUpperCase()} | ${language}`;

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
              accessibilityRole="button"
              accessibilityLabel={domain.tags?.includes('public_service_media') ? `${domain.domain}, ${formatDomainLocale(domain.country, domain.language)}, ${t('article.public_service_media', 'Public service media')}` : undefined}
              style={[styles.domainRow, index > 0 ? [styles.domainRowDivider, { borderTopColor: palette.border }] : null]}
            >
              <Text style={[styles.domainName, { color: palette.text }]}>{domain.domain}</Text>
              <View style={styles.domainMetaWrap}>
                {domain.tags?.includes('public_service_media') ? <Radio color={palette.secondaryText} size={14} strokeWidth={2.2} /> : null}
                <Text style={[styles.domainMeta, { color: palette.secondaryText }]}>{formatDomainLocale(domain.country, domain.language)}</Text>
              </View>
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
  domainMetaWrap: { flexDirection: 'row', alignItems: 'center', gap: 4, flexShrink: 0 },
  domainMeta: { fontSize: 14, fontWeight: '500' },
});
