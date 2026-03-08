WITH expanded_authors AS (
  SELECT
    feeds.root_domain AS domain,
    items.think_rating,
    UNNEST(items.authors) AS author
  FROM items
  JOIN feeds ON feeds.id = items.feed_id
), ranked_authors AS (
  SELECT
    COUNT(*) AS article_count,
    ROUND(AVG(think_rating), 2) AS avg,
    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY think_rating) AS median,
    ROUND(STDDEV_SAMP(think_rating), 2) AS stddev,
    domain,
    author,
    ROW_NUMBER() OVER (
      PARTITION BY domain
      ORDER BY COUNT(*) DESC, AVG(think_rating) DESC, author ASC
    ) AS rank_in_domain
  FROM expanded_authors
  WHERE domain IS NOT NULL
    AND domain <> ''
    AND author IS NOT NULL
    AND author <> ''
  GROUP BY domain, author
)
SELECT
  domain,
  author,
  article_count,
  avg,
  median,
  stddev
FROM ranked_authors
WHERE rank_in_domain <= 5
ORDER BY domain ASC, article_count DESC, avg DESC, author ASC;
