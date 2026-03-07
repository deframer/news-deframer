SELECT
  COUNT(*) FILTER (WHERE think_error_count > 0) AS think_error_count,
  COUNT(*) FILTER (
    WHERE think_error_count > 0
      AND think_error LIKE '%Context size has been exceeded%'
  ) AS context_size_exceeded_count
FROM items;
