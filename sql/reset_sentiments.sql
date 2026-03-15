BEGIN TRANSACTION;

-- SELECT * FROM trends WHERE sentiments <> '{}'::jsonb OR sentiments_deframed <> '{}'::jsonb LIMIT 10;

UPDATE trends
SET sentiments = '{}'::jsonb,
    sentiments_deframed = '{}'::jsonb;

ROLLBACK;
