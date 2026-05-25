BEGIN TRANSACTION;

SELECT
    COALESCE(NULLIF(items.think_result ->> 'llm_model', ''), '<missing>') AS llm_model,
    COUNT(*) AS item_count
FROM items
WHERE items.think_result IS NOT NULL
GROUP BY 1
ORDER BY item_count DESC, llm_model;

-- UPDATE items
-- SET think_result = jsonb_set(
--     items.think_result,
--     '{llm_model}',
--     to_jsonb('mistralai/ministral-3-3b'::text),
--     true
-- )
-- WHERE items.think_result IS NOT NULL
--   AND COALESCE(NULLIF(items.think_result ->> 'llm_model', ''), '') = '';

-- COMMIT;
ROLLBACK;
