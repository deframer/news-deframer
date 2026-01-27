-- PREPARE filter_reasons(float) AS
-- SELECT
--     id,
--     CASE WHEN (think_result ->> 'framing')::float > $1 THEN think_result ->> 'framing_reason' END AS framing_reason,
--     CASE WHEN (think_result ->> 'clickbait')::float > $1 THEN think_result ->> 'clickbait_reason' END AS clickbait_reason,
--     CASE WHEN (think_result ->> 'persuasive')::float > $1 THEN think_result ->> 'persuasive_reason' END AS persuasive_reason,
--     CASE WHEN (think_result ->> 'hyper_stimulus')::float > $1 THEN think_result ->> 'hyper_stimulus_reason' END AS hyper_stimulus_reason,
--     CASE WHEN (think_result ->> 'speculative')::float > $1 THEN think_result ->> 'speculative_reason' END AS speculative_reason,
--     CASE WHEN (think_result ->> 'overall')::float > $1 THEN think_result ->> 'overall_reason' END AS overall_reason
-- FROM
--     items
-- WHERE
--     think_result IS NOT NULL
--     AND (
--         (think_result ->> 'framing')::float > $1
--         OR (think_result ->> 'clickbait')::float > $1
--         OR (think_result ->> 'persuasive')::float > $1
--         OR (think_result ->> 'hyper_stimulus')::float > $1
--         OR (think_result ->> 'speculative')::float > $1
--         OR (think_result ->> 'overall')::float > $1
--     );

-- EXECUTE filter_reasons(0.7);

-- WITH config AS (
--     SELECT 0.7 AS threshold
-- )
-- SELECT
--     id,
--     array_to_string(ARRAY_REMOVE(ARRAY[
--         CASE WHEN (think_result ->> 'framing')::float > c.threshold
--              THEN 'Framing: ' || (think_result ->> 'framing_reason') END,
--         CASE WHEN (think_result ->> 'clickbait')::float > c.threshold
--              THEN 'Clickbait: ' || (think_result ->> 'clickbait_reason') END,
--         CASE WHEN (think_result ->> 'persuasive')::float > c.threshold
--              THEN 'Persuasive: ' || (think_result ->> 'persuasive_reason') END,
--         CASE WHEN (think_result ->> 'hyper_stimulus')::float > c.threshold
--              THEN 'Hyper Stimulus: ' || (think_result ->> 'hyper_stimulus_reason') END,
--         CASE WHEN (think_result ->> 'speculative')::float > c.threshold
--              THEN 'Speculative: ' || (think_result ->> 'speculative_reason') END,
--         CASE WHEN (think_result ->> 'overall')::float > c.threshold
--              THEN 'Overall: ' || (think_result ->> 'overall_reason') END
--     ], NULL), E'\n') AS reasons_text -- E'\n' puts each reason on a new line
-- FROM
--     items,
--     config c
-- WHERE
--     think_result IS NOT NULL
--     AND (
--         (think_result ->> 'framing')::float > c.threshold
--         OR (think_result ->> 'clickbait')::float > c.threshold
--         OR (think_result ->> 'persuasive')::float > c.threshold
--         OR (think_result ->> 'hyper_stimulus')::float > c.threshold
--         OR (think_result ->> 'speculative')::float > c.threshold
--         OR (think_result ->> 'overall')::float > c.threshold
--     );

WITH config AS (
    SELECT 0.7 AS threshold
)
SELECT
    id,
    array_to_string(ARRAY_REMOVE(ARRAY[
        CASE WHEN (think_result ->> 'framing')::float > c.threshold
             THEN (think_result ->> 'framing_reason') END,
        CASE WHEN (think_result ->> 'clickbait')::float > c.threshold
             THEN (think_result ->> 'clickbait_reason') END,
        CASE WHEN (think_result ->> 'persuasive')::float > c.threshold
             THEN (think_result ->> 'persuasive_reason') END,
        CASE WHEN (think_result ->> 'hyper_stimulus')::float > c.threshold
             THEN (think_result ->> 'hyper_stimulus_reason') END,
        CASE WHEN (think_result ->> 'speculative')::float > c.threshold
             THEN (think_result ->> 'speculative_reason') END,
        CASE WHEN (think_result ->> 'overall')::float > c.threshold
             THEN (think_result ->> 'overall_reason') END
    ], NULL), E'\n') AS reasons_text -- E'\n' puts each reason on a new line
FROM
    items,
    config c
WHERE
    think_result IS NOT NULL
    AND (
        (think_result ->> 'framing')::float > c.threshold
        OR (think_result ->> 'clickbait')::float > c.threshold
        OR (think_result ->> 'persuasive')::float > c.threshold
        OR (think_result ->> 'hyper_stimulus')::float > c.threshold
        OR (think_result ->> 'speculative')::float > c.threshold
        OR (think_result ->> 'overall')::float > c.threshold
    );
