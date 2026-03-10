import { WordCloudInput } from './types';

export const buildWordCloudWords = <T>(
  items: T[],
  getText: (item: T) => string,
  getValue: (item: T) => number,
  getData: (item: T, index: number) => T = (item) => item,
): WordCloudInput<T>[] =>
  items
    .map((item, index) => ({
      text: getText(item),
      value: getValue(item),
      data: getData(item, index),
    }))
    .sort((a, b) => b.value - a.value);
