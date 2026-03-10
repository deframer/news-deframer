import './src/i18n';

import { useEffect, useMemo, useState } from 'react';
import { ArrowLeft, Menu } from 'lucide-react-native';
import { Appearance, Pressable, StatusBar, StyleSheet, Text, View } from 'react-native';
import { useTranslation } from 'react-i18next';
import { SafeAreaProvider, SafeAreaView } from 'react-native-safe-area-context';

import { MenuDrawer } from './src/components/MenuDrawer';
import { LoadingSpinner } from './src/components/LoadingSpinner';
import { AboutScreen } from './src/screens/AboutScreen';
import { DashboardScreen } from './src/screens/DashboardScreen';
import { SessionScreen } from './src/screens/SessionScreen';
import { SettingsScreen } from './src/screens/SettingsScreen';
import { DomainEntry, NewsDeframerClient } from './src/services/newsDeframerClient';
import i18n from './src/i18n';
import { DEFAULT_SETTINGS, settingsService, Settings } from './src/services/settingsService';
import { themeService } from './src/services/themeService';
import { HostStatus } from './src/components/StatusBadge';

type Screen = 'dashboard' | 'settings' | 'about' | 'session';

const getResolvedLanguage = (language: string): string => {
  if (language !== 'default') {
    return language;
  }

  if (typeof navigator !== 'undefined') {
    const detected = navigator.language.split('-')[0];
    return ['de', 'en'].includes(detected) ? detected : 'en';
  }

  return 'en';
};

