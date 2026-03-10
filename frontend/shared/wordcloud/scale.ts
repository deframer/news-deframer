export const scaleWordValue = (
  value: number,
  minValue: number,
  maxValue: number,
  minFontSize: number,
  maxFontSize: number,
): number => {
  const safeValue = Math.max(value, 0.1);
  const safeMin = Math.max(minValue, 0.1);
  const safeMax = Math.max(maxValue, safeMin * 1.1);

  if (safeMax === safeMin) {
    return (minFontSize + maxFontSize) / 2;
  }

  const ratio = (Math.log(safeValue) - Math.log(safeMin)) / (Math.log(safeMax) - Math.log(safeMin));
  return minFontSize + ratio * (maxFontSize - minFontSize);
};
