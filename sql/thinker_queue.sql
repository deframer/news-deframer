SELECT
  think_error_count,
  COUNT(*) AS item_count
FROM items
WHERE think_result IS NULL
GROUP BY think_error_count
ORDER BY think_error_count;
