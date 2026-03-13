SELECT
    COUNT(*) FILTER (WHERE sentiments = '{}'::jsonb) AS empty_sentiments_trend_count,
    COUNT(*) FILTER (WHERE sentiments <> '{}'::jsonb) AS non_empty_sentiments_trend_count
FROM trends;
