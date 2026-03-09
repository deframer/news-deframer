import { fallbackLng, resources as sharedResources } from '@frontend-shared/i18n';
import i18n from 'i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import { initReactI18next } from 'react-i18next';

const hostResources = {
  en: {
    translation: {
      portal: {
        title: 'News Deframer',
        trends: 'Trends',
        hide: 'Hide',
      },
      article: {
        back: 'Back',
        back_tooltip: 'Go back to {{domain}} portal',
        hide: 'Hide',
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
      options: {
        title: 'Options',
        settings_title: 'Settings',
        loading: 'Loading...',
        status_connected: 'Connected',
        status_disabled: 'Disabled',
        status_error: 'Error',
        status_checking: 'Checking...',
        open_settings: 'Open Settings',
        apply: 'Apply',
      },
    },
  },
  de: {
    translation: {
      portal: {
        title: 'News Deframer',
        trends: 'Trends',
        hide: 'Ausblenden',
      },
      article: {
        back: 'Zurück',
        back_tooltip: 'Zurück zum {{domain}} Portal',
        hide: 'Ausblenden',
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
      options: {
        title: 'Einstellungen',
        settings_title: 'Einstellungen',
        loading: 'Laden...',
        status_connected: 'Verbunden',
        status_disabled: 'Deaktiviert',
        status_error: 'Fehler',
        status_checking: 'Prüfen...',
        open_settings: 'Einstellungen öffnen',
        apply: 'Anwenden',
      },
    },
  },
};

const resources = {
  en: {
    translation: {
      ...sharedResources.en.translation,
      ...hostResources.en.translation,
      article: {
        ...sharedResources.en.translation.article,
        ...hostResources.en.translation.article,
      },
      options: {
        ...sharedResources.en.translation.options,
        ...hostResources.en.translation.options,
      },
    },
  },
  de: {
    translation: {
      ...sharedResources.de.translation,
      ...hostResources.de.translation,
      article: {
        ...sharedResources.de.translation.article,
        ...hostResources.de.translation.article,
      },
      options: {
        ...sharedResources.de.translation.options,
        ...hostResources.de.translation.options,
      },
    },
  },
};

i18n.use(LanguageDetector).use(initReactI18next).init({
  resources,
  fallbackLng,
  interpolation: {
    escapeValue: false,
  },
  showSupportNotice: false,
});

export default i18n;
