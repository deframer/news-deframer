import './src/i18n';

import { useCallback, useEffect, useMemo, useState } from 'react';
import { ArrowLeft, Menu } from 'lucide-react-native';
import { Appearance, BackHandler, NativeModules, Platform, Pressable, StatusBar, StyleSheet, Text, View } from 'react-native';
import { useTranslation } from 'react-i18next';
import { SafeAreaProvider, SafeAreaView } from 'react-native-safe-area-context';

import { MenuDrawer } from './src/components/MenuDrawer';
import { LoadingSpinner } from './src/components/LoadingSpinner';
import { ArticleScreen } from './src/screens/ArticleScreen';
import { AboutScreen } from './src/screens/AboutScreen';
import { DashboardScreen } from './src/screens/DashboardScreen';
import { PortalScreen } from './src/screens/PortalScreen';
import { SettingsScreen } from './src/screens/SettingsScreen';
import { AnalyzedItem, DomainEntry, NewsDeframerClient } from './src/services/newsDeframerClient';
import i18n from './src/i18n';
import { DEFAULT_SETTINGS, settingsService, Settings } from './src/services/settingsService';
import { themeService } from './src/services/themeService';
import { HostStatus } from './src/components/StatusBadge';
import { logger } from './src/services/logger';

type Screen = 'dashboard' | 'settings' | 'about' | 'portal' | 'article';

const SUPPORTED_LANGUAGES = ['de', 'en'] as const;
type SupportedLanguage = (typeof SUPPORTED_LANGUAGES)[number];

const resolveSupportedLanguage = (rawLocale?: string): SupportedLanguage => {
  const normalized = rawLocale?.toLowerCase().split(/[-_]/)[0];
  return SUPPORTED_LANGUAGES.includes(normalized as SupportedLanguage) ? (normalized as SupportedLanguage) : 'en';
};

