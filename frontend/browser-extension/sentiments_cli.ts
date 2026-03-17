// npx ts-node sentiments_cli.ts

import { sentimentsToLabels, analyzeEmotionVectorToCodes } from '../shared/sentiments/sentimentLabels.ts';
import type { EmotionVector } from '../shared/sentiments/sentimentLabels.ts';

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
const de = sentimentsToLabels(input, "de");

console.log('Codes:', codes);
console.log('Translations (DE):', de);