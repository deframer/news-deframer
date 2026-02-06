SET duckdb.force_execution = true;

SET SESSION "vars.last" = '30 days';
SET SESSION "vars.top" = 5;

--SET SESSION "vars.domain_a" = 'spiegel.de';
--SET SESSION "vars.domain_b" = 'nius.de';
SET SESSION "vars.domain_a" = 'nius.de';
SET SESSION "vars.domain_b" = 'tagesschau.de';
--SET SESSION "vars.domain_a" = 'bbc.com';
--SET SESSION "vars.domain_b" = 'nytimes.com';

-- Parameters
WITH params AS (
    SELECT
        current_date - current_setting('vars.last')::INTERVAL AS start_date,
        current_date AS end_date,
        current_setting('vars.domain_a')::VARCHAR AS root_domain_a,
        current_setting('vars.domain_b')::VARCHAR AS root_domain_b,
        current_setting('vars.top')::INT AS result_limit
),
filtered_docs AS (
    SELECT td.*, p.root_domain_a, p.root_domain_b
    FROM trends td
    CROSS JOIN params p
    WHERE td.pub_date IS NOT NULL
      AND td.pub_date >= p.start_date
      AND td.pub_date < p.end_date + INTERVAL '1 day'
),
unnested_nouns AS (
    SELECT
        td.root_domain,
        unnested.noun
    FROM filtered_docs td
    CROSS JOIN UNNEST(noun_stems) AS unnested(noun)
    WHERE unnested.noun IS NOT NULL AND unnested.noun <> ''
),
ranked_nouns_a AS (
    SELECT
        noun,
        COUNT(*) AS noun_count,
        ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC, noun ASC) AS rank
    FROM unnested_nouns, params p
    WHERE root_domain = p.root_domain_a
    GROUP BY noun
),
ranked_nouns_b AS (
    SELECT
        noun,
        COUNT(*) AS noun_count,
        ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC, noun ASC) AS rank
    FROM unnested_nouns, params p
    WHERE root_domain = p.root_domain_b
    GROUP BY noun
),
set_nouns_a AS (
    SELECT noun, noun_count FROM ranked_nouns_a WHERE rank <= (SELECT result_limit * 10 FROM params)
),
set_nouns_b AS (
    SELECT noun, noun_count FROM ranked_nouns_b WHERE rank <= (SELECT result_limit * 10 FROM params)
),
intersection_nouns AS (
    SELECT
        a.noun,
        a.noun_count AS count_noun_a,
        b.noun_count AS count_noun_b,
        ROW_NUMBER() OVER (ORDER BY a.noun_count + b.noun_count DESC, a.noun) AS rn
    FROM set_nouns_a a
    INNER JOIN set_nouns_b b ON a.noun = b.noun
),
blindspots_nouns_a AS (
    SELECT
        a.noun,
        a.noun_count,
        ROW_NUMBER() OVER (ORDER BY a.noun_count DESC, a.noun) AS rn
    FROM set_nouns_a a
    LEFT JOIN set_nouns_b b ON a.noun = b.noun
    WHERE b.noun IS NULL
),
blindspots_nouns_b AS (
    SELECT
        b.noun,
        b.noun_count,
        ROW_NUMBER() OVER (ORDER BY b.noun_count DESC, b.noun) AS rn
    FROM set_nouns_b b
    LEFT JOIN set_nouns_a a ON a.noun = b.noun
    WHERE a.noun IS NULL
),
unioned_results AS (
    SELECT
        'intersection' AS category,
        n.noun,
        n.count_noun_a,
        n.count_noun_b
    FROM intersection_nouns n
    WHERE n.rn <= (SELECT result_limit FROM params)
    UNION ALL
    SELECT
        'blindspots_a_has' AS category,
        n.noun,
        n.noun_count,
        NULL
    FROM blindspots_nouns_a n
    WHERE n.rn <= (SELECT result_limit FROM params)
    UNION ALL
    SELECT
        'blindspots_b_has' AS category,
        n.noun,
        NULL,
        n.noun_count
    FROM blindspots_nouns_b n
    WHERE n.rn <= (SELECT result_limit FROM params)
)
SELECT * FROM unioned_results
ORDER BY category, noun;
