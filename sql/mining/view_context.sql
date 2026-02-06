---Query: "How is [Noun] being described?" This query finds the Adjectives/Verbs that appear most frequently in the same articles as your target Noun.

-- "Show me the context (Verbs/Adjectives) for the trend 'bitcoin'"
-- "How is 'bitcoin' being described in English news?"
SET duckdb.force_execution = true;

WITH trend_context AS (
    SELECT
        unnest(verb_stems) as context_word,
        'VERB' as type
    FROM public.trends
    WHERE 'trump' = ANY(noun_stems)      -- The Trend Trigger
    --WHERE 'trump' = ANY(noun_stems)      -- The Trend Trigger
    --  AND "language" = 'en'
      AND "language" = 'de'
      --AND pub_date >= NOW() - INTERVAL '7 DAYS'
)
SELECT context_word, count(*) as frequency
FROM trend_context
GROUP BY 1
ORDER BY 2 DESC, 1
LIMIT 10;
