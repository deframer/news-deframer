SELECT
    DATE_TRUNC('month', pub_date)::date AS month,
    COUNT(*) AS article_count,
    COUNT(DISTINCT feed_id) AS feed_count
FROM items
WHERE pub_date >= NOW() - INTERVAL '30 DAYS'
GROUP BY 1
ORDER BY 1;
