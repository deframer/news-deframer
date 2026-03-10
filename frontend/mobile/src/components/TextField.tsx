import React from 'react';
import { StyleSheet, Text, TextInput, View } from 'react-native';

import { AppPalette } from '../theme';

export const TextField = ({
  label,
  optional,
  value,
  onChangeText,
  palette,
  secureTextEntry,
  error,
  hideLabel = false,
}: {
  label: string;
  optional?: string;
  value: string;
  onChangeText: (text: string) => void;
  palette: AppPalette;
  secureTextEntry?: boolean;
  error?: string;
  hideLabel?: boolean;
}) => (
  <View style={styles.wrapper}>
    {!hideLabel ? (
      <Text style={[styles.label, { color: palette.text }]}> 
        {label}
        {optional ? <Text style={{ color: palette.secondaryText }}> {optional}</Text> : null}
      </Text>
    ) : null}
    <TextInput
      value={value}
      onChangeText={onChangeText}
      secureTextEntry={secureTextEntry}
      style={[
        styles.input,
        {
          color: palette.text,
          borderColor: error ? palette.danger : palette.buttonBorder,
          backgroundColor: palette.background,
        },
      ]}
      autoCapitalize="none"
    />
    {error ? <Text style={[styles.error, { color: palette.danger }]}>{error}</Text> : null}
  </View>
);

const styles = StyleSheet.create({
  wrapper: { marginBottom: 16 },
  label: { marginBottom: 6, fontSize: 16, fontWeight: '500' },
  input: {
    borderWidth: 1,
    borderRadius: 8,
    paddingHorizontal: 14,
    paddingVertical: 12,
    fontSize: 18,
  },
  error: { marginTop: 6, fontSize: 12 },
});
