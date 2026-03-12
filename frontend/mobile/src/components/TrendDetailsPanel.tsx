import React from 'react';
import { View } from 'react-native';
import { useTranslation } from 'react-i18next';

import { AppPalette } from '../theme';
import { AnalyzedItem } from '../services/newsDeframerClient';
import { Settings } from '../services/settingsService';
import { Card } from './Card';
import { SegmentedControl } from './SegmentedControl';
import { TrendArticleListPanel } from './TrendArticleListPanel';
import { TrendContextPanel } from './TrendContextPanel';
import { TrendLifecyclePanel } from './TrendLifecyclePanel';

type TrendDetailTab = 'lifecycle' | 'context' | 'articles';

export const TrendDetailsPanel = ({
  palette,
  term,
  domain,
  language,
  daysInPast,
  settings,
  activeTab,
  setActiveTab,
  onOpenArticle,
}: {
  palette: AppPalette;
  term: string;
  domain: string;
  language: string;
  daysInPast: number;
  settings: Settings;
  activeTab: TrendDetailTab;
  setActiveTab: (tab: TrendDetailTab) => void;
  onOpenArticle: (item: AnalyzedItem) => void;
}) => {
  const { t } = useTranslation();

  return (
    <Card palette={palette}>
      <View>
        <SegmentedControl
          palette={palette}
          value={activeTab}
          onChange={(value) => setActiveTab(value as TrendDetailTab)}
          options={[
            { label: t('trends.lifecycle'), value: 'lifecycle' },
            { label: t('trends.context'), value: 'context' },
            { label: t('trends.articles'), value: 'articles' },
          ]}
        />
      </View>

      {activeTab === 'lifecycle' ? <TrendLifecyclePanel palette={palette} term={term} domain={domain} language={language} daysInPast={daysInPast} settings={settings} onOpenArticle={onOpenArticle} /> : null}
      {activeTab === 'context' ? <TrendContextPanel palette={palette} term={term} domain={domain} language={language} daysInPast={daysInPast} settings={settings} /> : null}
      {activeTab === 'articles' ? <TrendArticleListPanel palette={palette} term={term} domain={domain} settings={settings} daysInPast={daysInPast} headerTitle={null} onOpenArticle={onOpenArticle} /> : null}
    </Card>
  );
};
