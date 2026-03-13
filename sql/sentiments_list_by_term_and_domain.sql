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
SELECT
    p.term,
    p.root_domain,
    COUNT(*) AS count,
    ROUND((percentile_cont(0.5) WITHIN GROUP (ORDER BY (t.sentiments->>'v')::double precision))::numeric, 2) AS valence,
    ROUND((percentile_cont(0.5) WITHIN GROUP (ORDER BY (t.sentiments->>'a')::double precision))::numeric, 2) AS arousal,
    ROUND((percentile_cont(0.5) WITHIN GROUP (ORDER BY (t.sentiments->>'d')::double precision))::numeric, 2) AS dominance,
    ROUND((percentile_cont(0.5) WITHIN GROUP (ORDER BY (t.sentiments->>'j')::double precision))::numeric, 2) AS joy,
    ROUND((percentile_cont(0.5) WITHIN GROUP (ORDER BY (t.sentiments->>'a_n')::double precision))::numeric, 2) AS anger,
    ROUND((percentile_cont(0.5) WITHIN GROUP (ORDER BY (t.sentiments->>'s')::double precision))::numeric, 2) AS sadness,
    ROUND((percentile_cont(0.5) WITHIN GROUP (ORDER BY (t.sentiments->>'f')::double precision))::numeric, 2) AS fear,
    ROUND((percentile_cont(0.5) WITHIN GROUP (ORDER BY (t.sentiments->>'d_g')::double precision))::numeric, 2) AS disgust
FROM term_domain_pairs p
JOIN public.trends t
    ON p.term = ANY(t.noun_stems)
   AND p.root_domain = t.root_domain
WHERE
    t."language" = 'de'
    AND t.pub_date >= NOW() - INTERVAL '90 DAYS'
    AND t.sentiments <> '{}'::jsonb
GROUP BY
    p.term,
    p.root_domain
ORDER BY
    p.term,
    p.root_domain;
