import { useTranslation } from 'react-i18next';

import { Settings } from '../../shared/settings';
import { Theme } from '../../shared/theme';
import { HostStatus } from './StatusBadge';
import { ToggleSwitch } from './ToggleSwitch';

interface SettingsFormProps {
  settings: Settings;
  lang: string;
  status: HostStatus;
  onSettingsChange: (settings: Settings) => void;
  onLanguageChange: (language: string) => void;
  onTestConnection: () => void;
}

export const SettingsForm = ({ settings, lang, status, onSettingsChange, onLanguageChange, onTestConnection }: SettingsFormProps) => {
  const { t } = useTranslation();
  const isLoading = status === 'loading';
  const isSearchUrlValid = settings.searchEngineUrl?.startsWith('https://');

  return (
    <div className="content-grid">
      <div className="settings-column">
        <div className="card">
          <h3 className="section-title">{t('options.section_connection')}</h3>
          <div className="form-group">
            <label className="input-label">{t('options.label_server_url')}</label>
            <input type="text" className="text-input" value={settings.backendUrl} onChange={(e) => onSettingsChange({ ...settings, backendUrl: e.target.value })} placeholder="http://localhost:8080" />
          </div>
          <div className="form-group">
            <label className="input-label">
              {t('options.label_username')} <span className="optional-text">{t('options.label_optional')}</span>
            </label>
            <input type="text" className="text-input" value={settings.username} onChange={(e) => onSettingsChange({ ...settings, username: e.target.value })} />
          </div>
          <div className="form-group-last">
            <label className="input-label">
              {t('options.label_password')} <span className="optional-text">{t('options.label_optional')}</span>
            </label>
            <input type="password" className="text-input" value={settings.password} onChange={(e) => onSettingsChange({ ...settings, password: e.target.value })} />
          </div>
          <button className="action-button action-button-enabled" onClick={onTestConnection} disabled={isLoading} type="button">
            {isLoading ? t('options.btn_testing') : t('options.btn_test_connection')}
          </button>
        </div>

        <div className="card">
          <h3 className="section-title">{t('options.label_search_engine')}</h3>
          <div className="form-group-last">
            <input type="text" className={`text-input ${!isSearchUrlValid ? 'error' : ''}`} value={settings.searchEngineUrl || ''} onChange={(e) => onSettingsChange({ ...settings, searchEngineUrl: e.target.value })} />
            {!isSearchUrlValid && <div className="input-error-message">{t('options.error_https_only')}</div>}
          </div>
        </div>
      </div>

      <div className="settings-column">
        <div className="card">
          <h3 className="section-title">{t('options.language_label')}</h3>
          <select className="select-input" value={lang} onChange={(e) => onLanguageChange(e.target.value)}>
            <option value="default">{t('options.default')}</option>
            <option value="en">English</option>
            <option value="de">Deutsch</option>
          </select>
        </div>
        <div className="card">
          <h3 className="section-title">{t('options.theme_label')}</h3>
          <div className="theme-toggle-group">
            {(['light', 'dark', 'system'] as const).map((themeMode: Theme) => (
              <button key={themeMode} className={`theme-button ${settings.theme === themeMode ? 'active' : ''}`} onClick={() => onSettingsChange({ ...settings, theme: themeMode })} type="button">
                {themeMode === 'system' ? t('options.default') : t(`options.theme_${themeMode}`)}
              </button>
            ))}
          </div>
        </div>

        <div className="card">
          <h3 className="section-title">{t('options.section_general')}</h3>
          <ToggleSwitch
            id="settings-enable-extension"
            label={t('options.label_enable_extension')}
            checked={settings.enabled}
            onChange={(enabled) => onSettingsChange({ ...settings, enabled })}
          />
        </div>

        <div className="card">
          <h3 className="section-title">{t('options.section_project')}</h3>
          <a
            href="https://deframer.github.io/"
            target="_blank"
            rel="noopener noreferrer"
            className="project-link"
          >
            {t('footer.github_link')}
          </a>
        </div>
      </div>
    </div>
  );
};
