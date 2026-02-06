SET duckdb.force_execution = true;

/*
  --------------------------------------------------------------------------------------
  TREND METRICS EXPLANATION
  Based on Streibel, O. (2013). "Trend Mining" (Dissertation), Chapter 10.
  --------------------------------------------------------------------------------------

  1. stem as trend_topic:
     The normalized word (Trigger/Noun or Indicator/Verb) being analyzed.
     Represents the 'Topic Area' or specific 'Actor' (if an Entity like 'Jeffrey Epstein')
     that is potentially emerging [3][4].

  2. frequency (Interest):
     The raw volume of mentions in the current time slice (e.g., today).
     Corresponds to the thesis concept of 'Interest': how much is this topic being
     talked about right now? [5].
     Note: High frequency alone does not equal a trend (common words like "Market"
     have high frequency but low trend scores).

  3. utility (Source Diversity):
     The count of DISTINCT feeds/domains discussing this topic.
     Thesis Definition (Def 10.2.4): Measures if the topic is useful/broad or just
     spam from a single source [6].
     - Low Utility (e.g., 1) = Echo chamber (only one source talking).
     - High Utility = Broad market consensus (a valid trend).

  4. outlier_ratio (Burstiness/Surprise):
     Calculated as: (Current Frequency) / (Historical Moving Average).
     Thesis Definition (Def 10.2.2): Identifies terms that appear "significantly often
     in a given time slice" compared to the baseline history [7].
     - Ratio > 1.0 : Growing attention.
     - Ratio > 1.5 : Significant Trend Signal (The "Outlier").
     - Ratio < 1.0 : Fading topic.
     NULL -> no historic data
  --------------------------------------------------------------------------------------
*/

WITH ranked_trends AS (
    SELECT
        stem as trend_topic,
        frequency,
        utility,
        outlier_ratio,
        time_slice,
        -- Calculate the rank here
        ROW_NUMBER() OVER (
            PARTITION BY stem
            ORDER BY outlier_ratio DESC
        ) as rn
    FROM view_trend_metrics_by_domain
    WHERE "language" = 'de'
      AND root_domain = 'tagesschau.de'
      AND stem_type = 'NOUN'
      --AND time_slice = DATE_TRUNC('day', NOW())  -- Or specific date
      AND time_slice >= NOW() - INTERVAL '7 DAYS'
      -- Note: Utility threshold might need to be lower for a single domain
      -- (e.g. > 1 means it appeared in at least 2 different RSS feeds on that site)
      AND utility >= 1
      AND outlier_ratio > 1.5
)
SELECT
    trend_topic,
    frequency,
    utility,
    outlier_ratio,
    time_slice
FROM ranked_trends
WHERE rn = 1 -- Filter for the top result here
ORDER BY outlier_ratio DESC NULLS LAST
LIMIT 20;