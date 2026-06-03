WITH config AS (
  SELECT
    365 AS days,
    5 AS min_articles,
    'de' AS language,
    'de' AS country
),
expanded_authors AS (
  SELECT
    feeds.root_domain AS root_domain,
    items.think_rating,
    items.pub_date,
    UNNEST(items.authors) AS author
  FROM items
  JOIN feeds ON feeds.id = items.feed_id
  CROSS JOIN config
  WHERE feeds.language = config.language
    AND feeds.country = config.country
    AND feeds.enabled = TRUE
    AND feeds.deleted_at IS NULL
    AND items.pub_date >= NOW() - (config.days * INTERVAL '1 DAY')
)
SELECT
  root_domain,
  author,
  MIN(pub_date) AS first_publication_date,
  COUNT(*) AS article_count
FROM expanded_authors
WHERE root_domain IS NOT NULL
  AND root_domain <> ''
  AND author IS NOT NULL
  AND author <> ''
GROUP BY root_domain, author
HAVING COUNT(*) >= (SELECT min_articles FROM config)
ORDER BY root_domain ASC, author ASC;
