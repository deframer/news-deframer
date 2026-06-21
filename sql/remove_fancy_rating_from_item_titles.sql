BEGIN TRANSACTION;

SELECT COUNT(*) AS affected_items
FROM items
WHERE items.content <> regexp_replace(
    items.content,
    '(<title[^>]*>)[[:space:]]*(?:[★☆]+[[:space:]]*)?',
    E'\\1',
    'g'
);

-- UPDATE items
-- SET content = regexp_replace(
--     items.content,
--     '(<title[^>]*>)[[:space:]]*(?:[★☆]+[[:space:]]*)?',
--     E'\\1',
--     'g'
-- )
-- WHERE items.content <> regexp_replace(
--     items.content,
--     '(<title[^>]*>)[[:space:]]*(?:[★☆]+[[:space:]]*)?',
--     E'\\1',
--     'g'
-- );

-- COMMIT;
ROLLBACK;
