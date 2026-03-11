import React from 'react';
import { StyleSheet, Text, View } from 'react-native';

import { AppPalette } from '../theme';

export type HostStatus = 'idle' | 'loading' | 'success' | 'error';

export const StatusBadge = ({ status, palette, labels }: { status: HostStatus; palette: AppPalette; labels: { connected: string; error: string; checking: string } }) => {
  if (status === 'idle') {
    return null;
  }

  const color = status === 'success' ? palette.success : status === 'error' ? palette.danger : palette.warning;
  const label = status === 'success' ? labels.connected : status === 'error' ? labels.error : labels.checking;

  return (
    <View style={[styles.badge, { borderColor: color, backgroundColor: palette.background }]}> 
      <View style={[styles.dot, { backgroundColor: color }]} />
      <Text style={[styles.text, { color }]}>{label}</Text>
    </View>
  );
};

const styles = StyleSheet.create({
  badge: {
    flexDirection: 'row',
    alignItems: 'center',
    borderWidth: 1,
    borderRadius: 999,
    paddingHorizontal: 12,
    paddingVertical: 6,
  },
  dot: {
    width: 8,
    height: 8,
    borderRadius: 4,
    marginRight: 8,
  },
  text: {
    fontSize: 12,
    fontWeight: '700',
  },
});
