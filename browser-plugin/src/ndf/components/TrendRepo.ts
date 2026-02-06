import log from '../../shared/logger';
import { getSettings } from '../../shared/settings';
import { NewsDeframerClient } from '../client';

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
   * Fetches context words (verbs/adjectives) for a specific trend topic.
   */
  static async getTrendContext(topic: string, domain?: string, days?: number): Promise<TrendContextMetric[]> {
    log.debug(`TrendRepo.getTrendContext called with topic=${topic}, domain=${domain}, days=${days}`);
    if (!domain) return [];

    const settings = await getSettings();
    const client = new NewsDeframerClient(settings);

    // Resolve language for the domain
    const domains = await client.getDomains();
    const domainEntry = domains.find((d) => d.domain === domain);
    const lang = domainEntry?.language || 'de';

    const data = await client.getContextByDomain(topic, domain, lang, days || 7);
    return data.map((d) => ({
      context_word: d.context,
      frequency: d.frequency,
      type: 'WORD', // API does not return type currently
    }));
  }

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

  /**
   * Returns list of available domains for comparison.
   */
  static getAvailableDomains(currentDomain?: string) {
    log.debug(`TrendRepo.getAvailableDomains called with currentDomain=${currentDomain}`);
    return AVAILABLE_DOMAINS;
  }

  /**
   * Fetches lifecycle metrics for a specific term.
   */
  static async getTrendLifecycle(term: string, domain: string, days?: number): Promise<TrendLifecycleMetric[]> {
    log.debug(`TrendRepo.getTrendLifecycle called with term=${term}, domain=${domain}, days=${days}`);
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve(STATIC_LIFECYCLE_DATA);
      }, 300);
    });
  }
}