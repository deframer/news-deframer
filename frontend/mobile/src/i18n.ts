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
      article: { back: 'Back', no_title: 'No Title', no_description: 'No Description', screen_title: 'Article', selected_url: 'Selected URL', placeholder_title: 'Article screen placeholder', placeholder_body: 'The dedicated mobile article screen will be implemented next.' },
      metadata: { just_now: 'a moment ago' },
      rating: { no_reason: 'No reason provided', reason_title: 'Rating reason', info_aria_label: 'Show rating reason' },
      footer: { github_link: 'News Deframer on GitHub' },
      options: {
        settings_title: 'Settings', language_label: 'Language', theme_label: 'Theme', theme_light: 'Light', theme_dark: 'Dark', theme_system: 'System', default: 'Default', loading: 'Loading...', status_connected: 'Connected', status_error: 'Error', status_loading: 'Loading...', section_connection: 'Connection', label_server_url: 'Server URL', label_username: 'Username', label_password: 'Password', label_optional: '(Optional)', btn_test_connection: 'Test Connection', btn_testing: 'Testing...', label_search_engine: 'Search Engine', error_https_only: 'URL must start with https://',
      },
      mobile: {
        dashboard_title: 'News Deframer', about_title: 'About', about_body: 'News Deframer', menu_dashboard: 'Dashboard', menu_settings: 'Settings', menu_about: 'About', no_domains: 'No domains available.', pick_domain: 'Choose a domain to open a portal.', missing_config: 'Add your server settings to continue.', portal_title: 'Portal', trend_mining: 'Trend Mining', trend_mining_placeholder: 'Trend mining will land here next. For now, use the Articles tab to browse the current portal feed.', portal_empty_title: 'No articles yet', portal_empty_body: 'This domain is configured, but the backend did not return any portal articles.', portal_load_error_title: 'Could not load articles', portal_load_error: 'The mobile app could not load this portal right now. Please try again.', retry: 'Retry', open_project: 'News Deframer on GitHub', ok: 'OK',
      },
    },
  },
  de: {
    translation: {
      portal: { title: 'News Deframer', trends: 'Trends', hide: 'Ausblenden' },
      trends: {
        cloud: 'Tag Cloud', compare_view: 'Vergleich', lifecycle: 'Verlauf', context: 'Kontext', articles: 'Artikel', search: 'Suche', no_data: 'Keine Trendthemen für diesen Zeitraum gefunden.', rank: 'Rang', trend: 'Trend', vol: 'Vol', time_ranges: { last_24h: '24h', last_7d: '7 T.', last_30d: '30 T.', last_90d: '90 T.', last_365d: '365 T.' },
      },
      article: { back: 'Zurück', no_title: 'Kein Titel', no_description: 'Keine Beschreibung', screen_title: 'Artikel', selected_url: 'Ausgewählte URL', placeholder_title: 'Artikelansicht Platzhalter', placeholder_body: 'Die eigene mobile Artikelansicht wird als Nächstes umgesetzt.' },
      metadata: { just_now: 'gerade eben' },
      rating: { no_reason: 'Keine Begründung vorhanden', reason_title: 'Bewertung Begründung', info_aria_label: 'Bewertungsbegründung anzeigen' },
      footer: { github_link: 'News Deframer auf GitHub' },
      options: {
        settings_title: 'Einstellungen', language_label: 'Sprache', theme_label: 'Design', theme_light: 'Hell', theme_dark: 'Dunkel', theme_system: 'System', default: 'Standard', loading: 'Laden...', status_connected: 'Verbunden', status_error: 'Fehler', status_loading: 'Laden...', section_connection: 'Verbindung', label_server_url: 'Server-URL', label_username: 'Benutzername', label_password: 'Passwort', label_optional: '(Optional)', btn_test_connection: 'Verbindung testen', btn_testing: 'Teste...', label_search_engine: 'Suchmaschine', error_https_only: 'URL muss mit https:// beginnen',
      },
      mobile: {
        dashboard_title: 'News Deframer', about_title: 'Über', about_body: 'News Deframer', menu_dashboard: 'Dashboard', menu_settings: 'Einstellungen', menu_about: 'Über', no_domains: 'Keine Domains verfügbar.', pick_domain: 'Wählen Sie eine Domain, um ein Portal zu öffnen.', missing_config: 'Fügen Sie Ihre Server-Einstellungen hinzu, um fortzufahren.', portal_title: 'Portal', trend_mining: 'Trend Mining', trend_mining_placeholder: 'Trend Mining erscheint hier als Nächstes. Verwenden Sie vorerst den Tab Artikel, um den aktuellen Portal-Feed zu durchsuchen.', portal_empty_title: 'Noch keine Artikel', portal_empty_body: 'Diese Domain ist konfiguriert, aber das Backend hat keine Portal-Artikel zurückgegeben.', portal_load_error_title: 'Artikel konnten nicht geladen werden', portal_load_error: 'Die Mobile-App konnte dieses Portal gerade nicht laden. Bitte versuchen Sie es erneut.', retry: 'Erneut laden', open_project: 'News Deframer auf GitHub', ok: 'OK',
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
