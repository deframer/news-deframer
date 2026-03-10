import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

const withSuppressedI18nextPromo = <T>(run: () => T): T => {
  const originalInfo = console.info;

  console.info = (...args: unknown[]) => {
    if (typeof args[0] === 'string' && args[0].includes('i18next is made possible by our own product, Locize')) {
      return;
    }

    originalInfo(...args);
  };

  try {
    return run();
  } finally {
    console.info = originalInfo;
  }
};

const resources = {
  en: {
    translation: {
      portal: { title: 'News Deframer', trends: 'Trends', hide: 'Hide' },
      trends: {
        cloud: 'Tag Cloud', compare_view: 'Compare', lifecycle: 'Lifecycle', context: 'Context', articles: 'Articles', search: 'Search', no_data: 'No trending topics found for this period.', rank: 'Rank', trend: 'Trend', vol: 'Vol', time_ranges: { last_24h: '24h', last_7d: '7d', last_30d: '30d', last_90d: '90d', last_365d: '365d' },
      },
      article: { back: 'Back' },
      footer: { github_link: 'News Deframer on GitHub' },
      options: {
        settings_title: 'Settings', language_label: 'Language', theme_label: 'Theme', theme_light: 'Light', theme_dark: 'Dark', theme_system: 'System', default: 'Default', loading: 'Loading...', status_connected: 'Connected', status_error: 'Error', status_loading: 'Loading...', section_connection: 'Connection', label_server_url: 'Server URL', label_username: 'Username', label_password: 'Password', label_optional: '(Optional)', btn_test_connection: 'Test Connection', btn_testing: 'Testing...', label_search_engine: 'Search Engine', error_https_only: 'URL must start with https://',
      },
      mobile: {
        dashboard_title: 'Dashboard', about_title: 'About', about_body: 'News Deframer', menu_dashboard: 'Dashboard', menu_settings: 'Settings', menu_about: 'About', no_domains: 'No domains available.', pick_domain: 'Choose a domain to open a session.', missing_config: 'Add your server settings to continue.', session_title: 'Session', open_project: 'News Deframer on GitHub', ok: 'OK',
      },
    },
  },
  de: {
    translation: {
      portal: { title: 'News Deframer', trends: 'Trends', hide: 'Ausblenden' },
      trends: {
        cloud: 'Tag Cloud', compare_view: 'Vergleich', lifecycle: 'Verlauf', context: 'Kontext', articles: 'Artikel', search: 'Suche', no_data: 'Keine Trendthemen für diesen Zeitraum gefunden.', rank: 'Rang', trend: 'Trend', vol: 'Vol', time_ranges: { last_24h: '24h', last_7d: '7 T.', last_30d: '30 T.', last_90d: '90 T.', last_365d: '365 T.' },
      },
      article: { back: 'Zurück' },
      footer: { github_link: 'News Deframer auf GitHub' },
      options: {
        settings_title: 'Einstellungen', language_label: 'Sprache', theme_label: 'Design', theme_light: 'Hell', theme_dark: 'Dunkel', theme_system: 'System', default: 'Standard', loading: 'Laden...', status_connected: 'Verbunden', status_error: 'Fehler', status_loading: 'Laden...', section_connection: 'Verbindung', label_server_url: 'Server-URL', label_username: 'Benutzername', label_password: 'Passwort', label_optional: '(Optional)', btn_test_connection: 'Verbindung testen', btn_testing: 'Teste...', label_search_engine: 'Suchmaschine', error_https_only: 'URL muss mit https:// beginnen',
      },
      mobile: {
        dashboard_title: 'Dashboard', about_title: 'Über', about_body: 'News Deframer', menu_dashboard: 'Dashboard', menu_settings: 'Einstellungen', menu_about: 'Über', no_domains: 'Keine Domains verfügbar.', pick_domain: 'Wählen Sie eine Domain, um eine Sitzung zu öffnen.', missing_config: 'Fügen Sie Ihre Server-Einstellungen hinzu, um fortzufahren.', session_title: 'Sitzung', open_project: 'News Deframer auf GitHub', ok: 'OK',
      },
    },
  },
} as const;

if (!i18n.isInitialized) {
  withSuppressedI18nextPromo(() => {
    i18n.use(initReactI18next).init({
      resources,
      fallbackLng: 'en',
      interpolation: { escapeValue: false },
    });
  });
}

export default i18n;
