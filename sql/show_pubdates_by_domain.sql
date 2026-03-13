SELECT
    f.root_domain AS domain,
    COUNT(i.id) AS articles,
    MIN(i.pub_date::date) AS oldest_article,
    MAX(i.pub_date::date) AS newest_article,
    MAX(i.pub_date::date) - MIN(i.pub_date::date) AS days
FROM items i
JOIN feeds f ON f.id = i.feed_id
WHERE
    i.think_result IS NOT NULL
    AND i.think_error IS NULL
    AND i.think_error_count = 0
    AND f.root_domain IS NOT NULL
    AND f.root_domain <> ''
    -- AND f.root_domain in ( 'apollo-news.net', 'nius.de', 'tagesschau.de', 'faz.net', 'welt.de' )
GROUP BY f.root_domain
ORDER BY f.root_domain ASC;
