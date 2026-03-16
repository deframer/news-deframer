import { analyzeEmotionVectorToCodes, translateCodeMap } from '../shared/sentiments/index.ts';
import type { EmotionVector } from '../shared/sentiments/index.ts';

// optional
const input: EmotionVector = {
  valence: 7.72,
  arousal: 4.77,
  dominance: 6.92,
  joy: 3.96,
  anger: 1.24,
  sadness: 1.29,
  fear: 1.37,
  disgust: 1.15,
};

const codes = analyzeEmotionVectorToCodes(input);
const de = translateCodeMap(codes, "de");

console.log(codes);
console.log(de);