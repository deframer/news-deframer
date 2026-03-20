WITH config AS (
    SELECT
        90 AS days,
        ARRAY['circular-technology.com'] AS ignore_domains
)
SELECT
    feeds.root_domain AS domain,
    items.language,
    TO_CHAR(MIN(items.pub_date), 'Mon DD, YYYY') || ' - ' || TO_CHAR(MAX(items.pub_date), 'Mon DD, YYYY') AS time,
    COUNT(items.id) AS article_count,
    ROUND(AVG(CASE WHEN (think_result ->> 'framing')::numeric BETWEEN 0 AND 1 THEN (think_result ->> 'framing')::numeric END), 2) AS framing,
    ROUND(AVG(CASE WHEN (think_result ->> 'clickbait')::numeric BETWEEN 0 AND 1 THEN (think_result ->> 'clickbait')::numeric END), 2) AS clickbait,
    ROUND(AVG(CASE WHEN (think_result ->> 'persuasive')::numeric BETWEEN 0 AND 1 THEN (think_result ->> 'persuasive')::numeric END), 2) AS persuasive,
    ROUND(AVG(CASE WHEN (think_result ->> 'hyper_stimulus')::numeric BETWEEN 0 AND 1 THEN (think_result ->> 'hyper_stimulus')::numeric END), 2) AS hyper_stimulus,
    ROUND(AVG(CASE WHEN (think_result ->> 'speculative')::numeric BETWEEN 0 AND 1 THEN (think_result ->> 'speculative')::numeric END), 2) AS speculative,
    ROUND(AVG(CASE WHEN (think_result ->> 'overall')::numeric BETWEEN 0 AND 1 THEN (think_result ->> 'overall')::numeric END), 2) AS overall
FROM
    items
JOIN
    feeds ON feeds.id = items.feed_id
CROSS JOIN
    config
WHERE
    think_result IS NOT NULL
    AND pub_date >= NOW() - (config.days * INTERVAL '1 DAY')
    AND feeds.enabled = true
    AND feeds.deleted_at IS NULL
    AND feeds.root_domain IS NOT NULL
    AND feeds.root_domain <> ''
    AND NOT (feeds.root_domain = ANY(config.ignore_domains))
GROUP BY
    feeds.root_domain,
    items.language
ORDER BY
    overall DESC,
    feeds.root_domain,
    items.language;
