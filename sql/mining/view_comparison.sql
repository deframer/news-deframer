SET duckdb.force_execution = true;

WITH comparison_stats AS (
    SELECT
        unnest(noun_stems) as stem,
        -- Create a label for the group
        CASE
            WHEN 'sports' = ANY(category_stems) THEN 'Sports'
            WHEN 'politics' = ANY(category_stems) THEN 'Politics'
        END as group_name,
        count(*) as freq
    FROM public.trends
    WHERE "language" = 'en'
      AND pub_date >= NOW() - INTERVAL '24 HOURS'
      AND (
          'sports' = ANY(category_stems) OR
          'politics' = ANY(category_stems)
      )
    GROUP BY 1, 2
)
-- Find the Intersection (Words mentioned in BOTH groups)
SELECT stem, sum(freq) as total_freq
FROM comparison_stats
GROUP BY stem
HAVING count(distinct group_name) = 2 -- Appears in both Sports AND Politics
ORDER BY total_freq DESC;