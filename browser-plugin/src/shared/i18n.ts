import i18n from 'i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import { initReactI18next } from 'react-i18next';

const resources = {
  en: {
    translation: {
      portal: {
        title: 'News Deframer: {{domain}}',
        hide: 'Hide',
      },
      article: {
        back: 'Back',
        back_tooltip: 'Go back to {{domain}} portal',
        hide: 'Hide',
        no_title: 'No title',
        no_description: 'No description',
        original_section: 'Original',
        btn_original_title: 'Original Title',
        btn_details: 'Details',
        btn_view_original: 'View Original',
      },
      metrics: {
        framing: 'Framing',
        clickbait: 'Clickbait',
        persuasive: 'Persuasive',
        hyper_stimulus: 'Hyper Stimulus',
        speculative: 'Speculative',
        overall_rating: 'Overall Rating',
      },
      rating: {
        label: 'Rating',
        aria_label: '{{label}}: {{percentage}}%',
        overlay_aria_label: 'Overall rating: {{percentage}}%. Reason: {{reason}}',
        no_reason: 'No reason provided.',
      },
      footer: {
        text_pre: 'This content was replaced by the ',
        text_post: ' browser plugin.',
      },
      metadata: {
        just_now: 'just now',
      },
      options: {
        title: 'Options',
        language_label: 'Language',
        theme_label: 'Theme',
        theme_light: 'Light',
        theme_dark: 'Dark',
        theme_system: 'System',
        system: 'System',
        default: 'Default',
        loading: 'Loading...',
        status_connected: 'Connected',
        status_error: 'Error',
        status_checking: 'Checking...',
        section_connection: 'Connection',
        label_server_url: 'Server URL',
        label_username: 'Username',
        label_password: 'Password',
        label_optional: '(Optional)',
        btn_test_connection: 'Test Connection',
        btn_testing: 'Testing...',
        section_general: 'General',
        label_enable_extension: 'Enable Extension',
      },
    },
  },
  de: {
    translation: {
      portal: {
        title: 'News Deframer: {{domain}}',
        hide: 'Ausblenden',
      },
      article: {
        back: 'Zurück',
        back_tooltip: 'Zurück zum {{domain}} Portal',
        hide: 'Ausblenden',
        no_title: 'Kein Titel',
        no_description: 'Keine Beschreibung',
        original_section: 'Original',
        btn_original_title: 'Originaltitel',
        btn_details: 'Details',
        btn_view_original: 'Original ansehen',
      },
      metrics: {
        framing: 'Framing',
        clickbait: 'Clickbait',
        persuasive: 'Beeinflussung',
        hyper_stimulus: 'Reizüberflutung',
        speculative: 'Spekulativ',
        overall_rating: 'Gesamtbewertung',
      },
      rating: {
        label: 'Bewertung',
        aria_label: '{{label}}: {{percentage}}%',
        overlay_aria_label: 'Gesamtbewertung: {{percentage}}%. Grund: {{reason}}',
        no_reason: 'Kein Grund angegeben.',
      },
      footer: {
        text_pre: 'Dieser Inhalt wurde durch das ',
        text_post: ' Browser-Plugin ersetzt.',
      },
      metadata: {
        just_now: 'gerade eben',
      },
      options: {
        title: 'Einstellungen',
        language_label: 'Sprache',
        theme_label: 'Design',
        theme_light: 'Hell',
        theme_dark: 'Dunkel',
        theme_system: 'System',
        system: 'System',
        default: 'Standard',
        loading: 'Laden...',
        status_connected: 'Verbunden',
        status_error: 'Fehler',
        status_checking: 'Prüfen...',
        section_connection: 'Verbindung',
        label_server_url: 'Server-URL',
        label_username: 'Benutzername',
        label_password: 'Passwort',
        label_optional: '(Optional)',
        btn_test_connection: 'Verbindung testen',
        btn_testing: 'Teste...',
        section_general: 'Allgemein',
        label_enable_extension: 'Erweiterung aktivieren',
      },
    },
  },
};

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    fallbackLng: 'en',
    interpolation: {
      escapeValue: false,
    },
  });

export default i18n;