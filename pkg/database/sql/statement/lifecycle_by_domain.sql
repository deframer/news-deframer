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
  AND time_slice >= COALESCE(
        CAST(@date AS timestamp) - ((CAST(@days AS INTEGER) - 1) * INTERVAL '1 DAY'),
        NOW() - (CAST(@days AS INTEGER) * INTERVAL '1 DAY')
  )
  AND time_slice < COALESCE(
        CAST(@date AS timestamp) + INTERVAL '1 DAY',
        NOW()
  )
ORDER BY time_slice ASC;
