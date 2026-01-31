# Architecture: Trend Mining Service

## 1. Executive Summary

This document outlines the architecture for the **Trend Mining Service**, a new component designed to broaden user perspectives by identifying "blind spots" and cross-feed trends.

Currently, users are confined to their subscribed feeds ("Filter Bubbles"). By analyzing what other feeds are covering, we can surface relevant content (local or global events) that a user might be missing.

**Core Concept:**
This system implements the Trend Mining methodologies described in the PhD Thesis **"Trend Mining im Social Web"** by Olga Streibel. It uses the semantic categories provided by feeds and items to build a global topic cluster, identifying "bursty" topics across the entire ecosystem.

---

## 2. Architectural Strategy

We are introducing a **Decoupled Worker Architecture**.

*   **The Syncer (Go):** Remains responsible for high-speed RSS fetching and database syncing.
*   **The Miner (Python):** A new, independent service dedicated to trend analysis.

### Why Decouple?
Trend analysis is computationally expensive and requires a different lifecycle than fetching.
1.  **Independent Schedules:** We can fetch news every 15 minutes but mine for trends every hour.
2.  **Non-Blocking:** The expensive "Mining" process never blocks the "Syncing" process. New news arrives instantly; trends are calculated asynchronously.
3.  **Tooling:** We use Python to leverage the rich Data Science ecosystem (DuckDB, Pandas, Scikit-learn) required for the Streibel algorithm, while keeping the core application in high-performance Go.

---

## 3. The Implementation Model

### A. Dual Scheduling System
We do not use a queue. Instead, we use the **Postgres-as-a-Queue** pattern, but we introduce a *second* schedule.

*   **Schedule 1 (Syncing):** Controls when to fetch RSS (`next_run_at`).
*   **Schedule 2 (Mining):** Controls when to analyze trends (`next_mining_at`).

This allows a single Feed to be "Syncing" (Fetching new items) and "Mining" (Analyzing old items) simultaneously without race conditions.

### B. The Miner Workflow (Python Worker)

The Miner runs as a daemon loop:

1.  **Locking (Feed Level):**
    *   It asks the DB: *"Give me a feed that needs mining (`next_mining_at <= NOW()`) and isn't currently being mined."*
    *   It locks this feed for mining (`mining_locked_until`).
    *   *Note:* This lock is distinct from the Syncer lock. Both workers can operate on the same feed ID at the same time.

2.  **Fetching (Item Level):**
    *   It loads all items for this feed that haven't been mined yet.
    *   Criteria: `mining_done_at IS NULL`.

3.  **Thinking (The Algorithm):**
    *   The items (Title, Description, Categories) are fed into the **Streibel Algorithm**.
    *   **Concept Unification:** The algorithm handles fuzzy matching. It is responsible for understanding that "Tech", "Technology", and "IT" belong to the same cluster. We do *not* normalize this data before ingestion.
    *   **Global Context:** The Miner maintains a persistent index (DuckDB) of *all* processed trends. This allows it to correlate the current feed's items with trends found in other feeds processed minutes or hours ago.

4.  **Committing:**
    *   It marks the items as done: `UPDATE items SET mining_done_at = NOW()`.
    *   It reschedules the feed: `UPDATE feed_schedules SET next_mining_at = NOW() + 1 HOUR`.

---

## 4. Data Structure Changes

To support this, we modify the existing schema:

### 1. `items` Table
*   **`mining_done_at` (Timestamp):** Indicates when the Miner successfully processed this item. If `NULL`, the item is pending analysis.
*   **`categories` (Array of Strings):** Stores the raw category tags from the RSS feed (e.g., `["Politics", "US Election"]`).

### 2. `feed_schedules` Table
*   **`next_mining_at` (Timestamp):** When the Miner should next look at this feed.
*   **`mining_locked_until` (Timestamp):** Acts as the mutex for the Mining Worker.
*   **`mining_error` (Text):** Stores the last error message from the Python Miner for debugging.

---

## 5. Summary of Responsibilities

| Responsibility | Component | Notes |
| :--- | :--- | :--- |
| **Fetching RSS** | `Go Syncer` | Fast, High Concurrency |
| **LLM Analysis** | `Go Syncer` | "Thinking" about content quality |
| **Trend Clustering** | `Python Miner` | "Thinking" about global topics |
| **Concept Unification** | `Streibel Algo` | Handles "Tech" vs "Technology" mapping |
| **State Management** | `PostgreSQL` | The single source of truth |
| **Trend Index** | `DuckDB` | The Miner's working memory (Persisted) |
