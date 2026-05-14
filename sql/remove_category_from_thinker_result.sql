BEGIN TRANSACTION;

-- SELECT
--     COUNT(*) AS item_count
-- FROM items
-- WHERE think_result IS NOT NULL
--   AND NULLIF(think_result ->> 'category', '') IS NOT NULL;

UPDATE items
SET think_result = think_result - 'category'
WHERE think_result IS NOT NULL
  AND think_result ? 'category';

ROLLBACK;