const getUrlHost = (url?: string): string => {
  if (!url) {
    return '';
  }

  try {
    return new URL(url).hostname;
  } catch {
    return url.replace(/^https?:\/\//, '').split('/')[0] || '';
  }
};

const getResolvedLanguage = (language: string): string => {
  if (language !== 'default') {
    return resolveSupportedLanguage(language);
  }

  const nativeLocale = (() => {
    if (Platform.OS === 'ios') {
      const settings = NativeModules.SettingsManager?.settings;
      const locale = settings?.AppleLocale;
      const languages = settings?.AppleLanguages;

      if (typeof locale === 'string' && locale.length > 0) {
        return locale;
      }

      if (Array.isArray(languages) && typeof languages[0] === 'string') {
        return languages[0];
      }
    }

    if (Platform.OS === 'android') {
      const locales = NativeModules.I18nManager?.localeIdentifier || NativeModules.I18nManager?.locale;
      if (typeof locales === 'string' && locales.length > 0) {
        return locales;
      }

      const platformLocales = NativeModules.PlatformConstants?.locales;
      if (Array.isArray(platformLocales) && typeof platformLocales[0] === 'string') {
        return platformLocales[0];
      }
    }

    return '';
  })();

  if (nativeLocale) {
    return resolveSupportedLanguage(nativeLocale);
  }

  const intlLocale = typeof Intl !== 'undefined' ? Intl.DateTimeFormat().resolvedOptions().locale : '';
  if (intlLocale) {
    return resolveSupportedLanguage(intlLocale);
  }

  if (typeof navigator !== 'undefined') {
    const rawLanguage =
      (typeof navigator.language === 'string' && navigator.language) ||
      (Array.isArray(navigator.languages) && typeof navigator.languages[0] === 'string' ? navigator.languages[0] : '') ||
      'en';
    return resolveSupportedLanguage(rawLanguage);
  }

  return 'en';
};

const withTimeout = async <T,>(promise: Promise<T>, timeoutMs: number): Promise<T> => {
  let timeoutHandle: ReturnType<typeof setTimeout> | undefined;
  const timeoutPromise = new Promise<never>((_, reject) => {
    timeoutHandle = setTimeout(() => reject(new Error(`Timeout after ${timeoutMs}ms`)), timeoutMs);
  });

  try {
    return await Promise.race([promise, timeoutPromise]);
  } finally {
    if (timeoutHandle) {
      clearTimeout(timeoutHandle);
    }
  }
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
  const [selectedArticle, setSelectedArticle] = useState<AnalyzedItem | null>(null);
  const [portalBackAction, setPortalBackAction] = useState<(() => void) | null>(null);
  const [settingsErrorMessage, setSettingsErrorMessage] = useState<string | null>(null);
  const palette = useMemo(() => themeService.getPalette(settings, colorScheme), [colorScheme, settings]);
  const visibleDomains = configured ? domains : [];
  const visibleDomainsLoading = configured ? domainsLoading : false;
  const activePortalBackAction = screen === 'portal' ? portalBackAction : null;

  useEffect(() => {
    if (status !== 'loading') {
      return;
    }

    const timeoutHandle = setTimeout(() => {
      setStatus((current) => (current === 'loading' ? 'error' : current));
    }, 5500);

    return () => clearTimeout(timeoutHandle);
  }, [status]);

  useEffect(() => {
    if (typeof document === 'undefined') {
      return;
    }

    const html = document.documentElement;
    const body = document.body;
    const previousHtmlHeight = html.style.height;
    const previousHtmlWidth = html.style.width;
    const previousHtmlOverflow = html.style.overflow;
    const previousBodyHeight = body.style.height;
    const previousBodyWidth = body.style.width;
    const previousBodyOverflow = body.style.overflow;
    const previousBodyMargin = body.style.margin;

    body.style.backgroundColor = palette.background;
    html.style.height = '100%';
    html.style.overflow = 'hidden';
    body.style.height = '100%';
    body.style.overflow = 'hidden';
    body.style.margin = '0';

    const root = document.getElementById('root');
    const previousRootHeight = root?.style.height ?? '';
    const previousRootOverflow = root?.style.overflow ?? '';
    const previousRootDisplay = root?.style.display ?? '';
    const previousRootFlexDirection = root?.style.flexDirection ?? '';

    if (root) {
      root.style.backgroundColor = palette.background;
      root.style.height = '100%';
      root.style.overflow = 'hidden';
      root.style.display = 'flex';
      root.style.flexDirection = 'column';
    }

    return () => {
      html.style.height = previousHtmlHeight;
      html.style.width = previousHtmlWidth;
      html.style.overflow = previousHtmlOverflow;
      body.style.height = previousBodyHeight;
      body.style.width = previousBodyWidth;
      body.style.overflow = previousBodyOverflow;
      body.style.margin = previousBodyMargin;
      if (root) {
        root.style.height = previousRootHeight;
        root.style.overflow = previousRootOverflow;
        root.style.display = previousRootDisplay;
        root.style.flexDirection = previousRootFlexDirection;
      }
    };
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
      try {
        const result = await withTimeout(settingsService.loadSettings(), 3000);
        if (!mounted) {
          return;
        }

        setSettings(result.settings);
        setConfigured(result.configured && settingsService.hasRequiredConfiguration(result.settings));
        setScreen(result.configured && settingsService.hasRequiredConfiguration(result.settings) ? 'dashboard' : 'settings');

        const lang = getResolvedLanguage(result.settings.language);
        await withTimeout(i18n.changeLanguage(lang), 2000);
      } catch {
        if (!mounted) {
          return;
        }

        setSettings(DEFAULT_SETTINGS);
        setConfigured(false);
        setScreen('settings');
      } finally {
        if (mounted) {
          setBooting(false);
        }
      }
    };

    boot();

    return () => {
      mounted = false;
    };
  }, []);

  useEffect(() => {
    if (booting) {
      return;
    }

    const handle = setTimeout(() => {
      settingsService.saveSettings(settings);
      i18n.changeLanguage(getResolvedLanguage(settings.language));
      setConfigured(settingsService.hasRequiredConfiguration(settings));
    }, 150);

    return () => clearTimeout(handle);
  }, [booting, settings]);

  useEffect(() => {
    if (booting || !configured) {
      return;
    }

    let mounted = true;
    const client = new NewsDeframerClient(settings);

    const loadDomains = async () => {
      setDomainsLoading(true);
      try {
        const loadedDomains = await withTimeout(client.getDomains(), 5000);
        if (mounted) {
          setDomains(loadedDomains);
        }
      } catch {
        if (mounted) {
          logger.error('Failed to load domains', { backendUrl: settings.backendUrl });
          setDomains([]);
        }
      } finally {
        if (mounted) {
          setDomainsLoading(false);
        }
      }
    };

    loadDomains();

    return () => {
      mounted = false;
    };
  }, [booting, configured, settings]);

  useEffect(() => {
    if (booting || screen !== 'settings') {
      return;
    }

    handleTestConnection();
  }, [booting, handleTestConnection, screen]);

  const handlePortalBackActionChange = useCallback((action: (() => void) | null) => {
    setPortalBackAction(() => action);
  }, []);

  const openScreen = (nextScreen: Screen) => {
    setDrawerOpen(false);
    setScreen(nextScreen);
  };

  const handleTestConnection = useCallback(async () => {
    setStatus('loading');
    setSettingsErrorMessage(null);

    const backendUrl = settings.backendUrl.trim();
    if (!backendUrl) {
      setStatus('error');
      setSettingsErrorMessage(t('mobile.connection_error_missing_url', 'Please enter a server URL before testing the connection.'));
      return;
    }

    try {
      const parsed = new URL(backendUrl);
      if (!['http:', 'https:'].includes(parsed.protocol)) {
        setStatus('error');
        setSettingsErrorMessage(t('mobile.connection_error_invalid_url', 'The server URL must start with http:// or https://.'));
        return;
      }
    } catch {
      setStatus('error');
      setSettingsErrorMessage(t('mobile.connection_error_invalid_url', 'The server URL must start with http:// or https://.'));
      return;
    }

    try {
      const client = new NewsDeframerClient(settings);
      await withTimeout(client.getDomains(), 5000);
      setStatus('success');
      setSettingsErrorMessage(null);
    } catch (error) {
      logger.error('Test connection failed', { backendUrl: settings.backendUrl });
      setStatus('error');
      setSettingsErrorMessage(error instanceof Error ? error.message : String(error));
    }
  }, [settings, t]);

  const renderScreen = () => {
    if (booting) {
      return (
        <View style={styles.loading}>
          <LoadingSpinner palette={palette} label={t('options.loading')} center />
        </View>
      );
    }

    if (screen === 'settings') {
      return (
        <SettingsScreen
          palette={palette}
          settings={settings}
          status={status}
          errorMessage={settingsErrorMessage}
          onSettingsChange={(nextSettings) => {
            setSettings(nextSettings);
            if (settingsErrorMessage) {
              setSettingsErrorMessage(null);
            }
          }}
          onTestConnection={handleTestConnection}
        />
      );
    }

    if (screen === 'about') {
      return <AboutScreen palette={palette} onClose={() => setScreen('dashboard')} />;
    }

    if (screen === 'portal' && selectedDomain) {
      return (
        <View style={styles.screenStack}>
          <View style={styles.screenLayer}>
            <PortalScreen
              palette={palette}
              domain={selectedDomain}
              settings={settings}
              onBackRequestChange={handlePortalBackActionChange}
              onOpenArticle={(item) => {
                setSelectedArticle(item);
              }}
            />
          </View>
          {selectedArticle ? (
            <View style={styles.screenOverlay}>
              <ArticleScreen palette={palette} item={selectedArticle} />
            </View>
          ) : null}
        </View>
      );
    }

    return (
      <DashboardScreen
        palette={palette}
        domains={visibleDomains}
        domainsLoading={visibleDomainsLoading}
        configured={configured}
        onOpenPortal={(domain) => {
          setSelectedDomain(domain);
          setScreen('portal');
        }}
      />
    );
  };

  const showBack = screen === 'settings' || screen === 'about' || screen === 'portal' || Boolean(selectedArticle);
  const showMenu = !showBack;
  const headerTitle = screen === 'settings'
    ? t('options.settings_title')
    : screen === 'about'
      ? t('mobile.about_title')
      : selectedArticle
        ? getUrlHost(selectedArticle?.url) || t('article.screen_title', 'Article')
        : screen === 'portal'
        ? selectedDomain?.domain || t('mobile.portal_title')
        : t('mobile.dashboard_title');
  const handleBack = useCallback(() => {
    if (selectedArticle || screen === 'article') {
      setSelectedArticle(null);
      setScreen('portal');
      return;
    }

    if (screen === 'portal' && activePortalBackAction) {
      activePortalBackAction();
      return;
    }

    setScreen('dashboard');
  }, [activePortalBackAction, screen, selectedArticle]);

  useEffect(() => {
    if (Platform.OS !== 'android' || !showBack) {
      return;
    }

    const subscription = BackHandler.addEventListener('hardwareBackPress', () => {
      handleBack();
      return true;
    });

    return () => subscription.remove();
  }, [handleBack, showBack]);

  return (
    <SafeAreaProvider>
      <StatusBar barStyle={palette.background === '#18191a' ? 'light-content' : 'dark-content'} backgroundColor={palette.background} />
      <SafeAreaView style={[styles.container, { backgroundColor: palette.background }]}> 
        <View style={[styles.appBar, { borderBottomColor: palette.border, backgroundColor: palette.background }]}> 
          <View style={styles.appBarLeft}>
            {showBack ? (
              <Pressable onPress={handleBack} style={styles.iconButton}>
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

        <View style={styles.screen}>
          {renderScreen()}
        </View>

        <MenuDrawer
          visible={drawerOpen}
          palette={palette}
          onClose={() => setDrawerOpen(false)}
          onNavigate={(nextScreen) => openScreen(nextScreen)}
          labels={{
            settings: t('mobile.menu_settings'),
            about: t('mobile.menu_about'),
          }}
        />
      </SafeAreaView>
    </SafeAreaProvider>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, overflow: 'hidden' },
  appBar: {
    position: 'sticky',
    top: 0,
    zIndex: 20,
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
  screen: { flex: 1, overflow: 'hidden' },
  loading: { flex: 1, alignItems: 'center', justifyContent: 'center' },
  screenStack: { flex: 1, position: 'relative' },
  screenLayer: { flex: 1 },
  screenOverlay: {
    ...StyleSheet.absoluteFillObject,
    zIndex: 5,
  },
});

export default App;
