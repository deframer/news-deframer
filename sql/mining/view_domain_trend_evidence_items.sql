/*
   QUERY: Domain-Specific Evidence Search
   - Logic: Retrieves a list of articles ('evidence') containing the trigger 'trump'.
   - Scope: Filters for the specific domain 'tagesschau.de' over the last 7 days.
   - Output: Returns the date and article URL for every matching instance found.
*/

SELECT
    -- 1. Time Dimension (Date)
    t.pub_date::DATE as trend_date,

    -- 2. Evidence Details
    -- t.item_id,
    i.url as item_url

FROM public.trends t
JOIN public.items i ON t.item_id = i.id

WHERE
    -- Trigger Selection ($T$)
    'trump' = ANY(t.noun_stems)

    -- Domain/Scope Filters
    AND t."language" = 'de'
    AND t.root_domain = 'tagesschau.de'

    -- Time Window ($TW$)
    AND t.pub_date >= NOW() - INTERVAL '7 DAYS'

-- Sort by newest articles first
ORDER BY t.pub_date DESC;
