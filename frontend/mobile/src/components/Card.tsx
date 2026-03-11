import React, { PropsWithChildren } from 'react';
import { StyleSheet, View } from 'react-native';

import { AppPalette } from '../theme';

export const Card = ({ children, palette }: PropsWithChildren<{ palette: AppPalette }>) => (
  <View style={[styles.card, { backgroundColor: palette.card, borderColor: palette.border, boxShadow: `0px 6px 12px ${palette.cardShadow}` }]}>{children}</View>
);

const styles = StyleSheet.create({
  card: {
    borderWidth: 1,
    borderRadius: 8,
    padding: 16,
    elevation: 3,
  },
});
