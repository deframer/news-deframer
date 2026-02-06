SET duckdb.force_execution = true;

SET SESSION "vars.last" = '30 days';
SET SESSION "vars.target_word" = 'trump';
SET SESSION "vars.language" = 'de';
SET SESSION "vars.domain" = '';

-- Parameters
WITH params AS (
    SELECT
        current_date - current_setting('vars.last')::INTERVAL AS start_date,
        current_date AS end_date,
        current_setting('vars.target_word')::VARCHAR AS target_word,
        NULLIF(current_setting('vars.domain'), '')::VARCHAR AS domain,
        current_setting('vars.language')::VARCHAR AS language
),
filtered_docs AS (
    SELECT td.*
    FROM trends td
    CROSS JOIN params p
    WHERE td.pub_date IS NOT NULL
      AND td.pub_date >= p.start_date
      AND td.pub_date < p.end_date + INTERVAL '1 day'
      AND (p.domain IS NULL OR td.root_domain = p.domain)
      AND (p.language IS NULL OR td.language = p.language)
),
daily_counts AS (
    SELECT
        DATE_TRUNC('day', pub_date) AS day,
        unnested.noun
    FROM filtered_docs
    CROSS JOIN UNNEST(noun_stems) AS unnested(noun)
    WHERE unnested.noun IS NOT NULL AND unnested.noun <> ''
),
per_day AS (
    SELECT
        day,
        SUM(CASE WHEN noun = (SELECT target_word FROM params) THEN 1 ELSE 0 END) AS target_count,
        COUNT(*) AS total_noun_count
    FROM daily_counts
    GROUP BY 1
)
SELECT
    day,
    CASE WHEN total_noun_count > 0
         THEN target_count::FLOAT / total_noun_count
         ELSE 0
    END AS relative_score
FROM per_day
ORDER BY day;
