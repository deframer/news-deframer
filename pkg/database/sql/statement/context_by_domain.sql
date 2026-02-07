/*
   QUERY: Domain-Specific Trend Context (Actions)
   - Focuses purely on 'Action' (Verbs) to determine the event structure or relations.
   - Uses Internal Utility (Feed Diversity) to filter out single-section noise.
   - Thesis: Context ($C$) is the set of relations/actions surrounding the trigger [Source 313].
*/

WITH settings AS (
    SELECT set_config('duckdb.force_execution', 'true', true)
),
domain_actions AS (
    SELECT
        -- 1. The Trigger ($T$): The Subject
        CAST(@term AS text) as trigger,

        -- Capture Feed ID to calculate Internal Utility (Section Diversity)
        feed_id,

        -- 2. Extract Context ($C$): The Action/Event (Verbs)
        -- Thesis: Verbs define the event structure or relation [Source 313]
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
    -- 3. The Action
    action as context_word,

    -- 4. Frequency (Intensity)
    count(*) as frequency,

    -- 5. Internal Utility (Section Diversity)
    -- Thesis: Measures if the term is useful or just noise [Source 33].
    -- Since domain is fixed, we count DISTINCT feeds to prove broad coverage.
    count(DISTINCT feed_id) as utility

FROM domain_actions
WHERE action IS NOT NULL
  -- FILTER: Remove self-references (e.g. if 'Trump' is misclassified as a verb)
  AND action != LOWER(CAST(@term AS text))
  -- Optional: Remove generic verb stopwords to improve signal quality
  -- AND action NOT IN ('be', 'have', 'do', 'say', 'make', 'go', 'take')

GROUP BY 1

-- 6. Apply Utility Threshold
-- Must appear in at least 2 different feeds/sections of the site to be considered a broad trend.
--HAVING count(DISTINCT feed_id) > 1

ORDER BY utility DESC, frequency DESC, context_word
LIMIT 20;
