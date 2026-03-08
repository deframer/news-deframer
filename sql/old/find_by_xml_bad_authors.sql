WITH authors AS (
  SELECT
    feeds.root_domain AS domain,
    items.think_rating,
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
), split_authors AS (
  SELECT
    domain,
    think_rating,
    BTRIM(
      regexp_split_to_table(
        regexp_replace(
          regexp_replace(first_author, '\s*(?:&amp;|&)\s*', ',', 'gi'),
          '\s+(?:and|und)\s+',
          ',',
          'gi'
        ),
        '\s*(?:,|/)\s*'
      )
    ) AS author
  FROM authors
), cleaned_authors AS (
  SELECT
    domain,
    think_rating,
    BTRIM(
      regexp_replace(
        regexp_replace(author, '\s*\|.*$', '', 'g'),
        '\s+',
        ' ',
        'g'
      )
    ) AS author
  FROM split_authors
)
SELECT
  AVG(think_rating) AS rating,
  COUNT(*) AS article_count,
  domain,
  author
FROM cleaned_authors
WHERE domain IS NOT NULL
  AND domain <> ''
  AND author IS NOT NULL
  AND author <> ''
GROUP BY domain, author
ORDER BY domain ASC, rating DESC, author ASC;
