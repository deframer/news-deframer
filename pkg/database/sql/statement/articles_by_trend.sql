SELECT DISTINCT ON (i.url)
    i.url AS url,
    COALESCE(NULLIF(i.think_result->>'title_corrected', ''), i.think_result->>'title_original') AS title,
    i.think_rating AS rating
FROM public.trends t
JOIN public.items i ON t.item_id = i.id
WHERE
    @term = ANY(t.noun_stems)
    AND t.root_domain = @domain
    AND t.pub_date::DATE = CAST(@date AS DATE)
    AND i.think_result IS NOT NULL
    AND i.think_error IS NULL
    AND i.think_error_count = 0
ORDER BY
    i.url,
    t.pub_date DESC
