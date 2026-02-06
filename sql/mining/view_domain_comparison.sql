-- Force DuckDB execution
SET duckdb.force_execution = true;

WITH domain_a AS (
    SELECT stem, outlier_ratio
    FROM view_trend_metrics_by_domain
    WHERE root_domain = 'spiegel.de'       -- Source A
      AND "language" = 'de'             -- Mandatory Language Filter
      AND time_slice = DATE_TRUNC('day', NOW())
      AND outlier_ratio > 1.5           -- Threshold for "Trending"
),
domain_b AS (
    SELECT stem, outlier_ratio
    FROM view_trend_metrics_by_domain
    WHERE root_domain = 'tagesschau.de'   -- Source B
      AND "language" = 'de'
      AND time_slice = DATE_TRUNC('day', NOW())
      AND outlier_ratio > 1.5
)
SELECT
    -- Get the word from whichever side has it
    COALESCE(a.stem, b.stem) as trend_topic,

    -- Show the scores for comparison
    ROUND(a.outlier_ratio::numeric, 2) as score_A,
    ROUND(b.outlier_ratio::numeric, 2) as score_B,

    -- Classify the Trend
    CASE
        WHEN a.stem IS NOT NULL AND b.stem IS NOT NULL
            THEN 'INTERSECT (Both)'
        WHEN a.stem IS NOT NULL AND b.stem IS NULL
            THEN 'BLINDSPOT (Only A)'
        WHEN a.stem IS NULL AND b.stem IS NOT NULL
            THEN 'BLINDSPOT (Only B)'
    END as classification

FROM domain_a a
FULL OUTER JOIN domain_b b ON a.stem = b.stem
ORDER BY classification, score_A DESC;