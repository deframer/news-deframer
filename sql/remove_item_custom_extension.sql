-- BEGIN TRANSACTION;

SELECT COUNT(*)
FROM items
WHERE content ~ '<deframer:(title_original|description_original)>';

-- WITH source AS (
--   SELECT
--     id,
--     COALESCE((regexp_match(content, '<deframer:title_original>(.*?)</deframer:title_original>'))[1], '') AS title_original,
--     COALESCE((regexp_match(content, '<deframer:description_original>(.*?)</deframer:description_original>'))[1], '') AS description_original
--   FROM items
--   WHERE content ~ '<deframer:(title_original|description_original)>'
-- )
-- UPDATE items i
-- SET content = regexp_replace(
--   regexp_replace(
--     regexp_replace(
--       regexp_replace(
--         regexp_replace(
--           i.content,
--           '<title>.*?</title>',
--           '<title>' || replace(replace(source.title_original, E'\\', E'\\\\'), '&', E'\\&') || '</title>',
--           's'
--         ),
--         '<description>.*?</description>',
--         '<description>' || replace(replace(source.description_original, E'\\', E'\\\\'), '&', E'\\&') || '</description>',
--         's'
--       ),
--       '<deframer:title_original>.*?</deframer:title_original>',
--       '',
--       'gs'
--     ),
--     '<deframer:description_original>.*?</deframer:description_original>',
--     '',
--     'gs'
--   ),
--   '[[:space:]]xmlns:deframer="[^"]+"',
--   '',
--   'g'
-- )
-- FROM source
-- WHERE i.id = source.id;

-- ROLLBACK;
-- COMMIT;
