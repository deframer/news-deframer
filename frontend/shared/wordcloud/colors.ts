const DEFAULT_COLORS = ['var(--trend-up)', 'var(--trend-down)', 'var(--trend-steady)'] as const;

export const getWordCloudColor = (text: string, colors: readonly string[] = DEFAULT_COLORS): string => {
  let hash = 0;

  for (let i = 0; i < text.length; i += 1) {
    hash = text.charCodeAt(i) + ((hash << 5) - hash);
  }

  return colors[Math.abs(hash) % colors.length];
};

export const wordCloudColors = DEFAULT_COLORS;
