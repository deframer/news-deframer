export interface WordCloudInput<T> {
  text: string;
  value: number;
  data: T;
}

export interface WordMeasurement {
  width: number;
  height: number;
}

export interface WordCloudWord<T> extends WordCloudInput<T> {
  rank: number;
  fontSize: number;
  color: string;
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface MeasureWordOptions {
  text: string;
  fontSize: number;
  fontFamily: string;
  fontWeight?: string;
}

export type MeasureWord = (options: MeasureWordOptions) => WordMeasurement;

export interface LayoutWordCloudOptions<T> {
  words: WordCloudInput<T>[];
  width: number;
  height: number;
  measureWord: MeasureWord;
  getColor: (text: string, index: number) => string;
  minFontSize?: number;
  maxFontSize?: number;
  padding?: number;
  fontFamily?: string;
  fontWeight?: string;
  spiralStep?: number;
  angleStep?: number;
  maxPlacementSteps?: number;
}
