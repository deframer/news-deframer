SET duckdb.force_execution = true;

WITH config AS (
    SELECT
        90 AS days,
        5 AS trend_limit,
        ARRAY['circular-technology.com'] AS ignore_domains
),
article_stats AS (
    SELECT
        feeds.root_domain AS domain,
        items.language,
        TO_CHAR(MIN(items.pub_date), 'Mon DD, YYYY') || ' - ' || TO_CHAR(MAX(items.pub_date), 'Mon DD, YYYY') AS time,
        COUNT(items.id) AS article_count
    FROM items
    JOIN feeds ON feeds.id = items.feed_id
    CROSS JOIN config
    WHERE items.pub_date >= NOW() - (config.days * INTERVAL '1 DAY')
      AND feeds.enabled = true
      AND feeds.deleted_at IS NULL
      AND feeds.root_domain IS NOT NULL
      AND feeds.root_domain <> ''
      AND NOT (feeds.root_domain = ANY(config.ignore_domains))
    GROUP BY feeds.root_domain, items.language
),
ranked_trends AS (
    SELECT
        root_domain AS domain,
        "language",
        stem AS trend_topic,
        frequency,
        utility,
        outlier_ratio,
        time_slice,
        ROW_NUMBER() OVER (
            PARTITION BY root_domain, "language", stem
            ORDER BY outlier_ratio DESC, frequency DESC, time_slice DESC
        ) AS rn
    FROM view_trend_metrics_by_domain
    CROSS JOIN config
    WHERE stem_type = 'NOUN'
      AND time_slice >= NOW() - (config.days * INTERVAL '1 DAY')
      AND root_domain IS NOT NULL
      AND root_domain <> ''
      AND NOT (root_domain = ANY(config.ignore_domains))
      AND EXISTS (
          SELECT 1
          FROM feeds
          WHERE feeds.root_domain = view_trend_metrics_by_domain.root_domain
            AND feeds.enabled = true
            AND feeds.deleted_at IS NULL
      )
      AND utility >= 1
      AND outlier_ratio > 1.5
),
top_trends AS (
    SELECT
        domain,
        "language",
        trend_topic,
        frequency,
        outlier_ratio,
        ROW_NUMBER() OVER (
            PARTITION BY domain, "language"
            ORDER BY outlier_ratio DESC, frequency DESC, trend_topic
        ) AS domain_rank
    FROM ranked_trends
    WHERE rn = 1
),
trend_agg AS (
    SELECT
        domain,
        "language",
        STRING_AGG(trend_topic || ' (' || frequency || ')', ', ' ORDER BY domain_rank) AS "trends-top-5"
    FROM top_trends
    CROSS JOIN config
    WHERE domain_rank <= config.trend_limit
    GROUP BY domain, "language"
)
SELECT
    article_stats.domain,
    article_stats.language,
    article_stats.time,
    article_stats.article_count,
    trend_agg."trends-top-5"
FROM article_stats
LEFT JOIN trend_agg
    ON trend_agg.domain = article_stats.domain
   AND trend_agg."language" = article_stats.language
ORDER BY
    article_stats.domain,
    article_stats.language;
