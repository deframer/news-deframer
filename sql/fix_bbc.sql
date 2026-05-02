BEGIN TRANSACTION;

-- SELECT COUNT(*)
-- FROM items
-- JOIN feeds ON feeds.id = items.feed_id
-- WHERE feeds.root_domain = 'bbc.com'
--  AND items.url LIKE '%?%';

UPDATE items
SET url = split_part(items.url, '?', 1)
FROM feeds
WHERE feeds.id = items.feed_id
  AND feeds.root_domain = 'bbc.com'
  AND items.url LIKE '%?%';

ROLLBACK;
