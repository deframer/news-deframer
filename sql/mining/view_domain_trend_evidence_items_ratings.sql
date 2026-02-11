/*
   QUERY: Domain-Specific Evidence Search with Quality Scores
   - Logic: Retrieves individual articles containing the trigger 'trump'.
   - Detail: Extracts raw AI quality scores (ThinkResult) for each article.
   - Handling Defaults: Uses COALESCE to treat missing JSON keys as 0.0.
   - Output: Date, URL, and the specific scores for Framing, Clickbait, etc.
*/

SELECT
    -- 1. Time Dimension
    t.pub_date::DATE as trend_date,

    -- 2. Evidence Details
    -- t.item_id,
    i.url as item_url,

    -- 3. Individual Quality Scores (extracted from ThinkResult JSONB)
    -- We cast to float and use COALESCE to handle 'omitempty' (NULL -> 0.0)
    COALESCE(CAST(i.think_result->>'framing' AS FLOAT), 0.0) as framing,
    COALESCE(CAST(i.think_result->>'clickbait' AS FLOAT), 0.0) as clickbait,
    COALESCE(CAST(i.think_result->>'persuasive' AS FLOAT), 0.0) as persuasive,
    COALESCE(CAST(i.think_result->>'hyper_stimulus' AS FLOAT), 0.0) as hyper_stimulus,
    COALESCE(CAST(i.think_result->>'speculative' AS FLOAT), 0.0) as speculative,
    COALESCE(CAST(i.think_result->>'overall' AS FLOAT), 0.0) as overall

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
