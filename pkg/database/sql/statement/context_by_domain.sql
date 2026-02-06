/*
   QUERY: Domain-Specific Narrative (Action + Quality)
   concatenates Verbs and Adjectives to form a unified context string
   filtered specifically for a domain.
*/

WITH settings AS (
    SELECT set_config('duckdb.force_execution', 'true', true)
),
domain_narrative AS (
    SELECT
        -- 1. The Trigger ($T$)
        CAST(@term AS text) as trigger,

        -- 2. Context ($C$) = Action (Verb)
        -- Unnest flattens the array. If the array is empty, the row is dropped (implicit filtering).
        unnest(verb_stems) as action

    FROM public.trends
    WHERE
      -- Ensure we match lower-case stems (Trend Mining convention)
      LOWER(CAST(@term AS text)) = ANY(noun_stems)

      AND "language" = CAST(@language AS text)
      AND root_domain = CAST(@domain AS text)
      AND pub_date >= NOW() - (CAST(@days_in_past AS INTEGER) * INTERVAL '1 DAY')
)
SELECT
    action as context_word,
    count(*) as frequency
FROM domain_narrative
WHERE action IS NOT NULL
  -- FILTER: Remove self-references (e.g. if 'Trump' is misclassified as a verb)
  AND action != LOWER(CAST(@term AS text))
GROUP BY 1
ORDER BY 2 DESC
LIMIT 20;
