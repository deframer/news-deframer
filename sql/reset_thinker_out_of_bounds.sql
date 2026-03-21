BEGIN TRANSACTION;

/*
SELECT
    id,
    url,
    think_result ->> 'framing' AS framing,
    think_result ->> 'clickbait' AS clickbait,
    think_result ->> 'persuasive' AS persuasive,
    think_result ->> 'hyper_stimulus' AS hyper_stimulus,
    think_result ->> 'speculative' AS speculative,
    think_result ->> 'overall' AS overall
FROM items
WHERE think_result IS NOT NULL
  AND (
    (think_result ->> 'framing')::numeric < 0
    OR (think_result ->> 'framing')::numeric > 1
    OR (think_result ->> 'clickbait')::numeric < 0
    OR (think_result ->> 'clickbait')::numeric > 1
    OR (think_result ->> 'persuasive')::numeric < 0
    OR (think_result ->> 'persuasive')::numeric > 1
    OR (think_result ->> 'hyper_stimulus')::numeric < 0
    OR (think_result ->> 'hyper_stimulus')::numeric > 1
    OR (think_result ->> 'speculative')::numeric < 0
    OR (think_result ->> 'speculative')::numeric > 1
    OR (think_result ->> 'overall')::numeric < 0
    OR (think_result ->> 'overall')::numeric > 1
  );
*/

UPDATE items SET think_result = NULL,
             think_error = NULL,
             think_rating = 0
WHERE think_result IS NOT NULL
  AND (
    (think_result ->> 'framing')::numeric < 0
    OR (think_result ->> 'framing')::numeric > 1
    OR (think_result ->> 'clickbait')::numeric < 0
    OR (think_result ->> 'clickbait')::numeric > 1
    OR (think_result ->> 'persuasive')::numeric < 0
    OR (think_result ->> 'persuasive')::numeric > 1
    OR (think_result ->> 'hyper_stimulus')::numeric < 0
    OR (think_result ->> 'hyper_stimulus')::numeric > 1
    OR (think_result ->> 'speculative')::numeric < 0
    OR (think_result ->> 'speculative')::numeric > 1
    OR (think_result ->> 'overall')::numeric < 0
    OR (think_result ->> 'overall')::numeric > 1
  );

ROLLBACK;
