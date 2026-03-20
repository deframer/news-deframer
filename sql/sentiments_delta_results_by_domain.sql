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
),
domain_sentiments AS (
    SELECT
        t.root_domain AS domain,
        t."language" AS language,
        TO_CHAR(MIN(t.pub_date), 'Mon DD, YYYY') || ' - ' || TO_CHAR(MAX(t.pub_date), 'Mon DD, YYYY') AS time,
        COUNT(t.item_id) AS all_article_count,
        AVG(CASE WHEN (t.sentiments->>'v')::numeric BETWEEN 1 AND 9 THEN (t.sentiments->>'v')::numeric END) AS valence_avg,
        AVG(CASE WHEN (t.sentiments->>'a')::numeric BETWEEN 1 AND 9 THEN (t.sentiments->>'a')::numeric END) AS arousal_avg,
        AVG(CASE WHEN (t.sentiments->>'d')::numeric BETWEEN 1 AND 9 THEN (t.sentiments->>'d')::numeric END) AS dominance_avg,
        AVG(CASE WHEN (t.sentiments->>'j')::numeric BETWEEN 1 AND 5 THEN (t.sentiments->>'j')::numeric END) AS joy_avg,
        AVG(CASE WHEN (t.sentiments->>'a_n')::numeric BETWEEN 1 AND 5 THEN (t.sentiments->>'a_n')::numeric END) AS anger_avg,
        AVG(CASE WHEN (t.sentiments->>'s')::numeric BETWEEN 1 AND 5 THEN (t.sentiments->>'s')::numeric END) AS sadness_avg,
        AVG(CASE WHEN (t.sentiments->>'f')::numeric BETWEEN 1 AND 5 THEN (t.sentiments->>'f')::numeric END) AS fear_avg,
        AVG(CASE WHEN (t.sentiments->>'d_g')::numeric BETWEEN 1 AND 5 THEN (t.sentiments->>'d_g')::numeric END) AS disgust_avg
    FROM public.trends t
    JOIN feeds ON feeds.id = t.feed_id
    CROSS JOIN config
    WHERE t.pub_date >= NOW() - (config.days * INTERVAL '1 DAY')
      AND t."language" = config.language
      AND t.sentiments IS NOT NULL
      AND t.sentiments <> '{}'::jsonb
      AND feeds.enabled = true
      AND feeds.deleted_at IS NULL
      AND t.root_domain IS NOT NULL
      AND t.root_domain <> ''
      AND NOT (t.root_domain = ANY(config.ignore_domains))
    GROUP BY
        t.root_domain,
        t."language"
),
trend_sentiments AS (
    SELECT
        config_terms.trend,
        t.root_domain AS domain,
        t."language" AS language,
        COUNT(t.item_id) AS article_count,
        AVG(CASE WHEN (t.sentiments->>'v')::numeric BETWEEN 1 AND 9 THEN (t.sentiments->>'v')::numeric END) AS valence_avg,
        AVG(CASE WHEN (t.sentiments->>'a')::numeric BETWEEN 1 AND 9 THEN (t.sentiments->>'a')::numeric END) AS arousal_avg,
        AVG(CASE WHEN (t.sentiments->>'d')::numeric BETWEEN 1 AND 9 THEN (t.sentiments->>'d')::numeric END) AS dominance_avg,
        AVG(CASE WHEN (t.sentiments->>'j')::numeric BETWEEN 1 AND 5 THEN (t.sentiments->>'j')::numeric END) AS joy_avg,
        AVG(CASE WHEN (t.sentiments->>'a_n')::numeric BETWEEN 1 AND 5 THEN (t.sentiments->>'a_n')::numeric END) AS anger_avg,
        AVG(CASE WHEN (t.sentiments->>'s')::numeric BETWEEN 1 AND 5 THEN (t.sentiments->>'s')::numeric END) AS sadness_avg,
        AVG(CASE WHEN (t.sentiments->>'f')::numeric BETWEEN 1 AND 5 THEN (t.sentiments->>'f')::numeric END) AS fear_avg,
        AVG(CASE WHEN (t.sentiments->>'d_g')::numeric BETWEEN 1 AND 5 THEN (t.sentiments->>'d_g')::numeric END) AS disgust_avg
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
)
SELECT
    trend_sentiments.trend,
    trend_sentiments.domain,
    trend_sentiments.language,
    domain_sentiments.time,
    trend_sentiments.article_count || '(' || domain_sentiments.all_article_count || ')' AS article_count,
    ROUND((trend_sentiments.valence_avg - domain_sentiments.valence_avg), 2) AS valence,
    ROUND((trend_sentiments.arousal_avg - domain_sentiments.arousal_avg), 2) AS arousal,
    ROUND((trend_sentiments.dominance_avg - domain_sentiments.dominance_avg), 2) AS dominance,
    ROUND((trend_sentiments.joy_avg - domain_sentiments.joy_avg), 2) AS joy,
    ROUND((trend_sentiments.anger_avg - domain_sentiments.anger_avg), 2) AS anger,
    ROUND((trend_sentiments.sadness_avg - domain_sentiments.sadness_avg), 2) AS sadness,
    ROUND((trend_sentiments.fear_avg - domain_sentiments.fear_avg), 2) AS fear,
    ROUND((trend_sentiments.disgust_avg - domain_sentiments.disgust_avg), 2) AS disgust
FROM trend_sentiments
JOIN domain_sentiments
    ON domain_sentiments.domain = trend_sentiments.domain
   AND domain_sentiments.language = trend_sentiments.language
ORDER BY
    trend_sentiments.trend,
    trend_sentiments.domain,
    trend_sentiments.language;
