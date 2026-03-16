import { analyzeEmotionVectorToCodes, translateCodeMap } from '../shared/sentiments/index.ts';


if (false) {
  const sampleInput = {
    valence: 7.72,
    arousal: 4.77,
    dominance: 6.92,
    joy: 3.96,
    anger: 1.24,
    sadness: 1.29,
    fear: 1.37,
    disgust: 1.15,
  };

  console.log('--- Sentiment Analysis CLI ---');
  console.log('Input Vector:', JSON.stringify(sampleInput, null, 2));

  const codes = analyzeEmotionVectorToCodes(sampleInput);
  console.log('\nGenerated Codes:', JSON.stringify(codes, null, 2));

  const translatedDe = translateCodeMap(codes, 'de');
  console.log('\nResults (Deutsch):');
  Object.entries(translatedDe).forEach(([key, value]) => {
    console.log(`  ${key}: ${value}`);
  });

  const translatedEn = translateCodeMap(codes, 'en');
  console.log('\nResults (English):');
  Object.entries(translatedEn).forEach(([key, value]) => {
    console.log(`  ${key}: ${value}`);
  });

}