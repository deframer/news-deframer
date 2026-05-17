SELECT
  root_domain,
  url,
  last_error AS error
FROM feeds
WHERE last_error IS NOT NULL
ORDER BY root_domain ASC, url ASC;
