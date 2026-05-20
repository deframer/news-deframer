SET duckdb.force_execution = true;

WITH config AS (
  SELECT
    7 AS days,
    20 AS item_limit,
    ARRAY['public_service_media'] AS tags
), ranked_items AS (
  SELECT
    feeds.root_domain AS domain,
    items.think_rating AS rating,
    NULLIF(items.think_result ->> 'category', '') AS category,
    items.pub_date,
    items.url AS article_url,
    ROW_NUMBER() OVER (
      PARTITION BY feeds.root_domain
      ORDER BY items.think_rating DESC, items.id ASC
    ) AS rn
  FROM items
  JOIN feeds ON feeds.id = items.feed_id
  CROSS JOIN config
  WHERE items.pub_date >= NOW() - (config.days * INTERVAL '1 DAY')
    AND feeds.tags @> config.tags
    AND feeds.enabled = true
    AND feeds.deleted_at IS NULL
    AND feeds.root_domain IS NOT NULL
    AND feeds.root_domain <> ''
), overall_stats AS (
  SELECT
    ROUND(AVG(rating), 2) AS avg,
    ROUND((PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY rating))::numeric, 2) AS median
  FROM ranked_items
)
SELECT
  domain,
  rating,
  overall_stats.avg,
  overall_stats.median,
  category,
  article_url
FROM ranked_items
CROSS JOIN overall_stats
WHERE rn <= (SELECT item_limit FROM config)
ORDER BY domain ASC, rating DESC, article_url ASC;
