BEGIN TRANSACTION;

-- SELECT count(id) FROM items WHERE think_error_count > 3;
-- SELECT count(id) FROM items WHERE think_error_count = 0 and think_error IS NOT NULL;

UPDATE items SET think_error_count = 0, think_error = NULL WHERE think_error_count > 3;

ROLLBACK;