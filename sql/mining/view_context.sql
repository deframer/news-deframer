/*
   QUERY: Trend Context (Actions) with Utility Score
   - Focuses purely on 'Action' (Verbs) to determine the event structure/relations.
   - Uses Utility to filter out "echo chamber" actions used by only one source.
   - Thesis: Verbs define the Context ($C$) (Actions/Events) [Source 313].
*/

-- Force DuckDB execution using CTE (as requested)
WITH settings AS (
    SELECT set_config('duckdb.force_execution', 'true', true)
),
trend_actions AS (
    SELECT
        -- 1. The Trigger
        'trump' as trigger,

        -- Capture Root Domain to calculate Utility (Source Diversity)
        root_domain,

        -- 2. Extract Context ($C$): The Action/Event (Verbs)
        -- We removed 'quality' (adjectives) to focus purely on the actions [Source 313]
        unnest(verb_stems) as action

    FROM public.trends
    WHERE 'trump' = ANY(noun_stems)
      AND "language" = 'de'
      AND pub_date >= NOW() - INTERVAL '7 DAYS'
)
SELECT
    -- 3. The Action
    action as context_word,

    -- 4. Frequency (Intensity)
    count(*) as frequency,

    -- 5. Utility (Source Diversity)
    -- Thesis: Measures if the term is a broad consensus or single-source noise [Source 33]
    count(DISTINCT root_domain) as utility

FROM trend_actions
WHERE action IS NOT NULL
  -- FILTER: Remove self-references and common verb stopwords
  AND action != 'trump'
  -- AND action NOT IN ('be', 'have', 'do', 'say') -- Optional: Filter generic verbs

GROUP BY 1

-- 6. Apply Utility Threshold
-- Only show actions reported by at least 3 different publishers to ensure it's a real trend.
--HAVING count(DISTINCT root_domain) > 2

ORDER BY utility DESC, frequency DESC, context_word
LIMIT 20;
