WITH config AS (
    SELECT
        'de' AS language,
        'de' AS country,
        'politics' as category,
        8 AS period_hours,
        3 AS associate_limit,
        2 AS min_shared_nouns
), driving_domains AS (
    SELECT
        COALESCE(ARRAY_AGG(DISTINCT f.root_domain), ARRAY[]::text[]) AS driving_domains
    FROM feeds f
    CROSS JOIN config s
    WHERE f.language = s.language
      AND f.country = s.country
      AND f.tags @> ARRAY['lead']::text[]
      AND f.enabled = true
      AND f.deleted_at IS NULL
      AND f.root_domain IS NOT NULL
      AND f.root_domain <> ''
), driver_candidates AS (
    SELECT
        i.id,
        f.root_domain AS domain,
        i.url,
        i.think_result,
        t.noun_stems,
        i.pub_date
    FROM config s
    CROSS JOIN driving_domains dd
    JOIN items i ON i.language = s.language
    JOIN feeds f ON f.id = i.feed_id
    JOIN public.trends t ON t.item_id = i.id
    WHERE f.root_domain = ANY(dd.driving_domains)
      AND f.country = s.country
      AND LOWER(COALESCE(i.think_result ->> 'category', '')) = LOWER(s.category)
      AND i.pub_date >= NOW() - (s.period_hours * INTERVAL '1 HOUR')
      AND f.enabled = true
      AND f.deleted_at IS NULL
      AND f.root_domain IS NOT NULL
      AND f.root_domain <> ''
), associate_universe AS (
    SELECT
        i.id,
        f.root_domain AS domain,
        i.url,
        i.think_result,
        t.noun_stems,
        i.pub_date
    FROM config s
    CROSS JOIN driving_domains dd
    JOIN items i ON i.language = s.language
    JOIN feeds f ON f.id = i.feed_id
    JOIN public.trends t ON t.item_id = i.id
    WHERE i.pub_date >= NOW() - (s.period_hours * INTERVAL '1 HOUR')
      AND f.country = s.country
      AND LOWER(COALESCE(i.think_result ->> 'category', '')) = LOWER(s.category)
      AND f.root_domain IS NOT NULL
      AND f.root_domain <> ''
      AND NOT (f.root_domain = ANY(dd.driving_domains))
      AND f.enabled = true
      AND f.deleted_at IS NULL
), driver_candidate_terms AS (
    SELECT
        d.id AS driver_id,
        d.domain AS driver_domain,
        unnest(d.noun_stems) AS stem
    FROM driver_candidates d
), associate_terms AS (
    SELECT
        a.id,
        a.domain,
        a.url,
        a.think_result,
        a.pub_date,
        unnest(a.noun_stems) AS stem
    FROM associate_universe a
), driver_associate_scores AS (
    SELECT
        d.driver_id,
        a.id,
        a.domain,
        a.url,
        a.think_result,
        a.pub_date,
        COUNT(DISTINCT d.stem) AS shared_stem_count
    FROM driver_candidate_terms d
    JOIN associate_terms a
      ON a.stem = d.stem
    WHERE d.stem IS NOT NULL
      AND d.stem <> ''
      AND a.id <> d.driver_id
      AND a.domain <> d.driver_domain
    GROUP BY
        d.driver_id,
        a.id,
        a.domain,
        a.url,
        a.think_result,
        a.pub_date
), ranked_driver_associates AS (
    SELECT
        ds.*,
        ROW_NUMBER() OVER (
            PARTITION BY ds.driver_id, ds.domain
            ORDER BY ds.shared_stem_count DESC, ds.pub_date DESC, ds.url ASC
        ) AS domain_rank
    FROM driver_associate_scores ds
), driver_strength AS (
    SELECT
        r.driver_id,
        COUNT(*) AS strong_associate_count
    FROM ranked_driver_associates r
    CROSS JOIN config s
    WHERE r.domain_rank = 1
      AND r.shared_stem_count >= s.min_shared_nouns
    GROUP BY r.driver_id
), driver_candidates_ranked AS (
    SELECT
        dc.id,
        dc.domain,
        dc.url,
        dc.think_result,
        dc.noun_stems,
        dc.pub_date,
        COALESCE(ds.strong_associate_count, 0) AS strong_associate_count
    FROM driver_candidates dc
    LEFT JOIN driver_strength ds ON ds.driver_id = dc.id
), driver AS (
    SELECT
        id,
        domain,
        url,
        think_result,
        noun_stems,
        pub_date
    FROM driver_candidates_ranked
    ORDER BY
        CASE
            WHEN strong_associate_count >= (SELECT associate_limit FROM config) THEN 0
            ELSE 1
        END,
        CASE
            WHEN strong_associate_count >= (SELECT associate_limit FROM config) THEN random()
            ELSE 0
        END,
        strong_associate_count DESC,
        random()
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
        a.id,
        a.domain,
        a.url,
        a.think_result,
        a.noun_stems,
        a.pub_date
    FROM associate_universe a
    CROSS JOIN driver d
    WHERE a.domain <> d.domain
      AND a.id <> d.id
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
        COUNT(DISTINCT d.stem) AS shared_stem_count
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
        id,
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
        id,
        domain,
        url,
        think_result
    FROM top_associates
    WHERE driver_rank <= (SELECT associate_limit FROM config)
), associate_slots AS (
    SELECT
        GREATEST((SELECT associate_limit FROM config) - COUNT(*), 0) AS needed
    FROM limited_associates
), related_fallback_candidates AS (
    SELECT
        s.id,
        s.domain,
        s.url,
        s.think_result,
        s.pub_date,
        s.shared_stem_count
    FROM scored_candidates s
    WHERE s.shared_stem_count >= 1
      AND NOT EXISTS (
          SELECT 1
          FROM limited_associates la
          WHERE la.id = s.id
      )
), related_fallback_associates AS (
    SELECT
        rfc.id,
        rfc.domain,
        rfc.url,
        rfc.think_result
    FROM related_fallback_candidates rfc
    ORDER BY rfc.shared_stem_count DESC, rfc.pub_date DESC, rfc.domain ASC, rfc.url ASC
    LIMIT (SELECT needed FROM associate_slots)
), all_associates AS (
    SELECT id, domain, url, think_result FROM limited_associates
    UNION ALL
    SELECT id, domain, url, think_result FROM related_fallback_associates
), combined AS (
    SELECT
        domain,
        'driver' AS type,
        json_build_object(
            'title_original', d.think_result ->> 'title_original',
            'description_original', d.think_result ->> 'description_original',
            'title_corrected', d.think_result ->> 'title_corrected',
            'description_corrected', d.think_result ->> 'description_corrected',
            'category', d.think_result ->> 'category'
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
            'description_corrected', think_result ->> 'description_corrected',
            'category', think_result ->> 'category'
        ) AS content,
        url,
        (SELECT domain FROM driver) AS driver_domain
    FROM all_associates
)
SELECT
    domain,
    type,
    content,
    url
FROM combined
ORDER BY CASE type WHEN 'driver' THEN 0 ELSE 1 END, domain ASC;
