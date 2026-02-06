SET duckdb.force_execution = true;

SET SESSION "vars.last" = '30 days';
SET SESSION "vars.top" = 10;
--SET SESSION "vars.domain" = 'apollo-news.net';
SET SESSION "vars.domain" = 'spiegel.de';

-- Parameters
WITH params AS (
    SELECT
        current_date - current_setting('vars.last')::INTERVAL AS start_date,
        current_date AS end_date,
        NULLIF(current_setting('vars.domain'), '')::VARCHAR AS domain,
        current_setting('vars.top')::INT AS result_limit
),
filtered_docs AS (
    SELECT td.*, p.domain
    FROM trends td
    CROSS JOIN params p
    WHERE td.pub_date IS NOT NULL
      AND td.pub_date >= p.start_date
      AND td.pub_date < p.end_date + INTERVAL '1 day'
      AND (p.domain IS NULL OR td.root_domain = p.domain)
),
unnested_nouns AS (
    SELECT UNNEST(noun_stems) AS noun
    FROM filtered_docs
    WHERE noun_stems IS NOT NULL
),
aggregated_nouns AS (
    SELECT
        noun,
        COUNT(*) AS noun_count
    FROM unnested_nouns
    WHERE noun IS NOT NULL AND noun <> ''
    GROUP BY 1
),
total_noun_counts AS (
    SELECT SUM(noun_count) AS total_noun_count FROM aggregated_nouns
),
ranked_nouns AS (
    SELECT
        ROW_NUMBER() OVER (ORDER BY noun_count DESC, noun ASC) AS rank,
        noun,
        noun_count,
        CASE WHEN total.total_noun_count > 0
             THEN noun_count::FLOAT / total.total_noun_count
             ELSE NULL
        END AS noun_rel_score
    FROM aggregated_nouns, total_noun_counts total
)
SELECT
    p.domain AS domain,
    n.rank,
    n.noun,
    n.noun_count,
    n.noun_rel_score
FROM ranked_nouns n
CROSS JOIN params p
WHERE n.rank <= p.result_limit
ORDER BY rank;
