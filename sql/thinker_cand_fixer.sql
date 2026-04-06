SELECT
  COUNT(*) FILTER (WHERE think_error_count > 3)
FROM items;
