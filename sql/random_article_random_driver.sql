WITH config AS (
    SELECT
        'spiegel.de' AS driving_domain,
        8 AS period_hours,
        'de' AS language,
        3 AS associate_limit,
        2 AS min_shared_nouns
), driver AS (
    SELECT
        i.id,
        f.root_domain AS domain,
        i.url,
        i.think_result,
        t.noun_stems,
        i.pub_date
    FROM items i
    JOIN feeds f ON f.id = i.feed_id
    JOIN public.trends t ON t.item_id = i.id
    CROSS JOIN config s
    WHERE f.root_domain = s.driving_domain
      AND i.language = s.language
      AND i.pub_date >= NOW() - (s.period_hours * INTERVAL '1 HOUR')
      AND f.enabled = true
      AND f.deleted_at IS NULL
      AND f.root_domain IS NOT NULL
      AND f.root_domain <> ''
    ORDER BY random()
    LIMIT 1
), driver_terms AS (
    SELECT
        d.id AS driver_id,
        d.domain AS driver_domain,
        d.url AS driver_url,
        d.pub_date AS driver_pub_date,
        'NOUN' AS stem_type,
        unnest(d.noun_stems) AS stem,
        1 AS weight
    FROM driver d
), associate_candidates AS (
    SELECT
        i.id,
        f.root_domain AS domain,
        i.url,
        i.think_result,
        t.noun_stems,
        i.pub_date
    FROM items i
    JOIN feeds f ON f.id = i.feed_id
    JOIN public.trends t ON t.item_id = i.id
    CROSS JOIN config s
    WHERE i.language = s.language
      AND i.pub_date >= NOW() - (s.period_hours * INTERVAL '1 HOUR')
      AND f.root_domain IS NOT NULL
      AND f.root_domain <> ''
      AND f.root_domain <> s.driving_domain
      AND f.enabled = true
      AND f.deleted_at IS NULL
      AND NOT EXISTS (
          SELECT 1
          FROM driver d
          WHERE d.id = i.id
      )
), candidate_terms AS (
    SELECT
        c.id,
        c.domain,
        c.url,
        c.think_result,
        c.pub_date,
        'NOUN' AS stem_type,
        unnest(c.noun_stems) AS stem,
        1 AS weight
    FROM associate_candidates c
), scored_candidates AS (
    SELECT
        c.id,
        c.domain,
        c.url,
        c.think_result,
        c.pub_date,
        SUM(d.weight) AS score,
        COUNT(DISTINCT d.stem) AS shared_stem_count,
        STRING_AGG(DISTINCT d.stem, ', ' ORDER BY d.stem) AS shared_stems
    FROM candidate_terms c
    JOIN driver_terms d
      ON d.stem = c.stem
     AND d.stem_type = c.stem_type
    WHERE c.stem IS NOT NULL
      AND c.stem <> ''
    GROUP BY
        c.id,
        c.domain,
        c.url,
        c.think_result,
        c.pub_date
), ranked_associates AS (
    SELECT
        *,
        ROW_NUMBER() OVER (
            PARTITION BY domain
            ORDER BY score DESC, shared_stem_count DESC, pub_date DESC, url ASC
        ) AS domain_rank
    FROM scored_candidates
), top_associates AS (
    SELECT
        domain,
        url,
        think_result
    FROM ranked_associates
    CROSS JOIN config s
    WHERE domain_rank = 1
      AND shared_stem_count >= s.min_shared_nouns
    ORDER BY score DESC, shared_stem_count DESC, pub_date DESC, domain ASC
    LIMIT (SELECT associate_limit FROM config)
), combined AS (
    SELECT
        domain,
        'driver' AS type,
        json_build_object(
            'title_original', driver.think_result ->> 'title_original',
            'description_original', driver.think_result ->> 'description_original',
            'title_corrected', driver.think_result ->> 'title_corrected',
            'description_corrected', driver.think_result ->> 'description_corrected',
            'category', driver.think_result ->> 'category'
        ) AS content,
        url
    FROM driver
    UNION ALL
    SELECT
        domain,
        'associate' AS type,
        json_build_object(
            'title_original', think_result ->> 'title_original',
            'description_original', think_result ->> 'description_original',
            'title_corrected', think_result ->> 'title_corrected',
            'description_corrected', think_result ->> 'description_corrected',
            'category', think_result ->> 'category'
        ) AS content,
        url
    FROM top_associates
)
SELECT
    domain,
    type,
    content,
    url
FROM combined
ORDER BY CASE type WHEN 'driver' THEN 0 ELSE 1 END, domain ASC;
