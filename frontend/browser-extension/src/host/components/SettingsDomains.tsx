import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { getDomain } from 'tldts';

import { DomainEntry } from '../../ndf/client';
import { Settings } from '../../shared/settings';

interface SettingsDomainsProps {
  domains: DomainEntry[];
  domainsLoading: boolean;
  domainsUnavailable: boolean;
  settings: Settings;
  onSelectedDomainsChange: (selectedDomains: string[]) => void;
}

const normalizeDomain = (value: string) => getDomain(value) || value;

export const SettingsDomains = ({ domains, domainsLoading, domainsUnavailable, settings, onSelectedDomainsChange }: SettingsDomainsProps) => {
  const { t } = useTranslation();
  const selectedDomains = settings.selectedDomains || [];
  const extensionDisabled = !settings.enabled;

  const domainOptions = useMemo(
    () => domains.map((entry) => ({ ...entry, rootDomain: normalizeDomain(entry.domain) })),
    [domains]
  );
  const allRootDomains = useMemo(
    () => Array.from(new Set(domainOptions.map((entry) => entry.rootDomain))),
    [domainOptions]
  );

  const handleToggle = (rootDomain: string, checked: boolean) => {
    const nextSelectedDomains = checked
      ? Array.from(new Set([...selectedDomains, rootDomain]))
      : selectedDomains.filter((domain) => domain !== rootDomain);

    onSelectedDomainsChange(nextSelectedDomains);
  };

  return (
    <div className="card">
      <div className="domains-header">
        <h2 className="domains-title">{t('options.domains_title')}</h2>
        {!domainsLoading && !domainsUnavailable && !extensionDisabled && domainOptions.length > 0 ? (
          <div className="domains-actions">
            <button className="domains-select-all" type="button" onClick={() => onSelectedDomainsChange(allRootDomains)}>
              {t('options.domains_select_all')}
            </button>
            <button className="domains-select-all" type="button" onClick={() => onSelectedDomainsChange([])}>
              {t('options.domains_unselect_all')}
            </button>
          </div>
        ) : null}
      </div>
      {domainsLoading ? <div className="domains-state">{t('options.domains_loading')}</div> : null}

      {!domainsLoading && extensionDisabled ? <div className="domains-state domains-state-error">{t('options.domains_disabled')}</div> : null}

      {!domainsLoading && !extensionDisabled && domainsUnavailable ? <div className="domains-state domains-state-error">{t('options.domains_unavailable')}</div> : null}

      {!domainsLoading && !domainsUnavailable && !extensionDisabled && domainOptions.length > 0 ? (
        <div className="domains-list" role="group" aria-label={t('options.domains_title')}>
          {domainOptions.map((entry) => {
            const checked = selectedDomains.includes(entry.rootDomain);

            return (
              <label key={`${entry.domain}-${entry.language}`} className="domain-option">
                <input
                  className="domain-checkbox-input"
                  type="checkbox"
                  checked={checked}
                  onChange={(e) => handleToggle(entry.rootDomain, e.target.checked)}
                />
                <span className={`domain-checkbox ${checked ? 'checked' : ''}`} aria-hidden="true">
                  <span className="domain-checkbox-mark" />
                </span>
                <span className="domain-option-content">
                  <span className="domain-option-header">
                    <span className="domain-option-name">{entry.rootDomain}</span>
                    <span className="domain-option-language">{entry.language}</span>
                  </span>
                  {entry.domain !== entry.rootDomain ? <span className="domain-option-source">{entry.domain}</span> : null}
                </span>
              </label>
            );
          })}
        </div>
      ) : null}
    </div>
  );
};
