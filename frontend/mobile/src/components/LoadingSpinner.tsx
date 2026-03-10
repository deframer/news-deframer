import React from 'react';
import { ActivityIndicator, StyleSheet, Text, View } from 'react-native';

import { AppPalette } from '../theme';

export const LoadingSpinner = ({
  palette,
  label,
  center = false,
}: {
  palette: AppPalette;
  label?: string;
  center?: boolean;
}) => (
  <View style={[styles.wrap, center ? styles.center : null]}>
    <ActivityIndicator size="small" color={palette.accent} />
    {label ? <Text style={[styles.label, { color: palette.secondaryText }]}>{label}</Text> : null}
  </View>
);

const styles = StyleSheet.create({
  wrap: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 10,
  },
  center: {
    justifyContent: 'center',
    alignItems: 'center',
    width: '100%',
  },
  label: {
    fontSize: 15,
    fontWeight: '500',
  },
});
