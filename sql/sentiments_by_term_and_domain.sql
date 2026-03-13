SELECT
    COUNT(*) AS count,
    ROUND((percentile_cont(0.5) WITHIN GROUP (ORDER BY (t.sentiments->>'v')::double precision))::numeric, 2) AS valence,
    ROUND((percentile_cont(0.5) WITHIN GROUP (ORDER BY (t.sentiments->>'a')::double precision))::numeric, 2) AS arousal,
    ROUND((percentile_cont(0.5) WITHIN GROUP (ORDER BY (t.sentiments->>'d')::double precision))::numeric, 2) AS dominance,
    ROUND((percentile_cont(0.5) WITHIN GROUP (ORDER BY (t.sentiments->>'j')::double precision))::numeric, 2) AS joy,
    ROUND((percentile_cont(0.5) WITHIN GROUP (ORDER BY (t.sentiments->>'a_n')::double precision))::numeric, 2) AS anger,
    ROUND((percentile_cont(0.5) WITHIN GROUP (ORDER BY (t.sentiments->>'s')::double precision))::numeric, 2) AS sadness,
    ROUND((percentile_cont(0.5) WITHIN GROUP (ORDER BY (t.sentiments->>'f')::double precision))::numeric, 2) AS fear,
    ROUND((percentile_cont(0.5) WITHIN GROUP (ORDER BY (t.sentiments->>'d_g')::double precision))::numeric, 2) AS disgust
FROM public.trends t
WHERE
    -- Trigger Selection ($T$)
    'trump' = ANY(t.noun_stems)

    -- Domain/Scope Filters
    AND t."language" = 'de'
    --AND t.root_domain = 'apollo-news.net'
    AND t.root_domain = 'tagesschau.de'

    -- Time Window ($TW$)
    AND t.pub_date >= NOW() - INTERVAL '90 DAYS'

    -- Sentiment Filter
    AND t.sentiments <> '{}'::jsonb;
