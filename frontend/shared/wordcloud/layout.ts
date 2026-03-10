import { scaleWordValue } from './scale';
import { LayoutWordCloudOptions, WordCloudWord } from './types';

// Inspired by the placement approach behind visx/wordcloud and d3-cloud,
// but implemented here as a small DOM-free layout engine that can be reused
// from both web and React Native wrappers.

interface Box {
  left: number;
  top: number;
  right: number;
  bottom: number;
}

const intersects = (a: Box, b: Box): boolean =>
  !(a.right <= b.left || a.left >= b.right || a.bottom <= b.top || a.top >= b.bottom);

const getWordBox = <T>(word: WordCloudWord<T>): Box => ({
  left: word.x - word.width / 2,
  top: word.y - word.height / 2,
  right: word.x + word.width / 2,
  bottom: word.y + word.height / 2,
});

const recenterPlacedWords = <T>(words: WordCloudWord<T>[], width: number, height: number): WordCloudWord<T>[] => {
  if (words.length === 0) {
    return words;
  }

  const bounds = words.reduce(
    (acc, word) => {
      const box = getWordBox(word);
      return {
        left: Math.min(acc.left, box.left),
        top: Math.min(acc.top, box.top),
        right: Math.max(acc.right, box.right),
        bottom: Math.max(acc.bottom, box.bottom),
      };
    },
    { left: Infinity, top: Infinity, right: -Infinity, bottom: -Infinity },
  );

  const dx = width / 2 - (bounds.left + bounds.right) / 2;
  const dy = height / 2 - (bounds.top + bounds.bottom) / 2;

  const centeredWords = words.map((word) => ({
    ...word,
    x: word.x + dx,
    y: word.y + dy,
  }));

  const centeredBounds = centeredWords.reduce(
    (acc, word) => {
      const box = getWordBox(word);
      return {
        left: Math.min(acc.left, box.left),
        top: Math.min(acc.top, box.top),
        right: Math.max(acc.right, box.right),
        bottom: Math.max(acc.bottom, box.bottom),
      };
    },
    { left: Infinity, top: Infinity, right: -Infinity, bottom: -Infinity },
  );

  const clampDx = centeredBounds.left < 0 ? -centeredBounds.left : centeredBounds.right > width ? width - centeredBounds.right : 0;
  const clampDy = centeredBounds.top < 0 ? -centeredBounds.top : centeredBounds.bottom > height ? height - centeredBounds.bottom : 0;

  if (clampDx === 0 && clampDy === 0) {
    return centeredWords;
  }

  return centeredWords.map((word) => ({
    ...word,
    x: word.x + clampDx,
    y: word.y + clampDy,
  }));
};

export const layoutWordCloud = <T>({
  words,
  width,
  height,
  measureWord,
  getColor,
  minFontSize = 14,
  maxFontSize = 50,
  padding = 4,
  fontFamily = 'sans-serif',
  fontWeight = 'normal',
  spiralStep = 2.5,
  angleStep = 0.35,
  maxPlacementSteps = 1200,
}: LayoutWordCloudOptions<T>): WordCloudWord<T>[] => {
  if (width <= 0 || height <= 0 || words.length === 0) {
    return [];
  }

  const values = words.map((word) => word.value);
  const minValue = Math.min(...values);
  const maxValue = Math.max(...values);
  const placedBoxes: Box[] = [];
  const centerX = width / 2;
  const centerY = height / 2;
  const horizontalScale = width >= height ? width / height : 1;
  const verticalScale = height > width ? height / width : 1;

  const placedWords = words.reduce<WordCloudWord<T>[]>((result, word, index) => {
    const fontSize = scaleWordValue(word.value, minValue, maxValue, minFontSize, maxFontSize);
    const measured = measureWord({
      text: word.text,
      fontSize,
      fontFamily,
      fontWeight,
    });

    const wordWidth = measured.width + padding * 2;
    const wordHeight = measured.height + padding * 2;

    let step = 0;
    let placedWord: WordCloudWord<T> | null = null;

    while (step < maxPlacementSteps) {
      const angle = step * angleStep;
      const radius = spiralStep * angle;
      const x = centerX + Math.cos(angle) * radius * horizontalScale;
      const y = centerY + Math.sin(angle) * radius * verticalScale;
      const box = {
        left: x - wordWidth / 2,
        top: y - wordHeight / 2,
        right: x + wordWidth / 2,
        bottom: y + wordHeight / 2,
      };

      const fitsBounds = box.left >= 0 && box.top >= 0 && box.right <= width && box.bottom <= height;
      const collides = placedBoxes.some((otherBox) => intersects(box, otherBox));

      if (fitsBounds && !collides) {
        placedBoxes.push(box);
        placedWord = {
          ...word,
          rank: index + 1,
          fontSize,
          color: getColor(word.text, index),
          x,
          y,
          width: measured.width,
          height: measured.height,
        };
        break;
      }

      step += 1;
    }

    if (placedWord) {
      result.push(placedWord);
    }

    return result;
  }, []);

  return recenterPlacedWords(placedWords, width, height);
};
