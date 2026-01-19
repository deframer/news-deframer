export const formatRatingPercent = (rating: number | undefined): number => Math.round((rating || 0.0) * 100);

export const getRatingColors = (percentage: number): { bg: string; text: string } => {
  if (percentage < 34) return { bg: '#198754', text: '#ffffff' }; // Accessible Green
  if (percentage < 67) return { bg: '#ffc107', text: '#000000' }; // Accessible Yellow
  return { bg: '#b02a37', text: '#ffffff' }; // Accessible Red
};