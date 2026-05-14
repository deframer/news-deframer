SELECT
    feeds.root_domain AS domain,
    COUNT(*) FILTER (WHERE COALESCE(cardinality(items.categories), 0) > 0) AS items_with_categories,
    COUNT(*) FILTER (WHERE COALESCE(cardinality(items.categories), 0) = 0) AS items_without_categories
FROM items
JOIN feeds ON feeds.id = items.feed_id
GROUP BY feeds.root_domain
ORDER BY feeds.root_domain;
