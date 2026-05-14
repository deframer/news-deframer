SELECT
    language,
    category,
    COUNT(*) AS article_count
FROM (
    SELECT DISTINCT
        id,
        language,
        NULLIF(think_result ->> 'category', '') AS category
    FROM items
    WHERE think_result IS NOT NULL
) item_categories
WHERE category IS NOT NULL
GROUP BY language, category
ORDER BY language, category;
