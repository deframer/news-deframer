WITH settings AS (
    SELECT set_config('duckdb.force_execution', 'true', true)
)
SELECT
    time_slice,
    frequency,
    growth_delta as velocity
FROM view_trend_metrics_by_domain
WHERE stem = CAST(@term AS text)
  AND root_domain = CAST(@domain AS text)
  AND "language" = CAST(@language AS text)
  AND time_slice >= NOW() - (CAST(@days_in_past AS INTEGER) * INTERVAL '1 DAY')
ORDER BY time_slice ASC;