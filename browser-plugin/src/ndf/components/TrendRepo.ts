export interface TrendMetric {
  trend_topic: string;
  frequency: number;
  utility: number;
  outlier_ratio: number;
  time_slice: string;
}

export interface TrendContextMetric {
  context_word: string;
  frequency: number;
  type: string;
}

export interface TrendComparisonMetric {
  classification: 'BLINDSPOT_A' | 'BLINDSPOT_B' | 'INTERSECT';
  rank_group: number;
  trend_topic: string;
  score_a: number;
  score_b: number;
}

export interface TrendLifecycleMetric {
  time_slice: string;
  frequency: number;
  velocity: number;
}

const STATIC_TREND_DATA: TrendMetric[] = [
  { trend_topic: 'berlin', frequency: 3, utility: 1, outlier_ratio: 2.333333333333333, time_slice: '2026-01-30 00:00:00+00' },
  { trend_topic: 'deutschland', frequency: 3, utility: 1, outlier_ratio: 2.1, time_slice: '2026-01-30 00:00:00+00' },
  { trend_topic: 'karin', frequency: 2, utility: 1, outlier_ratio: 2, time_slice: '2026-01-29 00:00:00+00' },
  { trend_topic: 'hamburg', frequency: 2, utility: 1, outlier_ratio: 2, time_slice: '2026-01-30 00:00:00+00' },
  { trend_topic: 'syrer', frequency: 2, utility: 1, outlier_ratio: 2, time_slice: '2026-02-03 00:00:00+00' },
  { trend_topic: 'prien', frequency: 2, utility: 1, outlier_ratio: 2, time_slice: '2026-01-29 00:00:00+00' },
  { trend_topic: 'köln', frequency: 2, utility: 1, outlier_ratio: 2, time_slice: '2026-02-03 00:00:00+00' },
  { trend_topic: 'oeter', frequency: 2, utility: 1, outlier_ratio: 2, time_slice: '2026-01-29 00:00:00+00' },
  { trend_topic: 'live', frequency: 2, utility: 1, outlier_ratio: 2, time_slice: '2026-02-03 00:00:00+00' },
  { trend_topic: 'regierung', frequency: 2, utility: 1, outlier_ratio: 2, time_slice: '2026-02-03 00:00:00+00' },
  { trend_topic: 'gericht', frequency: 2, utility: 1, outlier_ratio: 2, time_slice: '2026-01-31 00:00:00+00' },
  { trend_topic: 'merz', frequency: 2, utility: 1, outlier_ratio: 1.75, time_slice: '2026-01-30 00:00:00+00' },
  { trend_topic: 'land', frequency: 2, utility: 1, outlier_ratio: 1.7142857142857142, time_slice: '2026-02-03 00:00:00+00' },
  { trend_topic: 'nius', frequency: 3, utility: 1, outlier_ratio: 1.6153846153846154, time_slice: '2026-02-02 00:00:00+00' },
  { trend_topic: 'frau', frequency: 2, utility: 1, outlier_ratio: 1.6, time_slice: '2026-02-02 00:00:00+00' },
  { trend_topic: 'friedrich', frequency: 2, utility: 1, outlier_ratio: 1.5555555555555554, time_slice: '2026-01-30 00:00:00+00' },
  { trend_topic: 'cdu', frequency: 2, utility: 1, outlier_ratio: 1.5555555555555554, time_slice: '2026-01-31 00:00:00+00' }
];

const STATIC_CONTEXT_DATA: TrendContextMetric[] = [
  { context_word: 'anlaufen', frequency: 1, type: 'VERB' },
  { context_word: 'besänftigen', frequency: 1, type: 'VERB' },
  { context_word: 'einmischen', frequency: 1, type: 'VERB' },
  { context_word: 'engagieren', frequency: 1, type: 'VERB' },
  { context_word: 'fernbleiben', frequency: 1, type: 'VERB' },
  { context_word: 'funktionieren', frequency: 1, type: 'VERB' },
  { context_word: 'jagen', frequency: 1, type: 'VERB' },
  { context_word: 'lehnen', frequency: 1, type: 'VERB' },
  { context_word: 'mitproduzieren', frequency: 1, type: 'VERB' },
  { context_word: 'nehmen', frequency: 1, type: 'VERB' },
];

