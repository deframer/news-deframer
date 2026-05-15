WITH base AS (
    SELECT
        f.root_domain AS domain,
        f.language AS language,
        EXTRACT(HOUR FROM i.pub_date)::int AS hour_of_day
    FROM items i
    JOIN feeds f ON f.id = i.feed_id
    WHERE
        f.enabled = true
        AND f.language IS NOT NULL
        AND f.root_domain IS NOT NULL
        AND f.root_domain <> ''
), domains AS (
    SELECT DISTINCT domain, language
    FROM base
), hours AS (
    SELECT generate_series(0, 23) AS hour_of_day
), grid AS (
    SELECT
        d.domain,
        d.language,
        h.hour_of_day
    FROM domains d
    CROSS JOIN hours h
), counts AS (
    SELECT
        g.domain,
        g.language,
        g.hour_of_day,
        COUNT(b.hour_of_day) AS article_count
    FROM grid g
    LEFT JOIN base b
        ON b.domain = g.domain
       AND b.language = g.language
       AND b.hour_of_day = g.hour_of_day
    GROUP BY g.domain, g.language, g.hour_of_day
), percentages AS (
    SELECT
        domain,
        language,
        hour_of_day,
        ROUND(100.0 * article_count / NULLIF(SUM(article_count) OVER (PARTITION BY domain, language), 0), 2) AS percent
    FROM counts
)
SELECT
    domain,
    language,
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 0), 0) AS "0h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 1), 0) AS "1h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 2), 0) AS "2h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 3), 0) AS "3h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 4), 0) AS "4h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 5), 0) AS "5h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 6), 0) AS "6h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 7), 0) AS "7h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 8), 0) AS "8h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 9), 0) AS "9h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 10), 0) AS "10h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 11), 0) AS "11h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 12), 0) AS "12h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 13), 0) AS "13h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 14), 0) AS "14h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 15), 0) AS "15h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 16), 0) AS "16h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 17), 0) AS "17h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 18), 0) AS "18h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 19), 0) AS "19h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 20), 0) AS "20h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 21), 0) AS "21h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 22), 0) AS "22h",
    COALESCE(MAX(percent) FILTER (WHERE hour_of_day = 23), 0) AS "23h"
FROM percentages
GROUP BY domain, language
ORDER BY domain, language;
