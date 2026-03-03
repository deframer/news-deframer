SELECT
    i.url AS url,
    NULLIF(i.think_result->>'title_corrected', '') AS title,
    CASE WHEN i.think_result IS NOT NULL THEN i.think_rating ELSE NULL END AS rating,
    t.pub_date AS pub_date
FROM public.trends t
JOIN public.items i ON t.item_id = i.id
WHERE
    LOWER(CAST(@term AS text)) = ANY(t.noun_stems)
    AND t.root_domain = @domain
    AND t.pub_date >= COALESCE(
        CAST(NULLIF(@date, '') AS DATE)::timestamp - ((GREATEST(CAST(@days AS INTEGER), 1) - 1) * INTERVAL '1 DAY'),
        NOW() - (GREATEST(CAST(@days AS INTEGER), 1) * INTERVAL '1 DAY')
    )
    AND t.pub_date < COALESCE(
        CAST(NULLIF(@date, '') AS DATE)::timestamp + INTERVAL '1 DAY',
        NOW()
    )
ORDER BY
    t.pub_date DESC,
    title ASC NULLS LAST
