SELECT feeds.url, feed_schedules.* FROM
    feed_schedules
JOIN
    feeds ON feed_schedules.id = feeds.id
WHERE
    next_thinker_at <= NOW()
AND
    (thinker_locked_until IS NULL OR thinker_locked_until < NOW())
ORDER BY
    next_thinker_at ASC
LIMIT 1;