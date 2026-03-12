import React from 'react';
import { Pressable, ScrollView, StyleSheet, Text, View } from 'react-native';
import { useTranslation } from 'react-i18next';

import { Card } from '../components/Card';
import { SECTION_HEADER_GAP } from '../components/formStyles';
import { LoadingSpinner } from '../components/LoadingSpinner';
import { SelectField } from '../components/SelectField';
import { SegmentedControl } from '../components/SegmentedControl';
import { HostStatus, StatusBadge } from '../components/StatusBadge';
import { TextField } from '../components/TextField';
import { AppPalette } from '../theme';
import { Settings } from '../services/settingsService';

export const SettingsScreen = ({
  palette,
  settings,
  status,
  errorMessage,
  onSettingsChange,
  onTestConnection,
}: {
  palette: AppPalette;
  settings: Settings;
  status: HostStatus;
  errorMessage: string | null;
  onSettingsChange: (settings: Settings) => void;
  onTestConnection: () => void;
}) => {
  const { t } = useTranslation();
  const isLoading = status === 'loading';

  return (
    <ScrollView style={[styles.screen, { backgroundColor: palette.background }]} contentContainerStyle={styles.content}>
      <View style={styles.grid}>
        <View style={styles.column}>
          <Card palette={palette}>
            <View style={styles.cardHeader}>
              <Text style={[styles.sectionTitle, { color: palette.secondaryText }]}>{t('options.section_connection')}</Text>
              <StatusBadge
                status={status}
                palette={palette}
                labels={{
                  connected: t('options.status_connected'),
                  error: t('options.status_error'),
                  checking: t('options.status_loading'),
                }}
              />
            </View>
            {errorMessage ? (
              <View style={[styles.errorBox, { borderColor: palette.border, backgroundColor: palette.card }]}> 
                <Text style={[styles.errorTitle, { color: palette.text }]}>{t('options.status_error')}</Text>
                <Text style={[styles.errorBody, { color: palette.secondaryText }]}>{errorMessage}</Text>
              </View>
            ) : null}
            <TextField label={t('options.label_server_url')} value={settings.backendUrl} onChangeText={(backendUrl) => onSettingsChange({ ...settings, backendUrl })} palette={palette} />
            <TextField label={t('options.label_username')} optional={t('options.label_optional')} value={settings.username} onChangeText={(username) => onSettingsChange({ ...settings, username })} palette={palette} />
            <TextField label={t('options.label_password')} optional={t('options.label_optional')} value={settings.password} onChangeText={(password) => onSettingsChange({ ...settings, password })} palette={palette} secureTextEntry />
            <Pressable onPress={onTestConnection} disabled={isLoading} style={[styles.actionButton, { backgroundColor: palette.buttonBackground, borderColor: palette.buttonBorder }]}> 
              {isLoading ? (
                <LoadingSpinner palette={palette} label={t('options.btn_testing')} center />
              ) : (
                <Text style={[styles.actionButtonText, { color: palette.buttonText }]}>{t('options.btn_test_connection')}</Text>
              )}
            </Pressable>
          </Card>
        </View>

        <View style={styles.column}>
          <Card palette={palette}>
            <Text style={[styles.sectionTitle, styles.sectionGap, { color: palette.secondaryText }]}>{t('options.language_label')}</Text>
            <SelectField
              label={t('options.language_label')}
              value={settings.language}
              onValueChange={(language) => onSettingsChange({ ...settings, language })}
              items={[{ label: t('options.default'), value: 'default' }, { label: 'English', value: 'en' }, { label: 'Deutsch', value: 'de' }]}
              palette={palette}
              hideLabel
            />
          </Card>
          <Card palette={palette}>
            <Text style={[styles.sectionTitle, styles.sectionGap, { color: palette.secondaryText }]}>{t('options.theme_label')}</Text>
            <SegmentedControl
              palette={palette}
              value={settings.theme}
              onChange={(theme) => onSettingsChange({ ...settings, theme: theme as Settings['theme'] })}
              options={[
                { label: t('options.theme_light'), value: 'light' },
                { label: t('options.theme_dark'), value: 'dark' },
                { label: t('options.default'), value: 'system' },
              ]}
            />
          </Card>
        </View>
      </View>
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  screen: { flex: 1 },
  content: { padding: 24, gap: 16 },
  grid: { gap: 16 },
  column: { gap: 16 },
  cardHeader: { marginBottom: SECTION_HEADER_GAP, flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', gap: 12 },
  errorBox: {
    borderWidth: 1,
    borderRadius: 8,
    paddingHorizontal: 14,
    paddingVertical: 12,
    marginBottom: 12,
    gap: 6,
  },
  errorTitle: { fontSize: 16, fontWeight: '700' },
  errorBody: { fontSize: 15, lineHeight: 22 },
  sectionTitle: { fontSize: 14, letterSpacing: 1, textTransform: 'uppercase', fontWeight: '700' },
  sectionGap: { marginBottom: SECTION_HEADER_GAP },
  actionButton: { borderWidth: 1, borderRadius: 8, minHeight: 54, paddingVertical: 12, alignItems: 'center', justifyContent: 'center' },
  actionButtonText: { fontSize: 18, fontWeight: '600' },
});
