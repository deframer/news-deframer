/*
   QUERY: Domain-Specific Trend Evidence with Qualitative Scoring (Defaults Handled)
   - Logic: Aggregates item_ids and averages AI quality scores per day.
   - Fix: Uses COALESCE to treat missing JSON keys (omitempty) as 0.0.
     This ensures the average correctly accounts for "zero" scores rather than ignoring them.
   - Output: Date, Evidence IDs, and corrected average scores.
*/

SELECT
    -- 1. Time Dimension
    t.pub_date::DATE as trend_date,

    -- 2. Evidence (Array of IDs)
    array_agg(t.item_id) as evidence_ids,

    -- 3. Aggregated Quality Scores
    -- We cast to float, then use COALESCE to convert NULL (missing) to 0.0
    -- BEFORE averaging.
    ROUND(AVG(COALESCE(CAST(i.think_result->>'framing' AS FLOAT), 0.0))::NUMERIC, 2) as avg_framing,
    ROUND(AVG(COALESCE(CAST(i.think_result->>'clickbait' AS FLOAT), 0.0))::NUMERIC, 2) as avg_clickbait,
    ROUND(AVG(COALESCE(CAST(i.think_result->>'persuasive' AS FLOAT), 0.0))::NUMERIC, 2) as avg_persuasive,
    ROUND(AVG(COALESCE(CAST(i.think_result->>'hyper_stimulus' AS FLOAT), 0.0))::NUMERIC, 2) as avg_hyper_stimulus,
    ROUND(AVG(COALESCE(CAST(i.think_result->>'speculative' AS FLOAT), 0.0))::NUMERIC, 2) as avg_speculative,
    ROUND(AVG(COALESCE(CAST(i.think_result->>'overall' AS FLOAT), 0.0))::NUMERIC, 2) as avg_overall

FROM public.trends t
JOIN public.items i ON t.item_id = i.id

WHERE
    -- Trigger Selection ($T$)
    'ministerpräsident' = ANY(t.noun_stems)

    -- Domain/Scope Filters
    AND t."language" = 'de'
    AND t.root_domain = 'nius.de'

    -- Time Window ($TW$)
    AND t.pub_date >= NOW() - INTERVAL '30 DAYS'

-- Group by date
GROUP BY 1

-- Sort by newest date first
ORDER BY 1 DESC;
