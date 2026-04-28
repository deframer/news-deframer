import React from 'react';
import { Pressable, StyleSheet, Text, View } from 'react-native';

import { AppPalette } from '../theme';

export const ToggleSwitch = ({
  label,
  checked,
  onChange,
  palette,
  description,
}: {
  label: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
  palette: AppPalette;
  description?: string;
}) => (
  <View style={styles.row}>
    <Pressable
      accessibilityRole="switch"
      accessibilityState={{ checked }}
      onPress={() => onChange(!checked)}
      style={({ pressed }) => [styles.switchHitArea, pressed ? { opacity: 0.85 } : null]}
    >
      <View style={[styles.track, { backgroundColor: checked ? palette.accent : palette.buttonBorder }]}>
        <View style={[styles.thumb, { backgroundColor: palette.card, alignSelf: checked ? 'flex-end' : 'flex-start' }]} />
      </View>
    </Pressable>
    <View style={styles.textColumn}>
      <Text style={[styles.label, { color: palette.text }]}>{label}</Text>
      {description ? <Text style={[styles.description, { color: palette.secondaryText }]}>{description}</Text> : null}
    </View>
  </View>
);

const styles = StyleSheet.create({
  row: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
  },
  switchHitArea: {
    paddingVertical: 2,
  },
  track: {
    width: 48,
    height: 28,
    borderRadius: 999,
    padding: 3,
    justifyContent: 'center',
  },
  thumb: {
    width: 22,
    height: 22,
    borderRadius: 11,
  },
  textColumn: {
    flex: 1,
    gap: 4,
  },
  label: {
    fontSize: 16,
    fontWeight: '600',
  },
  description: {
    fontSize: 13,
    lineHeight: 18,
  },
});
