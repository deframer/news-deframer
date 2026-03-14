WITH term_domain_pairs(term, root_domain) AS (
    VALUES
        ('trump', 'tagesschau.de'),
        ('afd', 'tagesschau.de'),
        ('union', 'tagesschau.de'),
        ('ukraine', 'tagesschau.de'),
        ('kanzler', 'tagesschau.de'),
        ('merz', 'tagesschau.de'),
        ('selenskyj', 'tagesschau.de'),
        ('epstein', 'tagesschau.de'),
        ('krieg', 'tagesschau.de'),
        ('zdf', 'tagesschau.de'),
        ('grüne', 'tagesschau.de'),

        ('trump', 'apollo-news.net'),
        ('afd', 'apollo-news.net'),
        ('union', 'apollo-news.net'),
        ('ukraine', 'apollo-news.net'),
        ('kanzler', 'apollo-news.net'),
        ('merz', 'apollo-news.net'),
        ('selenskyj', 'apollo-news.net'),
        ('epstein', 'apollo-news.net'),
        ('krieg', 'apollo-news.net'),
        ('zdf', 'apollo-news.net'),
        ('grüne', 'apollo-news.net'),

        ('trump', 'nius.de'),
        ('afd', 'nius.de'),
        ('union', 'nius.de'),
        ('ukraine', 'nius.de'),
        ('kanzler', 'nius.de'),
        ('merz', 'nius.de'),
        ('selenskyj', 'nius.de'),
        ('epstein', 'nius.de'),
        ('krieg', 'nius.de'),
        ('zdf', 'nius.de'),
        ('grüne', 'nius.de')
)
-- Aggregation rationale:
-- - Median is avoided because emotionally salient words can appear as statistical
--   outliers while still carrying the main affective signal of the text.
-- - VAD (valence, arousal, dominance) uses AVG because these are continuous
--   mood dimensions and the mean better reflects overall emotional tone.
-- - BE5 (joy, anger, sadness, fear, disgust) uses MAX because discrete emotion
--   spikes are often the signal of interest and should not be diluted by
--   neutral terms in longer texts.
-- - Depending on the analysis goal, AVG is also possible for BE5 if you want
--   average emotional tone instead of peak emotional intensity.
-- - Stop-word removal should happen upstream before these scores are computed.
SELECT
    p.term,
    p.root_domain,
    COUNT(*) AS article_count,
    ROUND(AVG((t.sentiments->>'v')::double precision)::numeric, 2) AS valence,
    ROUND(AVG((t.sentiments->>'a')::double precision)::numeric, 2) AS arousal,
    ROUND(AVG((t.sentiments->>'d')::double precision)::numeric, 2) AS dominance,
    ROUND(MAX((t.sentiments->>'j')::double precision)::numeric, 2) AS joy,
    ROUND(MAX((t.sentiments->>'a_n')::double precision)::numeric, 2) AS anger,
    ROUND(MAX((t.sentiments->>'s')::double precision)::numeric, 2) AS sadness,
    ROUND(MAX((t.sentiments->>'f')::double precision)::numeric, 2) AS fear,
    ROUND(MAX((t.sentiments->>'d_g')::double precision)::numeric, 2) AS disgust
FROM term_domain_pairs p
JOIN public.trends t
    ON p.term = ANY(t.noun_stems)
   AND p.root_domain = t.root_domain
WHERE
    t."language" = 'de'
    AND t.pub_date >= NOW() - INTERVAL '7 DAYS'
    AND t.sentiments <> '{}'::jsonb
GROUP BY
    p.term,
    p.root_domain
ORDER BY
    p.term,
    p.root_domain;
