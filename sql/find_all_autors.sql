WITH expanded_authors AS (
  SELECT
    feeds.root_domain AS domain,
    items.think_rating,
    UNNEST(items.authors) AS author
  FROM items
  JOIN feeds ON feeds.id = items.feed_id
  WHERE feeds.language = 'en'
    -- AND feeds.country = 'de'
    AND feeds.enabled = TRUE
    AND feeds.deleted_at IS NULL
)
SELECT
  domain,
  author,
  COUNT(*) AS article_count
FROM expanded_authors
WHERE domain IS NOT NULL
  AND domain <> ''
  AND author IS NOT NULL
  AND author <> ''
GROUP BY domain, author
HAVING COUNT(*) > 5
ORDER BY domain ASC, author ASC;
