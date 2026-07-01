WITH config AS (
    SELECT
        14 AS days
), base AS (
    SELECT
        f.root_domain,
        i.id,
        i.created_at,
        i.updated_at,
        i.pub_date,
        i.think_result,
        i.think_error,
        i.think_error_count
    FROM feeds f
    JOIN items i ON i.feed_id = f.id
    CROSS JOIN config c
    WHERE
        f.deleted_at IS NULL
        AND f.enabled = true
        AND f.polling = true
        AND f.mining = true
        AND f.root_domain IS NOT NULL
        AND f.root_domain <> ''
        AND i.created_at >= NOW() - (c.days * INTERVAL '1 DAY')
), processing_times AS (
    SELECT
        b.root_domain,
        EXTRACT(EPOCH FROM (b.updated_at - b.created_at)) / 60.0 AS processing_minutes
    FROM base b
    JOIN public.trends t ON t.item_id = b.id
    WHERE
        b.think_result IS NOT NULL
        AND b.think_result <> '{}'::jsonb
        AND b.think_error IS NULL
        AND b.think_error_count = 0
), base_metrics AS (
    SELECT
        b.root_domain,
        COUNT(b.id) AS item_count,
        ROUND(
            100.0 * COUNT(*) FILTER (
                WHERE b.created_at - b.pub_date <= INTERVAL '24 HOURS'
            ) / NULLIF(COUNT(b.id), 0),
            2
        ) AS delta_pub_date_lt_24h,
        ROUND(
            100.0 * COUNT(*) FILTER (
                WHERE b.created_at - b.pub_date > INTERVAL '24 HOURS'
            ) / NULLIF(COUNT(b.id), 0),
            2
        ) AS delta_pub_date_gt_24h,
        ROUND(
            100.0 * COUNT(*) FILTER (
                WHERE b.created_at::date = b.pub_date::date
            ) / NULLIF(COUNT(b.id), 0),
            2
        ) AS delta_pub_date_same_day
    FROM base b
    GROUP BY b.root_domain
), processing_metrics AS (
    SELECT
        pt.root_domain,
        ROUND(
            percentile_cont(0.25) WITHIN GROUP (ORDER BY pt.processing_minutes)::numeric,
            2
        ) AS proc_min_p25,
        ROUND(
            percentile_cont(0.50) WITHIN GROUP (ORDER BY pt.processing_minutes)::numeric,
            2
        ) AS proc_min_p50,
        ROUND(
            percentile_cont(0.75) WITHIN GROUP (ORDER BY pt.processing_minutes)::numeric,
            2
        ) AS proc_min_p75,
        ROUND(
            percentile_cont(0.99) WITHIN GROUP (ORDER BY pt.processing_minutes)::numeric,
            2
        ) AS proc_min_p99,
        ROUND(
            MAX(pt.processing_minutes)::numeric,
            2
        ) AS proc_min_p100
    FROM processing_times pt
    GROUP BY pt.root_domain
), final_rows AS (
    SELECT
        bm.root_domain,
        bm.item_count::numeric AS item_count,
        pm.proc_min_p25,
        pm.proc_min_p50,
        pm.proc_min_p75,
        pm.proc_min_p99,
        pm.proc_min_p100,
        bm.delta_pub_date_lt_24h,
        bm.delta_pub_date_gt_24h,
        bm.delta_pub_date_same_day
    FROM base_metrics bm
    LEFT JOIN processing_metrics pm ON pm.root_domain = bm.root_domain
)
SELECT *
FROM (
    SELECT
        fr.root_domain,
        fr.item_count,
        fr.proc_min_p25,
        fr.proc_min_p50,
        fr.proc_min_p75,
        fr.proc_min_p99,
        fr.proc_min_p100,
        fr.delta_pub_date_lt_24h,
        fr.delta_pub_date_gt_24h,
        fr.delta_pub_date_same_day
    FROM final_rows fr

    UNION ALL

    SELECT
        'Overall' AS root_domain,
        ROUND(AVG(fr.item_count), 2) AS item_count,
        ROUND(AVG(fr.proc_min_p25), 2) AS proc_min_p25,
        ROUND(AVG(fr.proc_min_p50), 2) AS proc_min_p50,
        ROUND(AVG(fr.proc_min_p75), 2) AS proc_min_p75,
        ROUND(AVG(fr.proc_min_p99), 2) AS proc_min_p99,
        ROUND(AVG(fr.proc_min_p100), 2) AS proc_min_p100,
        ROUND(AVG(fr.delta_pub_date_lt_24h), 2) AS delta_pub_date_lt_24h,
        ROUND(AVG(fr.delta_pub_date_gt_24h), 2) AS delta_pub_date_gt_24h,
        ROUND(AVG(fr.delta_pub_date_same_day), 2) AS delta_pub_date_same_day
    FROM final_rows fr
) results
ORDER BY
    CASE WHEN results.root_domain = 'Overall' THEN 1 ELSE 0 END,
    results.root_domain ASC;
