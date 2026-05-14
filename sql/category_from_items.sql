SELECT
    language,
    category,
    COUNT(*) AS article_count
FROM (
    SELECT DISTINCT
        id,
        language,
        unnest(categories) AS category
    FROM items
) item_categories
WHERE category IS NOT NULL
  AND category <> ''
GROUP BY language, category
ORDER BY language, category;
