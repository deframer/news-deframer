SELECT
  feeds.root_domain AS domain,
  COUNT(*) AS item_count
FROM items
JOIN feeds ON feeds.id = items.feed_id
WHERE think_result IS NULL
  AND (think_error_count IS NULL OR think_error_count < 7)
  AND feeds.enabled = true
  AND feeds.deleted_at IS NULL
  AND feeds.root_domain IS NOT NULL
  AND feeds.root_domain <> ''
GROUP BY feeds.root_domain
ORDER BY item_count DESC, domain;
