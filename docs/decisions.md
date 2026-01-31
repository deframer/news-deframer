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
