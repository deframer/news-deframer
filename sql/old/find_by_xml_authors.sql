WITH authors AS (
  SELECT
    items.id,
    feeds.root_domain AS domain,
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
  JOIN feeds ON feeds.id = items.feed_id
)
SELECT
  id,
  domain,
  first_author
FROM authors
WHERE first_author IS NOT NULL
  AND first_author <> '';
