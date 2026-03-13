WITH base AS (
    SELECT
        f.root_domain AS domain,
        COUNT(*) AS articles,
        ROUND(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY COALESCE((i.think_result->>'overall')::float, 0.0))::numeric, 2) AS overall,
        ROUND(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY COALESCE((i.think_result->>'framing')::float, 0.0))::numeric, 2) AS framing,
        ROUND(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY COALESCE((i.think_result->>'clickbait')::float, 0.0))::numeric, 2) AS clickbait,
        ROUND(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY COALESCE((i.think_result->>'persuasive')::float, 0.0))::numeric, 2) AS persuasive,
        ROUND(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY COALESCE((i.think_result->>'hyper_stimulus')::float, 0.0))::numeric, 2) AS hyper_stimulus,
        ROUND(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY COALESCE((i.think_result->>'speculative')::float, 0.0))::numeric, 2) AS speculative
    FROM items i
    JOIN feeds f ON f.id = i.feed_id
    WHERE
        i.think_result IS NOT NULL
        AND i.think_error IS NULL
        AND i.think_error_count = 0
        AND f.root_domain IS NOT NULL
        AND f.root_domain <> ''
    GROUP BY f.root_domain
), known_authors AS (
    SELECT
        f.root_domain AS domain,
        COUNT(DISTINCT NULLIF(BTRIM(author_name), '')) AS known_authors
    FROM items i
    JOIN feeds f ON f.id = i.feed_id
    CROSS JOIN LATERAL unnest(i.authors) AS author_name
    WHERE
        i.think_result IS NOT NULL
        AND i.think_error IS NULL
        AND i.think_error_count = 0
        AND f.root_domain IS NOT NULL
        AND f.root_domain <> ''
    GROUP BY f.root_domain
)
SELECT
    b.domain,
    b.articles,
    COALESCE(a.known_authors, 0) AS known_authors,
    b.overall,
    b.framing,
    b.clickbait,
    b.persuasive,
    b.hyper_stimulus,
    b.speculative
FROM base b
LEFT JOIN known_authors a ON a.domain = b.domain
-- WHERE b.domain in ( 'apollo-news.net', 'nius.de', 'tagesschau.de', 'faz.net', 'welt.de'    )
ORDER BY b.domain ASC;
