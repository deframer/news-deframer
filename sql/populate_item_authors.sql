WITH extracted_authors AS (
  SELECT
    id,
    BTRIM(
      COALESCE(
        (regexp_match(content, '<dc:creator[^>]*>([^<]*)</dc:creator>', 'i'))[1],
        (regexp_match(content, '<atom:author[^>]*>(?:(?!</atom:author>).)*<atom:name[^>]*>([^<]*)</atom:name>(?:(?!</atom:author>).)*</atom:author>', 'in'))[1],
        (regexp_match(content, '<atom:name[^>]*>([^<]*)</atom:name>', 'i'))[1],
        (regexp_match(content, '<content:author[^>]*>([^<]*)</content:author>', 'i'))[1],
        (regexp_match(content, '<byline[^>]*>([^<]*)</byline>', 'i'))[1],
        (regexp_match(content, '<creator[^>]*>([^<]*)</creator>', 'i'))[1],
        (regexp_match(content, '<slash:author[^>]*>([^<]*)</slash:author>', 'i'))[1],
        (regexp_match(content, '<sy:author[^>]*>([^<]*)</sy:author>', 'i'))[1],
        (regexp_match(content, '<foaf:name[^>]*>([^<]*)</foaf:name>', 'i'))[1],
        (regexp_match(content, '<creativeCommons:attributionName[^>]*>([^<]*)</creativeCommons:attributionName>', 'i'))[1]
      )
    ) AS first_author
  FROM items
), split_authors AS (
  SELECT
    id,
    ord,
    BTRIM(author_part) AS author
  FROM extracted_authors,
  LATERAL unnest(
    regexp_split_to_array(
      regexp_replace(
        regexp_replace(first_author, '\s*(?:&amp;|&)\s*', ',', 'gi'),
        '\s+(?:and|und)\s+',
        ',',
        'gi'
      ),
      '\s*(?:,|/)\s*'
    )
  ) WITH ORDINALITY AS split(author_part, ord)
  WHERE first_author IS NOT NULL
    AND first_author <> ''
), cleaned_authors AS (
  SELECT
    id,
    ord,
    BTRIM(
      regexp_replace(
        regexp_replace(author, '\s*\|.*$', '', 'g'),
        '\s+',
        ' ',
        'g'
      )
    ) AS author
  FROM split_authors
), deduped_authors AS (
  SELECT DISTINCT ON (id, author)
    id,
    ord,
    author
  FROM cleaned_authors
  WHERE author IS NOT NULL
    AND author <> ''
  ORDER BY id, author, ord
), aggregated_authors AS (
  SELECT
    id,
    array_agg(author ORDER BY ord) AS authors
  FROM deduped_authors
  GROUP BY id
)
UPDATE items
SET authors = COALESCE(aggregated_authors.authors, '{}'::text[])
FROM aggregated_authors
WHERE items.id = aggregated_authors.id
  AND (items.authors IS NULL OR items.authors = '{}'::text[]);
