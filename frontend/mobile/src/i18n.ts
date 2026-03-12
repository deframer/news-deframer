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
        cloud: 'Tag Cloud', compare_view: 'Compare', lifecycle: 'Lifecycle', context: 'Context', articles: 'Articles', search: 'Search', search_placeholder: 'Enter term to analyze...', analyze: 'Search', no_data: 'No trending topics found for this period.', rank: 'Rank', trend: 'Trend', vol: 'Vol', freq: 'Freq', vel: 'Vel', lifecycle_no_data: 'No lifecycle data available for this topic.', context_no_data: 'No context data available for this topic.', context_header: 'How is "{{topic}}" being described?', frequency_label: 'Frequency', verb_label: 'Verb', search_aria_label: '{{date}}: Frequency {{frequency}}, Velocity {{velocity}}', rating_caption: 'Rating', date_caption: 'Date', author_caption: 'Author', article_caption: 'Article', compare: { shared: 'Shared Trends' }, time_ranges: { last_24h: '24h', last_7d: '7d', last_30d: '30d', last_90d: '90d', last_365d: '365d' },
      },
      article: {
        back: 'Back',
        no_title: 'No Title',
        no_description: 'No Description',
        screen_title: 'Article',
        selected_url: 'Selected URL',
        placeholder_title: 'Article screen placeholder',
        placeholder_body: 'The dedicated mobile article screen will be implemented next.',
        original_section: 'Original',
        btn_original_title: 'Original title',
        btn_details: 'Details',
        btn_open_article: 'Open article',
      },
      metrics: {
        framing: 'Framing',
        clickbait: 'Clickbait',
        persuasive: 'Persuasive',
        hyper_stimulus: 'Hyper stimulus',
        speculative: 'Speculative',
        overall_rating: 'Overall rating',
      },
      metadata: { just_now: 'a moment ago' },
      rating: { no_reason: 'No reason provided', reason_title: 'Rating reason', info_aria_label: 'Show rating reason' },
      footer: { github_link: 'News Deframer on GitHub' },
      options: {
        settings_title: 'Settings', language_label: 'Language', theme_label: 'Theme', theme_light: 'Light', theme_dark: 'Dark', theme_system: 'System', default: 'Default', loading: 'Loading...', status_connected: 'Connected', status_error: 'Error', status_loading: 'Loading...', section_connection: 'Connection', label_server_url: 'Server URL', label_username: 'Username', label_password: 'Password', label_optional: '(Optional)', btn_test_connection: 'Test Connection', btn_testing: 'Testing...',
      },
      mobile: {
        dashboard_title: 'News Deframer', about_title: 'About', about_body: 'News Deframer', menu_dashboard: 'Dashboard', menu_settings: 'Settings', menu_about: 'About', no_domains: 'No domains available.', pick_domain: 'Choose a domain to open a portal.', missing_config: 'Add your server settings to continue.', portal_title: 'Portal', trends: 'Trends', trend_mining: 'Trend Mining', trend_mining_placeholder: 'Trend mining will land here next. For now, use the Articles tab to browse the current portal feed.', trends_tag_cloud_title: 'Tag Cloud', trends_tag_cloud_body: 'Tag cloud panel scaffold. Select a term to open the Trend Details area (Lifecycle, Context, Articles).', trends_selected_range: 'Selected range: {{range}}', trend_details_title: 'Trend Details', trends_lifecycle_placeholder: 'Lifecycle chart scaffold for "{{term}}". Data fetching and chart rendering will be implemented next.', trends_context_placeholder: 'Context panel scaffold for "{{term}}". Context analysis data will be implemented next.', trends_article_list_placeholder: 'Article list scaffold for "{{term}}". Browser-like table behavior/pagination will be implemented next.', trends_article_list_placeholder_date: 'Article list scaffold for "{{term}}" on {{date}}. Browser-like table behavior/pagination will be implemented next.', trends_chart_mode_compact: 'Switch lifecycle chart to compact mode', trends_chart_mode_wide: 'Switch lifecycle chart to wide mode', trends_compare_title: 'Compare', trends_compare_body: 'Compare view placeholder. Domain comparison will be added in the next step.', compare_with_label: 'Compare with', connection_error_missing_url: 'Please enter a server URL before testing the connection.', connection_error_invalid_url: 'The server URL must start with http:// or https://.', trends_search_title: 'Search', trends_search_body: 'Search placeholder. Trend search and drill-down will be added next.', selected: 'Selected', not_selected: 'Not selected', selected_trend_term: 'Selected term', change_term: 'Change term', portal_empty_title: 'No articles yet', portal_empty_body: 'This domain is configured, but the backend did not return any portal articles.', portal_load_error_title: 'Could not load articles', portal_load_error: 'The mobile app could not load this portal right now. Please try again.', retry: 'Retry', open_project: 'News Deframer on GitHub', ok: 'OK',
      },
    },
  },
  de: {
    translation: {
      portal: { title: 'News Deframer', trends: 'Trends', hide: 'Ausblenden' },
      trends: {
        cloud: 'Tag Cloud', compare_view: 'Vergleich', lifecycle: 'Verlauf', context: 'Kontext', articles: 'Artikel', search: 'Suche', search_placeholder: 'Begriff eingeben...', analyze: 'Suche', no_data: 'Keine Trendthemen für diesen Zeitraum gefunden.', rank: 'Rang', trend: 'Trend', vol: 'Vol', freq: 'Hfg', vel: 'Vel', lifecycle_no_data: 'Keine Verlaufsdaten für dieses Thema verfügbar.', context_no_data: 'Keine Kontextdaten für dieses Thema verfügbar.', context_header: 'Wie wird "{{topic}}" beschrieben?', frequency_label: 'Häufigkeit', verb_label: 'Verb', search_aria_label: '{{date}}: Häufigkeit {{frequency}}, Geschwindigkeit {{velocity}}', rating_caption: 'Bewertung', date_caption: 'Datum', author_caption: 'Autor', article_caption: 'Artikel', compare: { shared: 'Gemeinsame Trends' }, time_ranges: { last_24h: '24h', last_7d: '7 T.', last_30d: '30 T.', last_90d: '90 T.', last_365d: '365 T.' },
      },
      article: {
        back: 'Zurück',
        no_title: 'Kein Titel',
        no_description: 'Keine Beschreibung',
        screen_title: 'Artikel',
        selected_url: 'Ausgewählte URL',
        placeholder_title: 'Artikelansicht Platzhalter',
        placeholder_body: 'Die eigene mobile Artikelansicht wird als Nächstes umgesetzt.',
        original_section: 'Original',
        btn_original_title: 'Originaltitel',
        btn_details: 'Details',
        btn_open_article: 'Artikel öffnen',
      },
      metrics: {
        framing: 'Framing',
        clickbait: 'Clickbait',
        persuasive: 'Beeinflussung',
        hyper_stimulus: 'Reizüberflutung',
        speculative: 'Spekulativ',
        overall_rating: 'Gesamtbewertung',
      },
      metadata: { just_now: 'gerade eben' },
      rating: { no_reason: 'Keine Begründung vorhanden', reason_title: 'Bewertung Begründung', info_aria_label: 'Bewertungsbegründung anzeigen' },
      footer: { github_link: 'News Deframer auf GitHub' },
      options: {
        settings_title: 'Einstellungen', language_label: 'Sprache', theme_label: 'Design', theme_light: 'Hell', theme_dark: 'Dunkel', theme_system: 'System', default: 'Standard', loading: 'Laden...', status_connected: 'Verbunden', status_error: 'Fehler', status_loading: 'Laden...', section_connection: 'Verbindung', label_server_url: 'Server-URL', label_username: 'Benutzername', label_password: 'Passwort', label_optional: '(Optional)', btn_test_connection: 'Verbindung testen', btn_testing: 'Teste...',
      },
      mobile: {
        dashboard_title: 'News Deframer', about_title: 'Über', about_body: 'News Deframer', menu_dashboard: 'Dashboard', menu_settings: 'Einstellungen', menu_about: 'Über', no_domains: 'Keine Domains verfügbar.', pick_domain: 'Wählen Sie eine Domain, um ein Portal zu öffnen.', missing_config: 'Fügen Sie Ihre Server-Einstellungen hinzu, um fortzufahren.', portal_title: 'Portal', trends: 'Trends', trend_mining: 'Trend Mining', trend_mining_placeholder: 'Trend Mining erscheint hier als Nächstes. Verwenden Sie vorerst den Tab Artikel, um den aktuellen Portal-Feed zu durchsuchen.', trends_tag_cloud_title: 'Tag Cloud', trends_tag_cloud_body: 'Tag-Cloud-Panel als Strukturstufe. Wählen Sie einen Begriff, um den Bereich Trend-Details (Verlauf, Kontext, Artikel) zu öffnen.', trends_selected_range: 'Ausgewählter Zeitraum: {{range}}', trend_details_title: 'Trend-Details', trends_lifecycle_placeholder: 'Verlaufsdiagramm als Struktur für "{{term}}". Datenabfrage und Diagramm folgen im nächsten Schritt.', trends_context_placeholder: 'Kontextbereich als Struktur für "{{term}}". Kontextanalyse-Daten folgen im nächsten Schritt.', trends_article_list_placeholder: 'Artikelliste als Struktur für "{{term}}". Browser-ähnliches Tabellenverhalten/Paginierung folgt im nächsten Schritt.', trends_article_list_placeholder_date: 'Artikelliste als Struktur für "{{term}}" am {{date}}. Browser-ähnliches Tabellenverhalten/Paginierung folgt im nächsten Schritt.', trends_chart_mode_compact: 'Verlaufsdiagramm auf kompakt umschalten', trends_chart_mode_wide: 'Verlaufsdiagramm auf breit umschalten', trends_compare_title: 'Vergleich', trends_compare_body: 'Vergleichsansicht als Platzhalter. Domain-Vergleich folgt im nächsten Schritt.', compare_with_label: 'Vergleichen mit', connection_error_missing_url: 'Bitte geben Sie vor dem Verbindungstest eine Server-URL ein.', connection_error_invalid_url: 'Die Server-URL muss mit http:// oder https:// beginnen.', trends_search_title: 'Suche', trends_search_body: 'Suchansicht als Platzhalter. Trendsuche und Drill-down folgen als Nächstes.', selected: 'Ausgewählt', not_selected: 'Nicht ausgewählt', selected_trend_term: 'Ausgewählter Begriff', change_term: 'Begriff ändern', portal_empty_title: 'Noch keine Artikel', portal_empty_body: 'Diese Domain ist konfiguriert, aber das Backend hat keine Portal-Artikel zurückgegeben.', portal_load_error_title: 'Artikel konnten nicht geladen werden', portal_load_error: 'Die Mobile-App konnte dieses Portal gerade nicht laden. Bitte versuchen Sie es erneut.', retry: 'Erneut laden', open_project: 'News Deframer auf GitHub', ok: 'OK',
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
