BEGIN TRANSACTION;

-- SELECT * FROM items WHERE think_result IS NOT NULL LIMIT 10;

UPDATE items SET think_result = NULL,
             think_error = NULL,
             think_rating = 0;

ROLLBACK;