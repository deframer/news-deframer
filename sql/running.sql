SELECT
  --feeds.id,
  feeds.url,
  --feeds.root_domain,
  feed_schedules.next_thinker_at,
  feed_schedules.thinker_locked_until,
  --feed_schedules.next_mining_at,
  feed_schedules.mining_locked_until
FROM
  feed_schedules
  INNER JOIN feeds ON (feed_schedules.id = feeds.id)
WHERE
  feeds.enabled = true
  AND feeds.deleted_at IS NULL
  AND (
    feed_schedules.thinker_locked_until IS NOT NULL
    OR
    feed_schedules.mining_locked_until IS NOT NULL
  )
