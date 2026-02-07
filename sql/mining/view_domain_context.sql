/*
   QUERY: Domain-Specific Trend Context (Actions)
   - Focuses purely on 'Action' (Verbs) to determine the event structure or relations on 'tagesschau.de'.
   - Uses Internal Utility (Feed Diversity) to filter out single-section noise.
   - Thesis: Context ($C$) is the set of relations/actions surrounding the trigger [Source 313].
*/

SET duckdb.force_execution = true;

WITH domain_actions AS (
    SELECT
        -- 1. The Trigger ($T$)
        'trump' as trigger,

        -- Capture Feed ID to calculate Internal Utility (Section Diversity)
        feed_id,

        -- 2. Extract Context ($C$): The Action/Event (Verbs)
        -- We switched from 'quality' (adjectives) to 'action' (verbs) [Source 313]
        unnest(verb_stems) as action

    FROM public.trends
    WHERE 'trump' = ANY(noun_stems)
      AND "language" = 'de'
      AND root_domain = 'tagesschau.de'        -- <--- FOCUS ON SPECIFIC DOMAIN
      AND pub_date >= NOW() - INTERVAL '7 DAYS'
)
SELECT
    -- 3. The Action
    action as context_word,

    -- 4. Frequency (Intensity)
    count(*) as frequency,

    -- 5. Internal Utility (Section Diversity)
    -- Thesis: Measures if the term is useful or just noise [Source 33].
    -- Since domain is fixed, we count DISTINCT feeds to prove broad coverage within the site.
    count(DISTINCT feed_id) as utility

FROM domain_actions
WHERE action IS NOT NULL
  -- FILTER: Remove self-references and common verb stopwords
  AND action != 'trump'
  -- AND action NOT IN ('be', 'have', 'do', 'say') -- Optional: Filter generic verbs

GROUP BY 1

-- 6. Apply Utility Threshold
-- Must appear in at least 2 different feeds/sections of the site to be a "trend".
--HAVING count(DISTINCT feed_id) > 1

ORDER BY utility DESC, frequency DESC, context_word
LIMIT 20;
