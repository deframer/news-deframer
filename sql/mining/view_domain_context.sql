/*
   QUERY: Domain-Specific Narrative (Action + Quality)
   concatenates Verbs and Adjectives to form a unified context string
   filtered specifically for 'bbc.com'.
*/

SET duckdb.force_execution = true;

WITH bbc_narrative AS (
    SELECT
        -- 1. The Trigger ($T$): The Subject
        'trump' as trigger,

        -- 2. Unroll arrays (DuckDB zips these lists, padding shorter ones with NULL)
        -- Thesis: Context ($C$) = Action (Verb) [Source 313]
        unnest(verb_stems) as action,

        -- Thesis: Diversificator = Quality/Sentiment (Adjective) [Source 284]
        unnest(adjective_stems) as quality

    FROM public.trends
    WHERE 'trump' = ANY(noun_stems)
      AND "language" = 'en'
      AND root_domain = 'bbc.com'        -- <--- FOCUS ON BBC ONLY
      AND pub_date >= NOW() - INTERVAL '7 DAYS' -- Extended to 7 days for single domain density
)
SELECT
    -- 3. Concat with ' ' separator.
    -- Result examples: "stated controversial", "attacked unfair", "won"
    CONCAT_WS(' ', action, quality) as context,

    count(*) as frequency
FROM bbc_narrative
WHERE (action IS NOT NULL OR quality IS NOT NULL)  -- Remove rows where both are empty
  -- FILTER HERE: Remove self-references to clean up the narrative
  AND quality != 'trump'
  AND action != 'trump'
GROUP BY 1
ORDER BY 2 DESC
LIMIT 20;
