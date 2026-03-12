import React, { useMemo } from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { Dropdown } from 'react-native-element-dropdown';

import { AppPalette } from '../theme';

type DomainOption = { id: string; name: string };

export const TrendComparePanel = ({
  palette,
  availableDomains,
  compareDomain,
  onSelectDomain,
}: {
  palette: AppPalette;
  availableDomains: DomainOption[];
  compareDomain: string | null;
  onSelectDomain: (domain: string) => void;
}) => {
  const data = useMemo(() => availableDomains.map((option) => ({ label: option.name, value: option.id })), [availableDomains]);

  return (
    <View>
      <Text style={[styles.label, { color: palette.secondaryText }]}>Compare to</Text>
      <Dropdown
        mode="modal"
        data={data}
        search
        disable={data.length === 0}
        labelField="label"
        valueField="value"
        value={compareDomain}
        placeholder={data.length === 0 ? 'No domains available.' : 'Select domain'}
        searchPlaceholder="Search..."
        onChange={(item) => onSelectDomain(item.value)}
        style={[styles.dropdown, { borderColor: palette.buttonBorder, backgroundColor: palette.background }]}
        containerStyle={[styles.container, { backgroundColor: palette.card, borderColor: palette.border }]}
        placeholderStyle={[styles.placeholder, { color: palette.secondaryText }]}
        selectedTextStyle={[styles.selectedText, { color: palette.text }]}
        itemTextStyle={[styles.itemText, { color: palette.text }]}
        inputSearchStyle={[styles.inputSearch, { color: palette.text, borderColor: palette.buttonBorder, backgroundColor: palette.background }]}
        activeColor={palette.accent}
        iconColor={palette.text}
        backgroundColor="rgba(0,0,0,0.35)"
        maxHeight={320}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  label: {
    marginBottom: 8,
    fontSize: 14,
    letterSpacing: 1,
    textTransform: 'uppercase',
    fontWeight: '700',
  },
  dropdown: {
    minHeight: 54,
    borderWidth: 1,
    borderRadius: 8,
    paddingHorizontal: 16,
  },
  container: {
    borderWidth: 1,
    borderRadius: 12,
    overflow: 'hidden',
  },
  placeholder: {
    fontSize: 18,
    fontWeight: '500',
  },
  selectedText: {
    fontSize: 18,
    fontWeight: '500',
  },
  itemText: {
    fontSize: 18,
    fontWeight: '500',
  },
  inputSearch: {
    minHeight: 46,
    borderWidth: 1,
    borderRadius: 8,
    margin: 12,
    paddingHorizontal: 14,
    fontSize: 16,
  },
});
