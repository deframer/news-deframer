import log from '../../shared/logger';

export interface TrendComparisonMetric {
  classification: 'BLINDSPOT_A' | 'BLINDSPOT_B' | 'INTERSECT';
  rank_group: number;
  trend_topic: string;
  score_a: number;
  score_b: number;
}

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

export class TrendRepo {
  /**
   * Fetches comparison data between two domains.
   */
  static async getTrendComparison(domainA: string, domainB: string, days?: number): Promise<TrendComparisonMetric[]> {
    log.debug(`TrendRepo.getTrendComparison called with domainA=${domainA}, domainB=${domainB}, days=${days}`);
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve(STATIC_COMPARISON_DATA);
      }, 300);
    });
  }
}