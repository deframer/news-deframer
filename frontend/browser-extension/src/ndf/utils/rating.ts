export interface RatingColorStyles {
  bg: string;
  text: string;
}

export const getRatingColors = (percentage: number): RatingColorStyles => {
  // Use default text color for very low ratings to ensure visibility on the empty track
  if (percentage < 11) return { bg: 'var(--success-color)', text: 'var(--text-color)' };
  if (percentage < 34) return { bg: 'var(--success-color)', text: '#ffffff' }; // Accessible Green
  if (percentage < 67) return { bg: 'var(--warning-color)', text: '#000000' }; // Accessible Yellow
  return { bg: 'var(--danger-color)', text: '#ffffff' }; // Accessible Red
};

export const formatRatingPercent = (rating: number | undefined): number =>
  Math.round((rating || 0.0) * 100);
