SET duckdb.force_execution = true;

SET SESSION "vars.last" = '30 days';
SET SESSION "vars.target_word" = 'usa';
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
),
raw_scores AS (
    SELECT
        day,
        target_count,
        total_noun_count,
        CASE WHEN total_noun_count > 0
             THEN target_count::FLOAT / total_noun_count
             ELSE 0
        END AS relative_score
    FROM per_day
)
SELECT
    day,
    -- target_count AS target_count,
    -- total_noun_count AS total_noun_count,
    -- relative_score,
    -- CASE WHEN MAX(relative_score) OVER () = 0 THEN 0
    --      ELSE relative_score / MAX(relative_score) OVER ()
    -- END AS scaled,
    CASE WHEN MAX(relative_score) OVER () = 0 THEN 0
         ELSE ROUND((relative_score / MAX(relative_score) OVER ()) * 100)::INT
    END AS scaled_percent
FROM raw_scores
ORDER BY day;
