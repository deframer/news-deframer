SELECT root_domain AS domain
FROM feeds
WHERE enabled = true
  AND deleted_at IS NULL
  AND root_domain IS NOT NULL
  AND root_domain <> ''
ORDER BY root_domain ASC;
