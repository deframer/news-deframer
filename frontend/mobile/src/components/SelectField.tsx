import React, { useMemo, useState } from 'react';
import { Modal, Pressable, ScrollView, StyleSheet, Text, View } from 'react-native';

import { SECTION_HEADER_GAP } from './formStyles';
import { AppPalette } from '../theme';

const Chevron = ({ color }: { color: string }) => (
  <View style={[styles.chevron, { borderColor: color }]} />
);

export const SelectField = ({
  label,
  value,
  onValueChange,
  items,
  palette,
  hideLabel = false,
}: {
  label: string;
  value: string;
  onValueChange: (value: string) => void;
  items: Array<{ label: string; value: string }>;
  palette: AppPalette;
  hideLabel?: boolean;
}) => {
  const [open, setOpen] = useState(false);

  const selectedLabel = useMemo(() => items.find((item) => item.value === value)?.label ?? items[0]?.label ?? '', [items, value]);

  return (
    <>
      <View>
        {!hideLabel ? <Text style={[styles.label, { color: palette.secondaryText }]}>{label}</Text> : null}
        <Pressable
          onPress={() => setOpen(true)}
          style={[styles.trigger, { borderColor: palette.buttonBorder, backgroundColor: palette.background }]}
        >
          <Text style={[styles.triggerText, { color: palette.text }]}>{selectedLabel}</Text>
          <Chevron color={palette.text} />
        </Pressable>
      </View>

      <Modal transparent animationType="fade" visible={open} onRequestClose={() => setOpen(false)}>
        <Pressable style={styles.overlay} onPress={() => setOpen(false)}>
          <Pressable style={[styles.dialog, { backgroundColor: palette.card, borderColor: palette.border }]} onPress={() => {}}>
            <Text style={[styles.dialogTitle, { color: palette.secondaryText }]}>{label}</Text>
            <ScrollView>
              {items.map((item, index) => {
                const active = item.value === value;
                return (
                  <Pressable
                    key={item.value}
                    onPress={() => {
                      onValueChange(item.value);
                      setOpen(false);
                    }}
                    style={[
                      styles.option,
                      { borderBottomColor: palette.border },
                      active ? { backgroundColor: palette.accent } : null,
                      index === items.length - 1 ? styles.lastOption : null,
                    ]}
                  >
                    <Text style={[styles.optionText, { color: active ? palette.accentText : palette.text }]}>{item.label}</Text>
                  </Pressable>
                );
              })}
            </ScrollView>
          </Pressable>
        </Pressable>
      </Modal>
    </>
  );
};

const styles = StyleSheet.create({
  label: {
    marginBottom: SECTION_HEADER_GAP,
    fontSize: 14,
    fontWeight: '700',
  },
  trigger: {
    minHeight: 54,
    borderWidth: 1,
    borderRadius: 8,
    paddingHorizontal: 16,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  triggerText: {
    fontSize: 18,
    fontWeight: '500',
  },
  chevron: {
    width: 10,
    height: 10,
    borderRightWidth: 2,
    borderBottomWidth: 2,
    transform: [{ rotate: '45deg' }],
    marginTop: -4,
  },
  overlay: {
    flex: 1,
    backgroundColor: 'rgba(0,0,0,0.35)',
    justifyContent: 'center',
    padding: 24,
  },
  dialog: {
    borderWidth: 1,
    borderRadius: 12,
    overflow: 'hidden',
    maxHeight: 320,
  },
  dialogTitle: {
    paddingHorizontal: 16,
    paddingTop: 16,
    paddingBottom: 8,
    fontSize: 13,
    letterSpacing: 1,
    textTransform: 'uppercase',
    fontWeight: '700',
  },
  option: {
    paddingHorizontal: 16,
    paddingVertical: 16,
    borderBottomWidth: 1,
  },
  optionText: {
    fontSize: 18,
    fontWeight: '500',
  },
  lastOption: {
    borderBottomWidth: 0,
  },
});
