BEGIN TRANSACTION;

-- SELECT * FROM trends WHERE sentiments <> '{}'::jsonb LIMIT 10;

UPDATE trends
SET sentiments = '{}'::jsonb;

ROLLBACK;
