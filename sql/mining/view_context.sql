/*
   QUERY: Combined Trend Narrative (Action + Quality)
   concatenates Verbs and Adjectives to form a unified context string.
*/

SET duckdb.force_execution = true;

WITH trend_components AS (
    SELECT
        -- 1. The Trigger
        'trump' as trigger,

        -- 2. Unroll arrays (DuckDB zips these lists, padding shorter ones with NULL)
        unnest(verb_stems) as action,
        unnest(adjective_stems) as quality

    FROM public.trends
    WHERE 'trump' = ANY(noun_stems)
      AND "language" = 'en'
      AND pub_date >= NOW() - INTERVAL '24 HOURS'
)
SELECT
    -- 3. Concat with ' ' separator.
    -- If action is NULL, result is "quality".
    -- If quality is NULL, result is "action".
    CONCAT_WS(' ', action, quality) as context,

    count(*) as frequency
FROM trend_components
WHERE (action IS NOT NULL OR quality IS NOT NULL)  -- Remove rows where both are empty
  -- FILTER HERE: Remove self-references to clean up the narrative
  AND quality != 'trump'
  AND action != 'trump'
GROUP BY 1
ORDER BY 2 DESC
LIMIT 20;
