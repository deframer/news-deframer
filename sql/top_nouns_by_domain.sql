WITH noun_counts AS (
    SELECT
        t.root_domain,
        stem,
        COUNT(*) AS frequency
    FROM public.trends t
    JOIN public.feeds f ON f.id = t.feed_id
    CROSS JOIN LATERAL unnest(t.noun_stems) AS stem
    WHERE f.enabled = true
      AND f.deleted_at IS NULL
      AND t.root_domain IS NOT NULL
      AND t.root_domain <> ''
      AND stem IS NOT NULL
      AND stem <> ''
    GROUP BY t.root_domain, stem
), ranked AS (
    SELECT
        root_domain,
        stem,
        frequency,
        ROW_NUMBER() OVER (
            PARTITION BY root_domain
            ORDER BY frequency DESC, stem ASC
        ) AS rn
    FROM noun_counts
)
SELECT
    root_domain,
    stem,
    frequency
FROM ranked
WHERE rn <= 20
ORDER BY root_domain ASC, frequency DESC, stem ASC;
