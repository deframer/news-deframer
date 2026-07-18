SELECT
    COALESCE(i.language, f.language) AS language,
    f.country,
    COUNT(*) AS article_count
FROM items i
JOIN feeds f ON f.id = i.feed_id
WHERE f.deleted_at IS NULL
  AND f.enabled = true
GROUP BY f.country, COALESCE(i.language, f.language)
ORDER BY f.country, COALESCE(i.language, f.language);