const STATIC_COMPARISON_DATA: TrendComparisonMetric[] = [
  { classification: 'BLINDSPOT_A', rank_group: 1, trend_topic: 'hans', score_a: 5, score_b: 0 },
  { classification: 'BLINDSPOT_A', rank_group: 2, trend_topic: 'müller', score_a: 5, score_b: 0 },
  { classification: 'BLINDSPOT_A', rank_group: 3, trend_topic: 'sprechen', score_a: 3, score_b: 0 },
  { classification: 'BLINDSPOT_A', rank_group: 4, trend_topic: 'borg', score_a: 2, score_b: 0 },
  { classification: 'BLINDSPOT_A', rank_group: 5, trend_topic: 'høiby', score_a: 2, score_b: 0 },
  { classification: 'BLINDSPOT_B', rank_group: 1, trend_topic: 'grenzübergang', score_a: 0, score_b: 7 },
  { classification: 'BLINDSPOT_B', rank_group: 2, trend_topic: 'rafah', score_a: 0, score_b: 6 },
  { classification: 'BLINDSPOT_B', rank_group: 3, trend_topic: 'gespräch', score_a: 0, score_b: 5 },
  { classification: 'BLINDSPOT_B', rank_group: 4, trend_topic: 'gazastreifen', score_a: 0, score_b: 5 },
  { classification: 'BLINDSPOT_B', rank_group: 5, trend_topic: 'öffnen', score_a: 0, score_b: 4 },
  { classification: 'INTERSECT', rank_group: 1, trend_topic: 'deutschland', score_a: 6, score_b: 6.18 },
  { classification: 'INTERSECT', rank_group: 2, trend_topic: 'frankreich', score_a: 6, score_b: 3.6 },
  { classification: 'INTERSECT', rank_group: 3, trend_topic: 'auto', score_a: 6, score_b: 3 },
  { classification: 'INTERSECT', rank_group: 4, trend_topic: 'mensch', score_a: 3, score_b: 4.5 },
  { classification: 'INTERSECT', rank_group: 5, trend_topic: 'kuh', score_a: 3, score_b: 2.55 },
];

const AVAILABLE_DOMAINS = [
  { id: 'tagesschau.de', name: 'Tagesschau' },
  { id: 'spiegel.de', name: 'Spiegel' },
  { id: 'zeit.de', name: 'Zeit' },
  { id: 'welt.de', name: 'Welt' },
  { id: 'bild.de', name: 'Bild' },
  { id: 'nytimes.com', name: 'NY Times' },
];

const STATIC_LIFECYCLE_DATA: TrendLifecycleMetric[] = [
  { time_slice: '2026-01-26 00:00:00+00', frequency: 2, velocity: 1 },
  { time_slice: '2026-01-23 00:00:00+00', frequency: 1, velocity: 0 },
  { time_slice: '2026-01-22 00:00:00+00', frequency: 1, velocity: -2 },
  { time_slice: '2026-01-21 00:00:00+00', frequency: 3, velocity: 2 },
  { time_slice: '2026-01-20 00:00:00+00', frequency: 1, velocity: 0 },
  { time_slice: '2026-01-19 00:00:00+00', frequency: 1, velocity: 0 },
  { time_slice: '2026-01-15 00:00:00+00', frequency: 1, velocity: 0 },
  { time_slice: '2026-01-14 00:00:00+00', frequency: 1, velocity: 0 },
  { time_slice: '2026-01-12 00:00:00+00', frequency: 1, velocity: 0 },
  { time_slice: '2026-01-09 00:00:00+00', frequency: 1, velocity: 0 },
  { time_slice: '2026-01-07 00:00:00+00', frequency: 1, velocity: 0 },
  { time_slice: '2026-01-06 00:00:00+00', frequency: 1, velocity: 1 },
];

export class TrendRepo {
  /**
   * Fetches trend metrics.
   * Currently returns static data from the dummy JSON.
   */
  static async getTrends(domain?: string): Promise<TrendMetric[]> {
    // Simulate a small network delay
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve(STATIC_TREND_DATA);
      }, 300);
    });
  }

  /**
   * Fetches context words (verbs/adjectives) for a specific trend topic.
   */
  static async getTrendContext(topic: string, domain?: string): Promise<TrendContextMetric[]> {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve(STATIC_CONTEXT_DATA);
      }, 300);
    });
  }

  /**
   * Fetches comparison data between two domains.
   */
  static async getTrendComparison(domainA: string, domainB: string, timeRange?: string): Promise<TrendComparisonMetric[]> {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve(STATIC_COMPARISON_DATA);
      }, 300);
    });
  }

  /**
   * Returns list of available domains for comparison.
   */
  static getAvailableDomains() {
    return AVAILABLE_DOMAINS;
  }

  /**
   * Fetches lifecycle metrics for a specific term.
   */
  static async getTrendLifecycle(term: string, domain: string, timeRange?: string): Promise<TrendLifecycleMetric[]> {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve(STATIC_LIFECYCLE_DATA);
      }, 300);
    });
  }
}