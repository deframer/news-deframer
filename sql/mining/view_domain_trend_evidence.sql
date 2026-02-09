/*
   QUERY: Domain-Specific Trend Evidence (Daily Arrays)
   - Logic: Aggregates item_ids into daily arrays for the trigger 'trump'.
   - Scope: Filters for 'tagesschau.de' over the last 7 days.
   - Output: Returns the date and an array of all matching item IDs for that day.
*/

SELECT
    -- 1. Time Dimension (Date)
    pub_date::DATE as trend_date,

    -- 2. Evidence (Array of IDs)
    -- We use array_agg to collect all IDs for the day into a single list
    array_agg(item_id) as evidence_ids

FROM public.trends

WHERE
    -- Trigger Selection ($T$)
    'trump' = ANY(noun_stems)

    -- Domain/Scope Filters
    AND "language" = 'de'
    AND root_domain = 'tagesschau.de'

    -- Time Window ($TW$)
    AND pub_date >= NOW() - INTERVAL '7 DAYS'

-- Group by the date to create the daily arrays
GROUP BY 1

-- Sort by newest date first
ORDER BY 1 DESC;
