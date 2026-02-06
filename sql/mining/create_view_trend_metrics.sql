CREATE OR REPLACE VIEW view_trend_metrics AS
WITH raw_unrolled AS (
    -- 1. Unnest stems and include LANGUAGE
    SELECT item_id, feed_id, pub_date, "language", unnest(noun_stems) as stem, 'NOUN' as stem_type
    FROM public.trends
    UNION ALL
    SELECT item_id, feed_id, pub_date, "language", unnest(verb_stems) as stem, 'VERB' as stem_type
    FROM public.trends
    UNION ALL
    SELECT item_id, feed_id, pub_date, "language", unnest(adjective_stems) as stem, 'ADJ' as stem_type
    FROM public.trends
),
daily_stats AS (
    -- 2. Aggregate per Stem + Time + LANGUAGE
    SELECT
        stem,
        stem_type,
        "language",
        date_trunc('day', pub_date) as time_slice,
        count(*) as frequency,
        count(distinct feed_id) as utility
    FROM raw_unrolled
    GROUP BY 1, 2, 3, 4
)
SELECT
    stem,
    stem_type,
    "language",
    time_slice,
    frequency,
    utility,
    -- 3. Calculate metrics relative to the specific LANGUAGE history
    frequency - LAG(frequency, 1, 0) OVER (
        PARTITION BY stem, stem_type, "language"
        ORDER BY time_slice
    ) AS growth_delta,

    frequency / NULLIF(
        AVG(frequency) OVER (
            PARTITION BY stem, stem_type, "language"
            ORDER BY time_slice
            ROWS BETWEEN 7 PRECEDING AND 1 PRECEDING
        ), 0
    ) AS outlier_ratio
FROM daily_stats;
