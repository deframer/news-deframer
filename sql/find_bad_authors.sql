WITH expanded_authors AS (
  SELECT
    feeds.root_domain AS domain,
    items.think_rating,
    UNNEST(items.authors) AS author
  FROM items
  JOIN feeds ON feeds.id = items.feed_id
)
SELECT
  AVG(think_rating) AS rating,
  COUNT(*) AS article_count,
  domain,
  author
FROM expanded_authors
WHERE domain IS NOT NULL
  AND domain <> ''
  AND author IS NOT NULL
  AND author <> ''
GROUP BY domain, author
ORDER BY domain ASC, rating DESC, author ASC;
