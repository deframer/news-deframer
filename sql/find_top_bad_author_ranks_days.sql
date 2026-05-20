SET duckdb.force_execution = true;

WITH config AS (
  SELECT
    90 AS days,
    5 AS author_limit
),
expanded_authors AS (
  SELECT
    feeds.root_domain AS domain,
    items.pub_date,
    items.think_rating,
    items.url AS article_url,
    UNNEST(
      CASE
        WHEN items.authors IS NULL OR COALESCE(cardinality(items.authors), 0) = 0 THEN ARRAY['unknown']
        ELSE items.authors
      END
    ) AS author
  FROM items
  JOIN feeds ON feeds.id = items.feed_id
  CROSS JOIN config
  WHERE items.pub_date >= NOW() - (config.days * INTERVAL '1 DAY')
),
normalized_authors AS (
  SELECT
    domain,
    pub_date,
    think_rating,
    article_url,
    COALESCE(NULLIF(TRIM(author), ''), 'unknown') AS author
  FROM expanded_authors
),
author_url_choice AS (
  SELECT
    domain,
    author,
    article_url,
    ROW_NUMBER() OVER (
      PARTITION BY domain, author
      ORDER BY think_rating DESC, random()
    ) AS rn
  FROM normalized_authors
  WHERE domain IS NOT NULL
    AND domain <> ''
    AND article_url IS NOT NULL
    AND article_url <> ''
),
ranked_named_authors AS (
  SELECT
    domain,
    author,
    TO_CHAR(MIN(pub_date), 'Mon DD, YYYY') AS start_date,
    COUNT(*) AS article_count,
    ROUND(AVG(think_rating), 2) AS avg,
    ROUND((PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY think_rating))::numeric, 2) AS median,
    ROUND((STDDEV_SAMP(think_rating))::numeric, 2) AS stddev,
    ROW_NUMBER() OVER (
      PARTITION BY domain
      ORDER BY AVG(think_rating) DESC, author ASC
    ) AS rn
  FROM normalized_authors
  WHERE domain IS NOT NULL
    AND domain <> ''
    AND author <> 'unknown'
  GROUP BY domain, author
),
unknown_authors AS (
  SELECT
    domain,
    'unknown' AS author,
    TO_CHAR(MIN(pub_date), 'Mon DD, YYYY') AS start_date,
    COUNT(*) AS article_count,
    ROUND(AVG(think_rating), 2) AS avg,
    ROUND((PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY think_rating))::numeric, 2) AS median,
    ROUND((STDDEV_SAMP(think_rating))::numeric, 2) AS stddev,
    1 AS rn
  FROM normalized_authors
  WHERE domain IS NOT NULL
    AND domain <> ''
    AND author = 'unknown'
  GROUP BY domain
),
selected_article_urls AS (
  SELECT
    domain,
    author,
    article_url
  FROM author_url_choice
  WHERE rn = 1
),
final_results AS (
  SELECT
    domain,
    author,
    start_date,
    article_count,
    avg,
    median,
    stddev,
    article_url,
    0 AS sort_key
  FROM ranked_named_authors
  CROSS JOIN config
  LEFT JOIN selected_article_urls USING (domain, author)
  WHERE rn <= config.author_limit
  UNION ALL
  SELECT
    domain,
    author,
    start_date,
    article_count,
    avg,
    median,
    stddev,
    article_url,
    1 AS sort_key
  FROM unknown_authors
  LEFT JOIN selected_article_urls USING (domain, author)
)
SELECT
  domain,
  author,
  start_date,
  article_count,
  avg,
  median,
  stddev,
  article_url
FROM final_results
ORDER BY domain ASC, sort_key ASC, avg DESC, author ASC;
