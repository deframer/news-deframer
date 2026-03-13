-- Force DuckDB execution
SET duckdb.force_execution = true;


WITH constants AS (
    SELECT
        'spiegel.de' AS domain_a,
        'tagesschau.de' AS domain_b,
        INTERVAL '30 DAYS' AS lookback_window
),

/*
   STEP 1: Prepare Domain A (e.g., Spiegel)
   - Filter for NOUNS (Triggers)
   - Deduplicate: Keep only the peak day for each word
*/
domain_a_unique AS (
    SELECT
        stem,
        outlier_ratio,
        utility
    FROM (
        SELECT
            t.stem,
            t.outlier_ratio,
            t.utility,
            ROW_NUMBER() OVER (
                PARTITION BY t.stem
                ORDER BY t.outlier_ratio DESC, t.utility DESC, t.stem ASC
            ) AS rn
        FROM view_trend_metrics_by_domain t
        CROSS JOIN constants c
        WHERE t.root_domain = c.domain_a
          AND t."language" = 'de'
          AND t.stem_type = 'NOUN'
          AND t.time_slice >= NOW() - c.lookback_window
          AND t.utility >= 1
          AND t.outlier_ratio > 1.5
    )
    WHERE rn = 1
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
            t.stem,
            t.outlier_ratio,
            t.utility,
            ROW_NUMBER() OVER (
                PARTITION BY t.stem
                ORDER BY t.outlier_ratio DESC, t.utility DESC, t.stem ASC
            ) AS rn
        FROM view_trend_metrics_by_domain t
        CROSS JOIN constants c
        WHERE t.root_domain = c.domain_b
          AND t."language" = 'de'
          AND t.stem_type = 'NOUN'
          AND t.time_slice >= NOW() - c.lookback_window
          AND t.utility >= 1
          AND t.outlier_ratio > 1.5
    )
    WHERE rn = 1
),

/*
   STEP 3: Combine & Classify
   - INTERSECT: Both are talking about it
   - BLINDSPOT: Only one is talking about it
*/
all_joined AS (
    SELECT
        COALESCE(a.stem, b.stem) AS trend_topic,
        COALESCE(a.outlier_ratio, 0) AS score_a,
        COALESCE(b.outlier_ratio, 0) AS score_b,
        CASE
            WHEN a.stem IS NOT NULL AND b.stem IS NOT NULL THEN 'INTERSECT'
            WHEN a.stem IS NOT NULL AND b.stem IS NULL THEN 'BLINDSPOT_A'
            WHEN a.stem IS NULL AND b.stem IS NOT NULL THEN 'BLINDSPOT_B'
        END AS classification
    FROM domain_a_unique a
    FULL OUTER JOIN domain_b_unique b ON a.stem = b.stem
),

/*
   STEP 4: Rank within Groups
   - Shows the top 5 shared topics, top 5 unique to A, top 5 unique to B
*/
ranked_trends AS (
    SELECT
        *,
        CASE
            WHEN classification = 'INTERSECT' THEN (score_a + score_b)
            WHEN classification = 'BLINDSPOT_A' THEN score_a
            WHEN classification = 'BLINDSPOT_B' THEN score_b
        END AS ranking_score,
        ROW_NUMBER() OVER (
            PARTITION BY classification
            ORDER BY
                CASE
                    WHEN classification = 'INTERSECT' THEN (score_a + score_b)
                    WHEN classification = 'BLINDSPOT_A' THEN score_a
                    WHEN classification = 'BLINDSPOT_B' THEN score_b
                END DESC,
                trend_topic ASC
        ) AS rank_group
    FROM all_joined
)

SELECT
    rank_group AS rank,
    trend_topic,
    CASE
        WHEN classification IN ('INTERSECT', 'BLINDSPOT_A') THEN 'spiegel.de'
        ELSE ''
    END AS domain_a,
    CASE
        WHEN classification IN ('INTERSECT', 'BLINDSPOT_B') THEN 'tagesschau.de'
        ELSE ''
    END AS domain_b,
    ROUND(score_a::numeric, 2) AS score_a,
    ROUND(score_b::numeric, 2) AS score_b
FROM ranked_trends
WHERE rank_group <= 5
ORDER BY classification, rank_group;
