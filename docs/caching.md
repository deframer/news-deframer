# Caching Strategy

## 1. Overview
This document defines the caching strategy for the News Deframer API. The goal is to minimize database load and network bandwidth without introducing complex external infrastructure (like Redis/Memcached) for the current scale.

**Constraint:** The backend updates feeds every **10-15 minutes**.
**Strategy:** We leverage standard HTTP caching mechanisms (`ETag`, `Last-Modified`) and a lightweight in-memory application cache.

---

## 2. HTTP Layer Caching (Client-Side)

We use HTTP headers to instruct downstream clients (browsers, RSS readers, CDNs) to cache responses.

### A. RSS Feeds (`GET /rss`)
*   **Content:** Pre-rendered XML stored in the database (`cached_feeds`).
*   **Update Frequency:** ~10-15 minutes.
*   **Header Strategy:**
    ```http
    Cache-Control: public, max-age=600
    ```
    *   `public`: Allows shared caches (CDNs/Proxies) to store it.
    *   `max-age=600`: Content is considered fresh for **10 minutes**.

### B. API Endpoints (`GET /api/*`)
*   **Content:** Dynamic JSON aggregated from multiple feeds (`handleSite`, `handleDomains`).
*   **Update Frequency:** High (aggregate of all feed updates).
*   **Header Strategy:**
    ```http
    Cache-Control: public, max-age=300
    ```
    *   `max-age=300`: Content is considered fresh for **5 minutes**. This balances freshness with load reduction.

---

## 3. Application-Level Caching (Server-Side)

To protect the database from redundant queries, we implement caching inside the Go application.

### A. ETag Middleware (All GET Requests)
We implement a global middleware that hashes the response body and manages the `If-None-Match` flow.

**Logic:**
1.  Buffer the response body.
2.  Generate a SHA256 hash of the content.
3.  Set response header: `ETag: "hash-value"`.
4.  Check request header: `If-None-Match`.
5.  **Hit:** If hashes match $\rightarrow$ Return `304 Not Modified` (Empty Body).
6.  **Miss:** Return `200 OK` with body.

**Benefit:** drastically reduces bandwidth usage for polling clients (RSS readers).

### B. In-Memory Cache (Specific Endpoints)
For expensive aggregation queries (like `handleSite`), we use a local in-memory cache (e.g., `go-cache` or `hashicorp/golang-lru`).

*   **Target:** `GET /api/site` (Aggregates items by root domain).
*   **Key:** `site:<root_domain>:<max_score>`
*   **TTL:** **5 minutes** (Matches `max-age`).
*   **Why:** A single request to this endpoint can trigger a complex JOIN/DISTINCT query. Caching it in RAM eliminates database load for 5 minutes at a time.

---

## 4. Infrastructure Decision (No External Cache)

**Decision:** We do **not** deploy Memcached, Redis, or Valkey at this stage.

**Reasoning:**
*   **Complexity:** Adding a container adds orchestration overhead.
*   **Performance:** Postgres `cached_feeds` (lookup by PK) is extremely fast.
*   **Scale:** For $<100$ req/sec, in-process memory is faster and simpler than network calls to a cache service.
*   **Traefik:** Traefik OSS is a router, not a cache. We avoid complex middleware plugins.

**Future Scaling:**
If we scale to multiple API replicas or $>100$ req/sec, we will introduce **Valkey** to share the cache state across instances.