function App() {
  const { t } = useTranslation();
  const colorScheme = Appearance.getColorScheme();
  const [booting, setBooting] = useState(true);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [configured, setConfigured] = useState(false);
  const [settings, setSettings] = useState<Settings>(DEFAULT_SETTINGS);
  const [screen, setScreen] = useState<Screen>('settings');
  const [status, setStatus] = useState<HostStatus>('idle');
  const [domains, setDomains] = useState<DomainEntry[]>([]);
  const [domainsLoading, setDomainsLoading] = useState(false);
  const [selectedDomain, setSelectedDomain] = useState<DomainEntry | null>(null);
  const palette = useMemo(() => themeService.getPalette(settings, colorScheme), [colorScheme, settings]);

  useEffect(() => {
    if (typeof document === 'undefined') {
      return;
    }

    document.body.style.backgroundColor = palette.background;
    const root = document.getElementById('root');
    if (root) {
      root.style.backgroundColor = palette.background;
    }
  }, [palette.background]);

  useEffect(() => {
    if (typeof document === 'undefined') {
      return;
    }

    document.documentElement.lang = getResolvedLanguage(settings.language);
  }, [settings.language]);

  useEffect(() => {
    let mounted = true;

    const boot = async () => {
      const result = await settingsService.loadSettings();
      if (!mounted) {
        return;
      }
      setSettings(result.settings);
      setConfigured(result.configured && settingsService.hasRequiredConfiguration(result.settings));
      setScreen(result.configured && settingsService.hasRequiredConfiguration(result.settings) ? 'dashboard' : 'settings');
      const lang = getResolvedLanguage(result.settings.language);
      await i18n.changeLanguage(lang);
      setBooting(false);
    };

    void boot();

    return () => {
      mounted = false;
    };
  }, []);

  useEffect(() => {
    if (booting) {
      return;
    }

    const handle = setTimeout(() => {
      void settingsService.saveSettings(settings);
      void i18n.changeLanguage(getResolvedLanguage(settings.language));
      setConfigured(settingsService.hasRequiredConfiguration(settings));
    }, 150);

    return () => clearTimeout(handle);
  }, [booting, settings]);

  useEffect(() => {
    if (booting || !configured) {
      setDomains([]);
      setDomainsLoading(false);
      return;
    }

    let mounted = true;
    const client = new NewsDeframerClient(settings);

    const loadDomains = async () => {
      setDomainsLoading(true);
      try {
        // used for debuging the loading spinner
        // await new Promise((resolve) => setTimeout(resolve, 10000));
        const loadedDomains = await client.getDomains();
        if (mounted) {
          setDomains(loadedDomains);
        }
      } catch {
        if (mounted) {
          setDomains([]);
        }
      } finally {
        if (mounted) {
          setDomainsLoading(false);
        }
      }
    };

    void loadDomains();

    return () => {
      mounted = false;
    };
  }, [booting, configured, settings]);

  useEffect(() => {
    if (booting || screen !== 'settings') {
      return;
    }

    void handleTestConnection();
  }, [booting, screen]);

  const openScreen = (nextScreen: Screen) => {
    setDrawerOpen(false);
    setScreen(nextScreen);
  };

  const handleTestConnection = async () => {
    setStatus('loading');
    try {
      const client = new NewsDeframerClient(settings);
      await client.getDomains();
      setStatus('success');
    } catch {
      setStatus('error');
    }
  };

  const renderScreen = () => {
    if (booting) {
      return (
        <View style={styles.loading}>
          <LoadingSpinner palette={palette} label={t('options.loading')} center />
        </View>
      );
    }

    if (screen === 'settings') {
      return <SettingsScreen palette={palette} settings={settings} status={status} onSettingsChange={setSettings} onTestConnection={handleTestConnection} />;
    }

    if (screen === 'about') {
      return <AboutScreen palette={palette} onClose={() => setScreen('dashboard')} />;
    }

    if (screen === 'session' && selectedDomain) {
      return <SessionScreen palette={palette} domain={selectedDomain} />;
    }

    return (
      <DashboardScreen
        palette={palette}
        domains={domains}
        domainsLoading={domainsLoading}
        configured={configured}
        onOpenSession={(domain) => {
          setSelectedDomain(domain);
          setScreen('session');
        }}
      />
    );
  };

  const showBack = screen === 'settings' || screen === 'about';
  const showMenu = !showBack;
  const headerTitle = screen === 'settings' ? t('options.settings_title') : screen === 'about' ? t('mobile.about_title') : screen === 'session' ? selectedDomain?.domain || t('mobile.session_title') : t('mobile.dashboard_title');

  return (
    <SafeAreaProvider>
      <StatusBar barStyle={palette.background === '#18191a' ? 'light-content' : 'dark-content'} backgroundColor={palette.background} />
      <SafeAreaView style={[styles.container, { backgroundColor: palette.background }]}> 
        <View style={[styles.appBar, { borderBottomColor: palette.border, backgroundColor: palette.background }]}> 
          <View style={styles.appBarLeft}>
            {showBack ? (
              <Pressable onPress={() => setScreen('dashboard')} style={styles.iconButton}>
                <ArrowLeft color={palette.text} size={22} strokeWidth={2.25} />
              </Pressable>
            ) : null}
            {showMenu ? (
              <Pressable onPress={() => setDrawerOpen(true)} style={styles.iconButton}>
                <Menu color={palette.text} size={22} strokeWidth={2.25} />
              </Pressable>
            ) : null}
          </View>
          <Text style={[styles.appBarTitle, { color: palette.text }]}>{headerTitle}</Text>
          <View style={styles.appBarSpacer} />
        </View>

        {renderScreen()}

        <MenuDrawer
          visible={drawerOpen}
          palette={palette}
          onClose={() => setDrawerOpen(false)}
          onNavigate={(nextScreen) => openScreen(nextScreen)}
          labels={{
            dashboard: t('mobile.menu_dashboard'),
            settings: t('mobile.menu_settings'),
            about: t('mobile.menu_about'),
          }}
        />
      </SafeAreaView>
    </SafeAreaProvider>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1 },
  appBar: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: 12,
    paddingVertical: 10,
    borderBottomWidth: 1,
  },
  appBarLeft: { flexDirection: 'row', alignItems: 'center', minWidth: 48 },
  appBarSpacer: { minWidth: 48 },
  iconButton: { paddingHorizontal: 8, paddingVertical: 6 },
  appBarTitle: { fontSize: 18, fontWeight: '700' },
  loading: { flex: 1, alignItems: 'center', justifyContent: 'center' },
});

export default App;
