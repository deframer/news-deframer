import { useTranslation } from 'react-i18next';

import { PROMPTS_URL, REFERENCE_TAB_TARGET } from '../../shared/links';

interface AnalysisSourceProps {
  llmModel?: string | null;
}

export const AnalysisSource = ({ llmModel }: AnalysisSourceProps) => {
  const { t } = useTranslation();

  return (
    <section className="analysis-source">
      <p
        className="analysis-source-line"
        style={{
          display: 'flex',
          flexWrap: 'wrap',
          justifyContent: 'center',
          alignItems: 'baseline',
          gap: '0.35rem',
          width: '100%',
          margin: 0,
          color: 'var(--secondary-text)',
          lineHeight: 1.45,
        }}
        >
        {llmModel ? (
          <>
            <span className="analysis-source-label">
              {t('article.analysis_source_llm_model_label')}
            </span>
            <em
              className="analysis-source-value"
            >
              {llmModel}
            </em>
            {' / '}
          </>
        ) : null}
        <a
          className="analysis-source-link"
          href={PROMPTS_URL}
          target={REFERENCE_TAB_TARGET}
        >
          {t('article.analysis_source_prompts_link')}
        </a>
      </p>
    </section>
  );
};
