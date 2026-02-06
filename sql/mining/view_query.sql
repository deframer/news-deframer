-- Remember to enable DuckDB execution
SET duckdb.force_execution = true;

-- Get top trends specifically for English
SELECT *
FROM view_trend_metrics
WHERE "language" = 'en'
  --AND time_slice = DATE_TRUNC('day', NOW())
  AND time_slice >= NOW() - INTERVAL '7 DAYS'
ORDER BY outlier_ratio DESC  NULLS LAST;
