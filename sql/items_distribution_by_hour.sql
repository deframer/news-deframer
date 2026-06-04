WITH config AS (
    SELECT
        365 AS days
), base AS (
    SELECT
        f.root_domain AS domain,
        f.language AS language,
        f.country AS country,
        EXTRACT(HOUR FROM timezone('UTC', i.pub_date))::int AS hour_utc
    FROM items i
    JOIN feeds f ON f.id = i.feed_id
    CROSS JOIN config c
    WHERE
        f.enabled = true
        AND f.deleted_at IS NULL
        AND f.language IS NOT NULL
        AND f.country IS NOT NULL
        AND f.root_domain IS NOT NULL
        AND f.root_domain <> ''
        AND i.pub_date >= NOW() - (c.days * INTERVAL '1 DAY')
), domains AS (
    SELECT DISTINCT domain, language, country
    FROM base
), hours AS (
    SELECT generate_series(0, 23) AS hour_of_day
), utc_grid AS (
    SELECT
        d.domain,
        d.language,
        d.country,
        h.hour_of_day
    FROM domains d
    CROSS JOIN hours h
), utc_counts AS (
    SELECT
        g.domain,
        g.language,
        g.country,
        g.hour_of_day,
        COUNT(b.hour_utc) AS article_count
    FROM utc_grid g
    LEFT JOIN base b
        ON b.domain = g.domain
        AND b.language = g.language
        AND b.country = g.country
        AND b.hour_utc = g.hour_of_day
    GROUP BY g.domain, g.language, g.country, g.hour_of_day
), utc_percentages AS (
    SELECT
        domain,
        language,
        country,
        hour_of_day,
        ROUND(100.0 * article_count / NULLIF(SUM(article_count) OVER (PARTITION BY domain, language, country), 0), 2) AS percent
    FROM utc_counts
)
SELECT
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 0), 0) AS "0h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 1), 0) AS "1h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 2), 0) AS "2h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 3), 0) AS "3h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 4), 0) AS "4h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 5), 0) AS "5h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 6), 0) AS "6h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 7), 0) AS "7h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 8), 0) AS "8h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 9), 0) AS "9h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 10), 0) AS "10h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 11), 0) AS "11h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 12), 0) AS "12h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 13), 0) AS "13h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 14), 0) AS "14h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 15), 0) AS "15h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 16), 0) AS "16h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 17), 0) AS "17h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 18), 0) AS "18h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 19), 0) AS "19h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 20), 0) AS "20h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 21), 0) AS "21h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 22), 0) AS "22h",
    COALESCE(MAX(utc_percentages.percent) FILTER (WHERE utc_percentages.hour_of_day = 23), 0) AS "23h"
FROM utc_percentages
GROUP BY utc_percentages.domain, utc_percentages.language, utc_percentages.country
ORDER BY utc_percentages.domain, utc_percentages.language, utc_percentages.country;
