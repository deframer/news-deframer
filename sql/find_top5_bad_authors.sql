WITH expanded_authors AS (
  SELECT
    feeds.root_domain AS domain,
    items.think_rating,
    UNNEST(items.authors) AS author
  FROM items
  JOIN feeds ON feeds.id = items.feed_id
), ranked_authors AS (
  SELECT
    AVG(think_rating) AS rating,
    COUNT(*) AS article_count,
    domain,
    author,
    ROW_NUMBER() OVER (
      PARTITION BY domain
      ORDER BY AVG(think_rating) DESC, author ASC
    ) AS rank_in_domain
  FROM expanded_authors
  WHERE domain IS NOT NULL
    AND domain <> ''
    AND author IS NOT NULL
    AND author <> ''
  GROUP BY domain, author
)
SELECT
  rating,
  article_count,
  domain,
  author
FROM ranked_authors
WHERE rank_in_domain <= 5
ORDER BY domain ASC, rating DESC, author ASC;
