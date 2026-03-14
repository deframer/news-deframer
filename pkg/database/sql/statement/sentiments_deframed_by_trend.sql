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
    ROUND(AVG((t.sentiments_deframed->>'v')::double precision)::numeric, 2) AS valence,
    ROUND(AVG((t.sentiments_deframed->>'a')::double precision)::numeric, 2) AS arousal,
    ROUND(AVG((t.sentiments_deframed->>'d')::double precision)::numeric, 2) AS dominance,
    ROUND(MAX((t.sentiments_deframed->>'j')::double precision)::numeric, 2) AS joy,
    ROUND(MAX((t.sentiments_deframed->>'a_n')::double precision)::numeric, 2) AS anger,
    ROUND(MAX((t.sentiments_deframed->>'s')::double precision)::numeric, 2) AS sadness,
    ROUND(MAX((t.sentiments_deframed->>'f')::double precision)::numeric, 2) AS fear,
    ROUND(MAX((t.sentiments_deframed->>'d_g')::double precision)::numeric, 2) AS disgust
FROM public.trends t
WHERE
    LOWER(CAST(@term AS text)) = ANY(t.noun_stems)
    AND t.root_domain = @domain
    AND t.pub_date >= COALESCE(
        CAST(@date AS timestamp) - ((GREATEST(CAST(@days AS INTEGER), 1) - 1) * INTERVAL '1 DAY'),
        NOW() - (GREATEST(CAST(@days AS INTEGER), 1) * INTERVAL '1 DAY')
    )
    AND t.pub_date < COALESCE(
        CAST(@date AS timestamp) + INTERVAL '1 DAY',
        NOW()
    )
    AND t.sentiments_deframed <> '{}'::jsonb;
