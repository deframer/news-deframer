BEGIN TRANSACTION;

-- SELECT count(id) FROM items WHERE think_error_count > 3;

UPDATE items SET think_error_count = 0 WHERE think_error_count > 3;

ROLLBACK;