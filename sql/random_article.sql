WITH config AS (
    SELECT
        ARRAY['spiegel.de', 'tagesschau.de', 'stern.de', 'nzz.de'] AS driving_domains,
        'de' AS language,
        8 AS period_hours,
        3 AS associate_limit,
        2 AS min_shared_nouns
), driver_candidates AS (
    SELECT
        i.id,
        f.root_domain AS domain,
        i.url,
        i.think_result,
        t.noun_stems,
        i.pub_date
    FROM config s
    JOIN items i ON i.language = s.language
    JOIN feeds f ON f.id = i.feed_id
    JOIN public.trends t ON t.item_id = i.id
    WHERE f.root_domain = ANY(s.driving_domains)
      AND i.pub_date >= NOW() - (s.period_hours * INTERVAL '1 HOUR')
      AND f.enabled = true
      AND f.deleted_at IS NULL
      AND f.root_domain IS NOT NULL
      AND f.root_domain <> ''
), driver AS (
    SELECT
        id,
        domain,
        url,
        think_result,
        noun_stems,
        pub_date
    FROM driver_candidates
    ORDER BY random()
    LIMIT 1
), driver_terms AS (
    SELECT
        d.domain AS driving_domain,
        d.id AS driver_id,
        d.domain AS driver_domain,
        d.url AS driver_url,
        d.pub_date AS driver_pub_date,
        unnest(d.noun_stems) AS stem
    FROM driver d
), associate_candidates AS (
    SELECT
        i.id,
        f.root_domain AS domain,
        i.url,
        i.think_result,
        t.noun_stems,
        i.pub_date
    FROM driver d
    CROSS JOIN config s
    JOIN items i ON i.language = s.language
    JOIN feeds f ON f.id = i.feed_id
    JOIN public.trends t ON t.item_id = i.id
    WHERE i.pub_date >= NOW() - (s.period_hours * INTERVAL '1 HOUR')
      AND f.root_domain IS NOT NULL
      AND f.root_domain <> ''
      AND f.root_domain <> d.domain
      AND NOT (f.root_domain = ANY(s.driving_domains))
      AND f.enabled = true
      AND f.deleted_at IS NULL
      AND NOT EXISTS (
          SELECT 1
          FROM driver x
          WHERE x.id = i.id
      )
), candidate_terms AS (
    SELECT
        c.id,
        c.domain,
        c.url,
        c.think_result,
        c.pub_date,
        unnest(c.noun_stems) AS stem
    FROM associate_candidates c
), scored_candidates AS (
    SELECT
        c.id,
        c.domain,
        c.url,
        c.think_result,
        c.pub_date,
        COUNT(DISTINCT d.stem) AS shared_stem_count,
        STRING_AGG(DISTINCT d.stem, ', ' ORDER BY d.stem) AS shared_stems
    FROM candidate_terms c
    JOIN driver_terms d
      ON d.stem = c.stem
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
            ORDER BY shared_stem_count DESC, pub_date DESC, url ASC
        ) AS domain_rank
    FROM scored_candidates
), top_associates AS (
    SELECT
        domain,
        url,
        think_result,
        shared_stem_count,
        pub_date,
        ROW_NUMBER() OVER (
            ORDER BY shared_stem_count DESC, pub_date DESC, domain ASC
        ) AS driver_rank
    FROM ranked_associates
    WHERE domain_rank = 1
      AND shared_stem_count >= (SELECT min_shared_nouns FROM config)
), limited_associates AS (
    SELECT
        domain,
        url,
        think_result
    FROM top_associates
    WHERE driver_rank <= (SELECT associate_limit FROM config)
), combined AS (
    SELECT
        domain,
        'driver' AS type,
        json_build_object(
            'title_original', d.think_result ->> 'title_original',
            'description_original', d.think_result ->> 'description_original',
            'title_corrected', d.think_result ->> 'title_corrected',
            'description_corrected', d.think_result ->> 'description_corrected'
        ) AS content,
        url,
        d.domain AS driver_domain
    FROM driver d
    UNION ALL
    SELECT
        domain,
        'associate' AS type,
        json_build_object(
            'title_original', think_result ->> 'title_original',
            'description_original', think_result ->> 'description_original',
            'title_corrected', think_result ->> 'title_corrected',
            'description_corrected', think_result ->> 'description_corrected'
        ) AS content,
        url,
        (SELECT domain FROM driver) AS driver_domain
    FROM limited_associates
)
SELECT
    domain,
    type,
    content,
    url
FROM combined
ORDER BY CASE type WHEN 'driver' THEN 0 ELSE 1 END, domain ASC;
