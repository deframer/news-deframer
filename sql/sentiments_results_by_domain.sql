WITH config AS (
    SELECT
        90 AS days,
        10 AS min_articles,
        'de' AS language,
        ARRAY['trump', 'ukraine', 'epstein'] AS terms,
        ARRAY['circular-technology.com'] AS ignore_domains
),
config_terms AS (
    SELECT
        config.days,
        config.min_articles,
        config.language,
        config.ignore_domains,
        UNNEST(config.terms) AS trend
    FROM config
)
SELECT
    config_terms.trend,
    t.root_domain AS domain,
    t."language" AS language,
    TO_CHAR(MIN(t.pub_date), 'Mon DD, YYYY') || ' - ' || TO_CHAR(MAX(t.pub_date), 'Mon DD, YYYY') AS time,
    COUNT(t.item_id) AS article_count,
    ROUND(AVG(CASE WHEN (t.sentiments->>'v')::numeric BETWEEN 1 AND 9 THEN (t.sentiments->>'v')::numeric END), 2) AS valence,
    ROUND(AVG(CASE WHEN (t.sentiments->>'a')::numeric BETWEEN 1 AND 9 THEN (t.sentiments->>'a')::numeric END), 2) AS arousal,
    ROUND(AVG(CASE WHEN (t.sentiments->>'d')::numeric BETWEEN 1 AND 9 THEN (t.sentiments->>'d')::numeric END), 2) AS dominance,
    ROUND(AVG(CASE WHEN (t.sentiments->>'j')::numeric BETWEEN 1 AND 5 THEN (t.sentiments->>'j')::numeric END), 2) AS joy,
    ROUND(AVG(CASE WHEN (t.sentiments->>'a_n')::numeric BETWEEN 1 AND 5 THEN (t.sentiments->>'a_n')::numeric END), 2) AS anger,
    ROUND(AVG(CASE WHEN (t.sentiments->>'s')::numeric BETWEEN 1 AND 5 THEN (t.sentiments->>'s')::numeric END), 2) AS sadness,
    ROUND(AVG(CASE WHEN (t.sentiments->>'f')::numeric BETWEEN 1 AND 5 THEN (t.sentiments->>'f')::numeric END), 2) AS fear,
    ROUND(AVG(CASE WHEN (t.sentiments->>'d_g')::numeric BETWEEN 1 AND 5 THEN (t.sentiments->>'d_g')::numeric END), 2) AS disgust
FROM public.trends t
JOIN feeds ON feeds.id = t.feed_id
JOIN config_terms ON config_terms.trend = ANY(t.noun_stems)
WHERE t.pub_date >= NOW() - (config_terms.days * INTERVAL '1 DAY')
  AND t."language" = config_terms.language
  AND t.sentiments IS NOT NULL
  AND t.sentiments <> '{}'::jsonb
  AND feeds.enabled = true
  AND feeds.deleted_at IS NULL
  AND t.root_domain IS NOT NULL
  AND t.root_domain <> ''
  AND NOT (t.root_domain = ANY(config_terms.ignore_domains))
GROUP BY
    config_terms.trend,
    t.root_domain,
    t."language"
HAVING
    COUNT(t.item_id) >= MAX(config_terms.min_articles)
ORDER BY
    config_terms.trend,
    t.root_domain,
    t."language";
