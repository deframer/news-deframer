SELECT DISTINCT ON (i.url)
    i.url AS url,
    i.think_result->>'title_corrected' AS title,
    CAST(i.think_result->>'overall' AS FLOAT) AS rating
FROM public.trends t
JOIN public.items i ON t.item_id = i.id
WHERE
    @term = ANY(t.noun_stems)
    AND t.root_domain = @domain
    AND t.pub_date::DATE = @date::DATE
ORDER BY
    i.url,
    t.pub_date DESC
