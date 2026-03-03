/*
   QUERY: Daily Article Retrieval by Domain and Term
   - Logic: Retrieves the ID, URL, and extracted title for articles matching a specific term on a specific date.
   - Scope: Filters for the given domain ('tagesschau.de'), exact date ('2026-02-10'), and term ('trump').
   - Output: Returns the article ID, URL, and the title extracted from the JSON.
*/

SELECT
    i.id,
    i.url,
    i.think_result->>'title_corrected' AS title,
    i.think_result->>'overall' AS overall

FROM public.trends t
JOIN public.items i ON t.item_id = i.id

WHERE
    -- Trigger Selection ($T$)
    'trump' = ANY(t.noun_stems)

    -- Domain/Scope Filter
    AND t.root_domain = 'tagesschau.de'

    -- Exact Date Filter ($TW$)
    AND t.pub_date::DATE = '2026-02-10'

-- Sort by newest articles first
ORDER BY
    t.pub_date DESC,
    title
