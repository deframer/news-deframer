WITH config AS (
    SELECT 'theepochtimes.com' AS root_domain
)
SELECT
    items.*
FROM items
JOIN feeds ON feeds.id = items.feed_id
CROSS JOIN config
WHERE feeds.root_domain = config.root_domain
ORDER BY items.pub_date DESC, items.id DESC
LIMIT 10;
