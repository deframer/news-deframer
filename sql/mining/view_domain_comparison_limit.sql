-- Force DuckDB execution
SET duckdb.force_execution = true;

/*
   STEP 1: Prepare Domain A (e.g., Spiegel)
   - Filter for NOUNS (Triggers)
   - Deduplicate: Keep only the peak day for each word
*/
WITH domain_a_unique AS (
    SELECT
        stem,
        outlier_ratio,
        utility
    FROM (
        SELECT
            stem,
            outlier_ratio,
            utility,
            ROW_NUMBER() OVER (PARTITION BY stem ORDER BY outlier_ratio DESC) as rn
        FROM view_trend_metrics_by_domain
        WHERE root_domain = 'spiegel.de'
          AND "language" = 'de'
          AND stem_type = 'NOUN'          -- <--- Thesis: Focus on Triggers/Topics [1][2]
          AND time_slice >= NOW() - INTERVAL '7 DAYS'
          AND utility >= 1
          AND outlier_ratio > 1.5
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
            ROW_NUMBER() OVER (PARTITION BY stem ORDER BY outlier_ratio DESC) as rn
        FROM view_trend_metrics_by_domain
        WHERE root_domain = 'tagesschau.de'
          AND "language" = 'de'
          AND stem_type = 'NOUN'
          AND time_slice >= NOW() - INTERVAL '7 DAYS'
          AND utility >= 1
          AND outlier_ratio > 1.5
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
        COALESCE(a.stem, b.stem) as trend_topic,
        COALESCE(a.outlier_ratio, 0) as score_a,
        COALESCE(b.outlier_ratio, 0) as score_b,
        CASE
            WHEN a.stem IS NOT NULL AND b.stem IS NOT NULL THEN 'INTERSECT'
            WHEN a.stem IS NOT NULL AND b.stem IS NULL THEN 'BLINDSPOT_A' -- Unique to A
            WHEN a.stem IS NULL AND b.stem IS NOT NULL THEN 'BLINDSPOT_B' -- Unique to B
        END as classification
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
        ROW_NUMBER() OVER (
            PARTITION BY classification
            ORDER BY
                -- Dynamic ordering: Prioritize the relevant score for the group
                CASE
                    WHEN classification = 'INTERSECT' THEN (score_a + score_b) -- Combined heat
                    WHEN classification = 'BLINDSPOT_A' THEN score_a           -- A's heat
                    WHEN classification = 'BLINDSPOT_B' THEN score_b           -- B's heat
                END DESC
        ) as rank_group
    FROM all_joined
)

SELECT
    classification,
    rank_group,
    trend_topic,
    ROUND(score_a::numeric, 2) as score_a,
    ROUND(score_b::numeric, 2) as score_b
FROM ranked_trends
WHERE rank_group <= 5
ORDER BY classification, rank_group;
