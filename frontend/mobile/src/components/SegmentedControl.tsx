import React from 'react';
import { Pressable, StyleSheet, Text, View } from 'react-native';

import { AppPalette } from '../theme';

export const SegmentedControl = ({
  options,
  value,
  onChange,
  palette,
}: {
  options: Array<{ label: string; value: string }>;
  value: string;
  onChange: (value: string) => void;
  palette: AppPalette;
}) => (
  <View style={[styles.wrapper, styles.group, { borderColor: palette.border, backgroundColor: palette.buttonBackground }]}>
    {options.map((option, index) => {
      const active = option.value === value;
      return (
        <Pressable
          key={option.value}
          onPress={() => onChange(option.value)}
          style={[
            styles.button,
            index < options.length - 1 ? [styles.buttonWithDivider, { borderRightColor: palette.border }] : null,
            active ? { backgroundColor: palette.accent } : null,
          ]}
        >
          <Text style={[styles.label, { color: active ? palette.accentText : palette.text }]}>{option.label}</Text>
        </Pressable>
      );
    })}
  </View>
);

const styles = StyleSheet.create({
  wrapper: {
    width: '100%',
  },
  group: {
    flexDirection: 'row',
    borderWidth: 1,
    borderRadius: 8,
    overflow: 'hidden',
  },
  button: {
    flex: 1,
    paddingVertical: 12,
    paddingHorizontal: 10,
    alignItems: 'center',
  },
  label: {
    fontSize: 16,
    fontWeight: '500',
  },
});
