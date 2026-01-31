# Capacity Planning & Infrastructure Guide

This document outlines the resource requirements, scaling limits, and cost analysis for deploying News Deframer. It is designed to help you decide between cloud-based API usage and home-hosted hardware.

## 1. The Core Challenge: The "Thinking" Bottleneck

The primary performance constraint in this architecture is the **AI Processing ("Thinking")** phase.
While fetching RSS feeds takes milliseconds, analyzing an article with an LLM (Large Language Model) typically takes **3 to 10 seconds** per item.

Because of this high latency, a single worker processing feeds sequentially would fall hopelessly behind. To keep up with the influx of news, we must use **concurrent workers**.

---

## 2. Load Calculations

We assume the following baseline behavior for a news-heavy setup:
*   **Update Cycle:** Feeds are checked every **15 minutes**.
*   **Feed Velocity:** Average of **50 items** per feed, with **20%** being new or updated each cycle.
*   **Result:** 10 items to process per feed, every 15 minutes.

| Scenario | Total Feeds | Items / 15 min | Rate (Items/sec) | Daily Volume |
| :--- | :--- | :--- | :--- | :--- |
| **A: Enthusiast** | 100 | 1,000 | **~1.1** | 96,000 |
| **B: Power User** | 1,000 | 10,000 | **~11.1** | 960,000 |

---

## 3. Worker Scaling (How many do I need?)

We determine the required number of workers using [Little's Law](https://en.wikipedia.org/wiki/Little%27s_law):

$$ \text{Workers Needed} = \text{Arrival Rate (items/sec)} \times \text{Average Latency (seconds)} $$

Assuming an average AI processing time of **5 seconds**:

*   **Scenario A (100 Feeds):**
    $$ 1.1 \text{ items/sec} \times 5 \text{ sec} \approx \mathbf{6 \text{ Workers}} $$
    *(Safe Recommendation: Start with 10)*

*   **Scenario B (1,000 Feeds):**
    $$ 11.1 \text{ items/sec} \times 5 \text{ sec} \approx \mathbf{56 \text{ Workers}} $$
    *(Safe Recommendation: Start with 80-100)*

---

## 4. Cost Analysis: Cloud APIs (OpenAI / Gemini)

This model assumes you use a paid external provider.

**Assumptions:**
*   **Model:** Low-cost tier (GPT-4o-mini or Gemini Flash). High-end models (GPT-4o) are cost-prohibitive for this volume.
*   **Tokens:** ~1,500 tokens per feed item (1k context + 500 output).
*   **Price:** ~$0.15 - $0.30 per 1M tokens.

| Scenario | Daily Cost | Monthly Cost |
| :--- | :--- | :--- |
| **100 Feeds** | ~$25 | **~$750** |
| **1,000 Feeds** | ~$250 | **~$7,500** |

*Note: As you can see, even "cheap" models become expensive at 24/7 volume.*

---

## 5. Home Hosting Strategy (Mac Silicon)

For a "Home Data Center" approach, we replace the Cloud API with a local Mac (Mini or Studio) running a local LLM (e.g., Llama 3 8B via Ollama).

### A. Throughput Feasibility
Unlike the cloud, your local hardware has a hard throughput limit.
*   **Mac Mini / Studio (M-Series Pro/Max):** Can comfortably inference at **2–4 items/second** (batched or parallel requests) for 8B-parameter models.
*   **Scenario A (1.1 items/sec):** ✅ **Feasible.** One Mac can handle the load.
*   **Scenario B (11.1 items/sec):** ❌ **Not Feasible.** You would need a cluster of 3-4 Mac Studios.

### B. Return on Investment (ROI)
Let's compare buying hardware vs. paying the cloud bill for the **100 Feed** scenario.

#### Option 1: Mac Studio (High Performance)
*   **Hardware:** Mac Studio (M2 Max, 64GB RAM) $\approx$ **$2,000**.
*   **Electricity:** ~60W avg load @ $0.30/kWh $\approx$ **$13/month**.
*   **Cloud Cost:** **$750/month**.
*   **Break Even:** $\approx$ **2.7 Months**.

#### Option 2: Mac Mini M4 (Cheap)
The M4 chip has a significantly faster Neural Engine than the M2, making it an excellent budget choice.
*   **Hardware:** Mac Mini M4 (16GB RAM) $\approx$ **$700**.
*   **Electricity:** ~30W avg load $\approx$ **$7/month**.
*   **Cloud Cost:** **$750/month**.

$$ \text{Break Even Time} = \frac{\$700}{\$750 - \$7} \approx \mathbf{0.94 \text{ Months}} $$

**Verdict:** The **Mac Mini M4** is the optimal "Home Data Center" machine for this workload. It pays for itself in under 1 month.

### C. Network Considerations for Home Hosting
While the application's data usage is minimal, there are a few network factors to consider.

1.  **Data Usage & Bandwidth:**
    *   The application only downloads small RSS/XML feed files, not the full article webpages.
    *   Data consumption is very low, amounting to a few gigabytes per month at most.
    *   This will not strain your download bandwidth or put you at risk of exceeding typical ISP data caps.

2.  **Upstream (Upload) Stability:**
    *   Your Mac must upload the analysis results to your database.
    *   If your database is hosted in the cloud (e.g., AWS, Heroku), a stable, reliable internet connection is important.
    *   Brief internet outages are fine, but prolonged instability could delay your feeds from updating.
