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
),
unnested_verbs AS (
    SELECT UNNEST(verb_stems) AS verb
    FROM filtered_docs
    WHERE verb_stems IS NOT NULL
),
aggregated_verbs AS (
    SELECT
        verb,
        COUNT(*) AS verb_count
    FROM unnested_verbs
    WHERE verb IS NOT NULL AND verb <> ''
    GROUP BY 1
),
total_verb_counts AS (
    SELECT SUM(verb_count) AS total_verb_count FROM aggregated_verbs
),
ranked_verbs AS (
    SELECT
        ROW_NUMBER() OVER (ORDER BY verb_count DESC, verb ASC) AS rank,
        verb,
        verb_count,
        CASE WHEN total.total_verb_count > 0
             THEN verb_count::FLOAT / total.total_verb_count
             ELSE NULL
        END AS verb_rel_score
    FROM aggregated_verbs, total_verb_counts total
)
SELECT
    p.domain AS domain,
    COALESCE(n.rank, v.rank) AS rank,
    n.noun,
    n.noun_count,
    n.noun_rel_score,
    v.verb,
    v.verb_count,
    v.verb_rel_score
FROM ranked_nouns n
FULL OUTER JOIN ranked_verbs v ON n.rank = v.rank
CROSS JOIN params p
WHERE COALESCE(n.rank, v.rank) <= p.result_limit
ORDER BY rank;
