SET duckdb.force_execution = true;

--SET SESSION "vars.day" = '2025-01-03';
SET SESSION "vars.day" = ''; -- today
SET SESSION "vars.top" = 10;
--SET SESSION "vars.domain" = 'apollo-news.net';
--SET SESSION "vars.domain" = 'spiegel.de';
SET SESSION "vars.domain" = '';
--SET SESSION "vars.language" = 'de';
SET SESSION "vars.language" = 'en';

-- Parameters
WITH params AS (
    SELECT
        COALESCE(NULLIF(current_setting('vars.day'), '')::DATE, current_date) AS target_date,
        NULLIF(current_setting('vars.domain'), '')::VARCHAR AS domain,
        NULLIF(current_setting('vars.language'), '')::VARCHAR AS language,
        current_setting('vars.top')::INT AS result_limit
),
filtered_docs AS (
    SELECT td.*, p.domain
    FROM trends td
    CROSS JOIN params p
    WHERE td.pub_date IS NOT NULL
      AND td.pub_date >= p.target_date
      AND td.pub_date < p.target_date + INTERVAL '1 day'
      AND (p.domain IS NULL OR td.root_domain = p.domain)
      AND (p.language IS NULL OR td.language = p.language)
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
