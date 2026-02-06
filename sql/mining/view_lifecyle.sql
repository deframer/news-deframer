SET duckdb.force_execution = true;

-- Check the lifecycle of 'trump'
SELECT
    time_slice,
    frequency,
    growth_delta as velocity
FROM view_trend_metrics
WHERE stem = 'trump'
  AND "language" = 'de'
--AND time_slice = DATE_TRUNC('day', NOW())  -- Or specific date
  AND time_slice >= NOW() - INTERVAL '7 DAYS'
ORDER BY time_slice DESC;
