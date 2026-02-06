CREATE OR REPLACE VIEW view_trend_metrics_by_domain AS
WITH raw_unrolled AS (
    -- 1. Unnest stems and include ROOT_DOMAIN
    SELECT item_id, feed_id, pub_date, "language", root_domain, unnest(noun_stems) as stem, 'NOUN' as stem_type
    FROM public.trends
    UNION ALL
    SELECT item_id, feed_id, pub_date, "language", root_domain, unnest(verb_stems) as stem, 'VERB' as stem_type
    FROM public.trends
    UNION ALL
    SELECT item_id, feed_id, pub_date, "language", root_domain, unnest(adjective_stems) as stem, 'ADJ' as stem_type
    FROM public.trends
),
domain_stats AS (
    -- 2. Aggregate per Stem + Time + Language + DOMAIN
    SELECT
        stem,
        stem_type,
        "language",
        root_domain,  -- <--- Grouping by Domain here
        date_trunc('day', pub_date) as time_slice,
        count(*) as frequency,
        count(distinct feed_id) as utility
    FROM raw_unrolled
    GROUP BY 1, 2, 3, 4, 5
)
SELECT
    stem,
    stem_type,
    "language",
    root_domain,
    time_slice,
    frequency,
    utility,
    -- 3. Calculate metrics relative to the SPECIFIC DOMAIN history
    -- We add root_domain to the PARTITION BY clause to isolate the math
    frequency - LAG(frequency, 1, 0) OVER (
        PARTITION BY stem, stem_type, "language", root_domain
        ORDER BY time_slice
    ) AS growth_delta,

    frequency / NULLIF(
        AVG(frequency) OVER (
            PARTITION BY stem, stem_type, "language", root_domain
            ORDER BY time_slice
            ROWS BETWEEN 7 PRECEDING AND 1 PRECEDING
        ), 0
    ) AS outlier_ratio
FROM domain_stats;
