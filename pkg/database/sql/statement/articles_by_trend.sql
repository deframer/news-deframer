WITH ranked_articles AS (
    SELECT
        i.url AS url,
        NULLIF(i.think_result->>'title_corrected', '') AS title,
        i.think_rating AS rating,
        t.pub_date AS pub_date,
        ROW_NUMBER() OVER (PARTITION BY i.url ORDER BY t.pub_date DESC) AS rn
    FROM public.trends t
    JOIN public.items i ON t.item_id = i.id
    WHERE
        @term = ANY(t.noun_stems)
        AND t.root_domain = @domain
        AND t.pub_date >= (
            COALESCE(CAST(NULLIF(@date, '') AS DATE), CURRENT_DATE)::timestamp
            - ((CAST(@days AS INTEGER) - 1) * INTERVAL '1 DAY')
        )
        AND t.pub_date < (COALESCE(CAST(NULLIF(@date, '') AS DATE), CURRENT_DATE)::timestamp + INTERVAL '1 DAY')
        AND i.think_result IS NOT NULL
        AND i.think_error IS NULL
        AND i.think_error_count = 0
)
SELECT
    url,
    title,
    rating,
    pub_date
FROM ranked_articles
WHERE rn = 1
ORDER BY
    pub_date DESC,
    title ASC NULLS LAST
