-- Force DuckDB execution
SET duckdb.force_execution = true;

SET SESSION "vars.domain_a" = 'nius.de';
SET SESSION "vars.domain_b" = 'tagesschau.de';
SET SESSION "vars.language" = 'de';
SET SESSION "vars.last" = '30 DAYS';

/*
   STEP 1: Prepare Domain A (e.g., Spiegel)
   - Filter for NOUNS (Triggers)
   - Deduplicate: Keep only the peak day for each word
*/
WITH params AS (
    SELECT
        current_setting('vars.domain_a')::VARCHAR AS domain_a,
        current_setting('vars.domain_b')::VARCHAR AS domain_b,
        current_setting('vars.language')::VARCHAR AS language,
        current_setting('vars.last')::INTERVAL AS lookback_window
),
domain_a_unique AS (
    SELECT
        stem,
        outlier_ratio,
        utility
    FROM (
        SELECT
            stem,
            outlier_ratio,
            utility,
            ROW_NUMBER() OVER (
                PARTITION BY stem
                ORDER BY outlier_ratio DESC, utility DESC, stem ASC
            ) as rn
        FROM view_trend_metrics_by_domain v
        CROSS JOIN params p
        WHERE v.root_domain = p.domain_a
          AND v."language" = p.language
          AND v.stem_type = 'NOUN'          -- <--- Thesis: Focus on Triggers/Topics [1][2]
          AND v.time_slice >= NOW() - p.lookback_window
          AND v.utility >= 1
          AND v.outlier_ratio > 1.5
    )
    WHERE rn = 1 -- Keep only the strongest signal per word
),

/*
   STEP 2: Prepare Domain B (e.g., Tagesschau)
   - Same cleaning logic
*/
domain_b_unique AS (
    SELECT
        stem,
        outlier_ratio,
        utility
    FROM (
        SELECT
            stem,
            outlier_ratio,
            utility,
            ROW_NUMBER() OVER (
                PARTITION BY stem
                ORDER BY outlier_ratio DESC, utility DESC, stem ASC
            ) as rn
        FROM view_trend_metrics_by_domain v
        CROSS JOIN params p
        WHERE v.root_domain = p.domain_b
          AND v."language" = p.language
          AND v.stem_type = 'NOUN'
          AND v.time_slice >= NOW() - p.lookback_window
          AND v.utility >= 1
          AND v.outlier_ratio > 1.5
    )
    WHERE rn = 1
),

/*
   STEP 3: Combine & Classify
   - Keep an internal category for ranking
   - Expose domain names as separate output columns
*/
all_joined AS (
    SELECT
        COALESCE(a.stem, b.stem) as trend_topic,
        COALESCE(a.outlier_ratio, 0) as score_a,
        COALESCE(b.outlier_ratio, 0) as score_b,
        CASE
            WHEN a.stem IS NOT NULL AND b.stem IS NOT NULL THEN 'shared'
            WHEN a.stem IS NOT NULL AND b.stem IS NULL THEN 'domain_a_only'
            WHEN a.stem IS NULL AND b.stem IS NOT NULL THEN 'domain_b_only'
        END as category_key
    FROM domain_a_unique a
    FULL OUTER JOIN domain_b_unique b ON a.stem = b.stem
),

/*
   STEP 4: Rank within Groups
   - Shows the top 5 shared topics and top 5 unique topics per domain
*/
ranked_trends AS (
    SELECT
        *,
        CASE
            WHEN category_key = 'shared' THEN (score_a + score_b)
            WHEN category_key = 'domain_a_only' THEN score_a
            WHEN category_key = 'domain_b_only' THEN score_b
        END as ranking_score,
        ROW_NUMBER() OVER (
            PARTITION BY category_key
            ORDER BY
                -- Dynamic ordering: Prioritize the relevant score for the group
                CASE
                    WHEN category_key = 'shared' THEN (score_a + score_b) -- Combined heat
                    WHEN category_key = 'domain_a_only' THEN score_a
                    WHEN category_key = 'domain_b_only' THEN score_b
                END DESC,
                trend_topic ASC
        ) as rank_group
    FROM all_joined
)

SELECT
    CASE
        WHEN r.category_key IN ('shared', 'domain_a_only') THEN p.domain_a
        ELSE NULL
    END as domain_a,
    CASE
        WHEN r.category_key IN ('shared', 'domain_b_only') THEN p.domain_b
        ELSE NULL
    END as domain_b,
    rank_group,
    trend_topic,
    ROUND(score_a::numeric, 2) as score_a,
    ROUND(score_b::numeric, 2) as score_b
FROM ranked_trends r
CROSS JOIN params p
WHERE rank_group <= 5
ORDER BY r.category_key, rank_group;
