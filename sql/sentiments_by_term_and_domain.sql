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
    COUNT(*) AS article_count,
    ROUND(AVG((t.sentiments->>'v')::double precision)::numeric, 2) AS valence,
    ROUND(AVG((t.sentiments->>'a')::double precision)::numeric, 2) AS arousal,
    ROUND(AVG((t.sentiments->>'d')::double precision)::numeric, 2) AS dominance,
    ROUND(MAX((t.sentiments->>'j')::double precision)::numeric, 2) AS joy,
    ROUND(MAX((t.sentiments->>'a_n')::double precision)::numeric, 2) AS anger,
    ROUND(MAX((t.sentiments->>'s')::double precision)::numeric, 2) AS sadness,
    ROUND(MAX((t.sentiments->>'f')::double precision)::numeric, 2) AS fear,
    ROUND(MAX((t.sentiments->>'d_g')::double precision)::numeric, 2) AS disgust
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
