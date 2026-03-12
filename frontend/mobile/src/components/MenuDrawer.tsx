import React from 'react';
import { Modal, Pressable, StyleSheet, Text } from 'react-native';

import { AppPalette } from '../theme';

export const MenuDrawer = ({
  visible,
  palette,
  onClose,
  onNavigate,
  labels,
}: {
  visible: boolean;
  palette: AppPalette;
  onClose: () => void;
  onNavigate: (screen: 'dashboard' | 'settings' | 'about') => void;
  labels: { settings: string; about: string };
}) => {
  const items: Array<{ screen: 'settings' | 'about'; label: string }> = [
    { screen: 'settings', label: labels.settings },
    { screen: 'about', label: labels.about },
  ];

  return (
    <Modal animationType="fade" transparent visible={visible} onRequestClose={onClose}>
      <Pressable style={styles.overlay} onPress={onClose}>
        <Pressable style={[styles.sheet, { backgroundColor: palette.card, borderColor: palette.border }]} onPress={() => {}}>
          {items.map((item, index) => (
            <Pressable
              key={item.screen}
              onPress={() => onNavigate(item.screen)}
              style={[
                styles.item,
                { borderBottomColor: palette.border },
                index === items.length - 1 ? styles.lastItem : null,
              ]}
            >
              <Text style={[styles.itemText, { color: palette.text }]}>{item.label}</Text>
            </Pressable>
          ))}
        </Pressable>
      </Pressable>
    </Modal>
  );
};

const styles = StyleSheet.create({
  overlay: {
    flex: 1,
    backgroundColor: 'rgba(0,0,0,0.4)',
    justifyContent: 'flex-start',
    alignItems: 'flex-start',
  },
  sheet: {
    marginTop: 72,
    marginLeft: 16,
    width: 240,
    borderWidth: 1,
    borderRadius: 12,
    overflow: 'hidden',
  },
  item: {
    paddingHorizontal: 16,
    paddingVertical: 16,
    borderBottomWidth: 1,
  },
  itemText: {
    fontSize: 18,
    fontWeight: '600',
  },
  lastItem: {
    borderBottomWidth: 0,
  },
});
