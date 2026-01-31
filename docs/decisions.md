# Decisions

## Why RSS?

Why did we chose RSS as primary input?

- **Standardized Input**: RSS feeds are an universal data source.
- **Flexible Sources**: You can consume existing feeds directly from publishers or employ `rssbridge` to generate feeds from any HTML website.
- **Optimized Content**: We provide `rssbridge` configuration examples. Generating custom feeds via scraping is often superior to official feeds, as it enables the exclusion of paywalled or irrelevant content.

## No External Cache

**Decision:** We do **not** deploy Memcached, Redis, or Valkey at this stage.

**Reasoning:**
*   **Complexity:** Adding a container adds orchestration overhead.
*   **Performance:** Postgres `cached_feeds` (lookup by PK) is extremely fast.
*   **Scale:** For $<100$ req/sec, in-process memory is faster and simpler than network calls to a cache service.

**Implementation:**

ETag Middleware (All GET Requests)


## Scaling

**Context** The current architecture processes feeds sequentially within a worker. A worker performs the following steps:
1. Locks a feed.
2. Fetches RSS items.
3. **Synchronously processes** each item via the LLM ("Thinking").
4. Saves to the database.
5. Unlocks the feed.

Since LLM calls are slow (2–30 seconds per item), this effectively holds the worker "hostage," preventing it from syncing other feeds. The proposal was to offload the "Thinking" phase to a [Temporal.io](https://temporal.io) workflow to unblock the feed syncer.

- For the current scale (< 1,000 feeds), horizontal scaling of simple Go workers is sufficient and cost-effective.

