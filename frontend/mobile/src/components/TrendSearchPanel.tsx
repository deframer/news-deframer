import React, { useEffect, useState } from 'react';
import { Search } from 'lucide-react-native';
import { Pressable, StyleSheet, TextInput, View } from 'react-native';
import { useTranslation } from 'react-i18next';

import { AnalyzedItem } from '../services/newsDeframerClient';
import { Settings } from '../services/settingsService';
import { AppPalette } from '../theme';
import { Card } from './Card';
import { TrendDetailsPanel } from './TrendDetailsPanel';

type TrendDetailTab = 'lifecycle' | 'context' | 'articles';

export const TrendSearchPanel = ({
  palette,
  domain,
  language,
  daysInPast,
  settings,
  onOpenArticle,
  onBackRequestChange,
}: {
  palette: AppPalette;
  domain: string;
  language: string;
  daysInPast: number;
  settings: Settings;
  onOpenArticle: (item: AnalyzedItem) => void;
  onBackRequestChange: (action: (() => void) | null) => void;
}) => {
  const { t } = useTranslation();
  const [term, setTerm] = useState('');
  const [activeTerm, setActiveTerm] = useState('');
  const [activeDetailTab, setActiveDetailTab] = useState<TrendDetailTab>('lifecycle');

  useEffect(() => {
    onBackRequestChange(null);

    return () => onBackRequestChange(null);
  }, [onBackRequestChange]);

  const handleSearch = () => {
    const normalized = term.trim();
    if (!normalized) {
      return;
    }

    setActiveTerm(normalized);
  };

  return (
    <View style={styles.stack}>
      <Card palette={palette}>
        <View style={styles.searchRow}>
          <TextInput
            value={term}
            onChangeText={setTerm}
            onSubmitEditing={handleSearch}
            placeholder={t('trends.search_placeholder', 'Enter term to analyze...')}
            placeholderTextColor={palette.secondaryText}
            style={[styles.searchInput, { color: palette.text, borderColor: palette.buttonBorder, backgroundColor: palette.background }]}
            returnKeyType="search"
          />
          <Pressable
            onPress={handleSearch}
            style={[styles.searchButton, { borderColor: palette.buttonBorder, backgroundColor: palette.buttonBackground }]}
            accessibilityRole="button"
            accessibilityLabel={t('trends.analyze', 'Search')}
          >
            <Search color={palette.text} size={18} strokeWidth={2.2} />
          </Pressable>
        </View>
      </Card>

      {activeTerm ? (
        <TrendDetailsPanel
          palette={palette}
          term={activeTerm}
          domain={domain}
          language={language}
          daysInPast={daysInPast}
          settings={settings}
          activeTab={activeDetailTab}
          setActiveTab={setActiveDetailTab}
          onOpenArticle={onOpenArticle}
        />
      ) : null}
    </View>
  );
};

const styles = StyleSheet.create({
  stack: { gap: 16 },
  searchRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 10,
  },
  searchInput: {
    flex: 1,
    borderWidth: 1,
    borderRadius: 8,
    minHeight: 44,
    paddingHorizontal: 12,
    fontSize: 16,
    fontWeight: '500',
  },
  searchButton: {
    width: 44,
    height: 44,
    borderWidth: 1,
    borderRadius: 8,
    alignItems: 'center',
    justifyContent: 'center',
  },
});
