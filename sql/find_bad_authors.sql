WITH expanded_authors AS (
  SELECT
    feeds.root_domain AS domain,
    items.think_rating,
    UNNEST(items.authors) AS author
  FROM items
  JOIN feeds ON feeds.id = items.feed_id
)
SELECT
  domain,
  author,
  COUNT(*) AS article_count,
  ROUND(AVG(think_rating), 2) AS avg,
  PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY think_rating) AS median,
  ROUND(STDDEV_SAMP(think_rating), 2) AS stddev
FROM expanded_authors
WHERE domain IS NOT NULL
  AND domain <> ''
  AND author IS NOT NULL
  AND author <> ''
GROUP BY domain, author
ORDER BY domain ASC, avg DESC, author ASC;
