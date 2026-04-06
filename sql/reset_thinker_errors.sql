BEGIN TRANSACTION;

-- SELECT count(id) FROM items WHERE think_error_count > 3;
-- SELECT count(id) FROM items WHERE think_error_count = 0 and think_error IS NOT NULL;

UPDATE items SET think_error_count = 4 WHERE think_error_count < 3 AND think_error IS NULL AND think_result IS NULL;

--SELECT updated_at, think_error, think_result, think_error_count FROM items WHERE think_error_count >0;

-- UPDATE items SET think_error_count = 3 WHERE think_error_count > 3;

ROLLBACK;