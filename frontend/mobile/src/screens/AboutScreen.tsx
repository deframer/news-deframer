import React from 'react';
import { Github } from 'lucide-react-native';
import { Linking, Pressable, ScrollView, StyleSheet, Text, View } from 'react-native';
import { useTranslation } from 'react-i18next';

import { APP_VERSION } from '../appInfo';
import { AppLogo } from '../components/AppLogo';
import { Card } from '../components/Card';
import { AppPalette } from '../theme';

export const AboutScreen = ({ palette, onClose }: { palette: AppPalette; onClose: () => void }) => {
  const { t } = useTranslation();

  return (
    <ScrollView style={{ flex: 1, backgroundColor: palette.background }} contentContainerStyle={styles.content}>
      <Card palette={palette}>
        <View style={styles.logoWrap}>
          <AppLogo size={72} />
        </View>
        <Text style={[styles.title, { color: palette.text }]}>News Deframer</Text>
        <Text style={[styles.subtitle, { color: palette.secondaryText }]}>Mobile</Text>
        <Text style={[styles.version, { color: palette.secondaryText }]}>{APP_VERSION}</Text>
        <Pressable onPress={() => Linking.openURL('https://deframer.github.io/')} style={styles.linkRow}> 
          <Github color={palette.accent} size={18} strokeWidth={2.1} />
          <Text style={[styles.linkText, { color: palette.accent }]}>{t('mobile.open_project')}</Text>
        </Pressable>
        <Pressable onPress={onClose} style={[styles.okButton, { backgroundColor: palette.accent }]}> 
          <Text style={[styles.okText, { color: palette.accentText }]}>{t('mobile.ok')}</Text>
        </Pressable>
      </Card>
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  content: { flexGrow: 1, padding: 24, alignItems: 'center', justifyContent: 'center' },
  logoWrap: { alignItems: 'center', marginBottom: 12 },
  title: { marginBottom: 4, fontSize: 32, fontWeight: '700', textAlign: 'center' },
  subtitle: { marginBottom: 2, fontSize: 18, fontWeight: '600', textAlign: 'center' },
  version: { marginBottom: 18, fontSize: 16, textAlign: 'center' },
  linkRow: { flexDirection: 'row', alignItems: 'center', justifyContent: 'center', gap: 10, marginBottom: 20 },
  linkText: { fontSize: 18, fontWeight: '600', textDecorationLine: 'underline' },
  okButton: { borderRadius: 8, paddingHorizontal: 24, paddingVertical: 12, alignItems: 'center', alignSelf: 'center', minWidth: 120 },
  okText: { fontSize: 18, fontWeight: '700' },
});
