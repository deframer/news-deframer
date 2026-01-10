# Specification: News Deframer System

## 1. Test and Helper Systems

These systems provide the controlled environment required to develop and test the core Deframer logic without relying on unstable external networks or rate-limited APIs.

### 1.1 RSS-Bridge
- **Purpose**:
  - Serves as the "Connector" logic.
  - Generates standard RSS Feeds for websites that lack them.
  - Sanitizes/Standardizes broken or limited upstream feeds before the Deframer touches them.
- **Technology**:
  - PHP Service (Community Standard).
  - Fetches `index.html`, parses via XPath/CSS Selectors.
  - **Constraint**: Does not render JavaScript (server-side scraping only).
- **Network**:
  - **Internal URL**: `http://rssbridge` (Accessible by Deframer).
  - **External URL**: `http://localhost:8002` (Accessible by Developer).

### 1.2 Dummy News Site
- **Purpose**:
  - Acts as the "Upstream Source."
  - Provides a controllable environment to test edge cases (broken images, paywalls, weird encodings).
  - Allows injection of specific test articles via API.
- **Technology**:
  - WordPress Stack.
  - Theme: `ColorMag` (Simulates a standard magazine layout).
- **Network**:
  - **Internal URL**: `http://wordpress` (Accessible by Deframer).
  - **External URL**: `http://localhost:8003` (Accessible by Developer).

### 1.3 FreshRSS
- **Purpose**:
  - Simple Web based RSS Reader
  - Can read internal and external RSS Feeds
- **Technology**:
  - Docker
- **Network**:
  - **External URL**: `http://localhost:8001` (Accessible by Developer).

---

## 2. News Deframer Architecture

The system is designed as a **Producer-Consumer** architecture. It decouples the high-throughput HTTP handling (API) from the latency-heavy AI processing (Worker).

### 2.1 Core Components
1.  **API Gateway (Service)**:
    - Written in **Golang**.
    - Stateless HTTP Server.
    - **Capacity**: Engineered to handle **10k concurrent users**.
    - Reads from **Valkey** (Hot Cache) and **PostgreSQL** (Cold Store/Source of Truth).
2.  **Deframer Worker (Service)**:
    - Written in **Golang**.
    - Background processor split into logical roles: **Download Worker**, **AI Worker**, and **Feed Builder**.
    - Scales horizontally based on Queue Depth.
3.  **Valkey (Infrastructure)**:
    - **Hot Cache**: Stores the ready-to-serve XML strings (`Final Feed Cache`).
    - **Message Queue**: Handles tasks (`IngestQueue`, `FeedQueue`).
    - **State Machine**: Handles Distributed Locks and Pending Counters.
4.  **PostgreSQL (Infrastructure)**:
    - **Persistent Storage**: Stores configuration (`Feeds`), processed data (`Items`), and cold cache (`cached_feeds`).
    - Ensures data survival across container restarts.

---

## 3. The RSS Deframer Proxy (API)

### Purpose
To serve valid RSS 2.0 documents instantly. It acts as the "Read Model" of the system.

### Endpoints

#### A. The Feed Proxy
```bash
GET /?url=${ENCODED_URL}&embedded=true&max_score=0.5
```
**Behavior (The Fallback Chain)**:
1.  **Resolution**: Map `url` to `uuid` using Valkey `feed_url:{encoded_url}`.
2.  **Hot Cache Hit**: If Valkey contains `feed:{uuid}:final`, return XML immediately (< 50ms).
3.  **Cold Cache Hit**: If Valkey is empty, query PostgreSQL `cached_feeds` table. If found, return XML and populate Valkey.
4.  **Miss (New/Empty Feed)**:
    - Register feed in DB if new.
    - **Push Feed UUID** to `IngestQueue`.
    - **Immediate Response**: Serve an "Empty" RSS Feed (Valid Header, No Items) to prevent blocking.

#### B. The JSON Lookup
```bash
GET /api/lookup?link=${ARTICLE_URL}
```
**Behavior**:
- Used to check if a specific article URL has already been deframed.
- **Status Codes**:
    - `200 OK`: Item found and deframed. Returns JSON object.
    - `202 Accepted`: Item found but currently processing.
    - `404 Not Found`: Domain or Link unknown.

---

## 4. The Background Worker

### Purpose
To handle the heavy lifting: Polling, Scraping, AI Processing, and XML Reconstruction.

### Routines & Locking

#### 1. The Download Worker (Ingest)
- **Triggers**:
  - **Cron**: Every 15 minutes. Iterates feeds where **`enabled=true` AND `polling=true`**.
  - **On-Demand**: Consumes `IngestQueue` (triggered by API Cache Miss).
- **Concurrency Control**:
  - Uses Valkey `SET key value NX EX 900` on the **Feed UUID**.
  - **Deduplication**: If triggered by Queue but lock is held, the task is discarded (prevents duplicate processing).
- **Logic**:
  - Fetches upstream XML.
  - **Validate Domain**: Checks if the Item URL matches the Feed's domain (if `EnforceFeedDomain` is true).
  - Diffs content against DB.
  - Saves new raw items to Postgres (`items` table) with `AnalyzerResult` initialized to empty JSON.
  - Sets Valkey `Pending Counter` = `Count(New Items)`.
  - Pushes Item SHA256s to the specific `FeedQueue`.

#### 2. The AI Worker (Processor)
- **Trigger**: Consumes `FeedQueue`.
- **Concurrency Control**:
  - Acquires **Item Lock** (SHA256) with strict **5-minute TTL**.
- **Logic**:
  - Executes the **Deframing Algorithm** (via Interface).
  - **Success**: Updates Postgres `AnalyzerResult`, decrements `Pending Counter`.
  - **Failure/Timeout**: Logs error, decrements `Pending Counter` (effectively drops item to prevent blocking).
  - **Retry**: Optional/TBD (Must avoid infinite loops).

#### 3. The Feed Builder (Aggregator)
- **Trigger**: Runs when `Pending Counter` for a feed reaches `0`.
- **Logic**:
  - Loads `cached_feeds` structure (XML Header) from Postgres.
  - Fetches all valid items (with populated `AnalyzerResult`) from Postgres.
  - Updates `pubDate` and `UpdatedAt`.
  - Builds the Final XML.
  - **Write 1**: Updates Valkey `Final Feed Cache`.
  - **Write 2**: Updates Postgres `cached_feeds` (Cold Store).

---

## 5. The Hybrid-State Algorithm

This algorithm ensures the system handles Bootstraps, Reboots, and Updates without ever serving an empty response, utilizing a Hot/Cold storage strategy.

### Phase 1: Ingest (The Hybrid/Pending View)
When new items arrive, they are stored in Postgres. The `Pending Counter` in Valkey prevents partial builds. The API continues to serve the *previous* successful build (from Hot or Cold cache) until processing is complete.

### Phase 2: Completion (The Upgrade)
Once the `Feed Builder` runs (all items processed):
1.  The XML is regenerated including the new Deframed items.
2.  Valkey and Postgres Caches are updated.
3.  **Result**: On the next refresh, the user sees the new content.

### Phase 3: Bootstrap (Cold Start)
If the system restarts and Valkey is empty:
1.  API checks Postgres `cached_feeds`.
2.  If data exists, it is served immediately and re-cached in Valkey.
3.  If no data exists, an Empty Feed is served and a download is triggered.

---

## 6. The Deframing Algorithm (AI Service)

This component is an abstraction layer over Large Language Models (LLMs).

### Interface
- **Input**: `feeds.RssItem` (from `gorilla.feed` package).
- **Output**: `JsonDocument` (Structured extraction/`AnalyzerResult`).

### Implementation Details
- **Mock Mode**: For development, a Mock Implementation simulates latency (`sleep(5)`) and returns dummy data to save costs.
- **Future Scalability**: The AI Worker is an interface. In the future, this can be replaced by an external decentralized network where multiple AIs process items in parallel and vote on the best result.
- **Prompts**: Configurable system prompts embedded in the application.

---

## 7. Data Schema

Based on the provided Go GORM models and Valkey implementation.

### Database (PostgreSQL)

**Table: `feeds`**
*Embeds `Base` (ID, CreatedAt, UpdatedAt, DeletedAt)*
- `id`: UUID (PK, default `uuid_generate_v4()`)
- `url`: String (**Indexed**)
- `enforce_feed_domain`: Boolean (Default: `true`. Items must come from the feed domain).
- `enabled`: Boolean (Default: `false`, **Indexed**).
- `polling`: Boolean (Default: `false`. If true, background worker actively refreshes content).
- `deleted_at`: Timestamp (Soft Delete support via GORM).

**Table: `items`**
*Does NOT embed `Base`*
- `id`: UUID (PK)
- `created_at`: Timestamp (Default: `now()`)
- `updated_at`: Timestamp (Default: `now()`)
- `hash`: String (64 chars, **Unique Indexes**: `idx_hash_feed_url`, `idx_hash_feed`).
- `feed_id`: UUID (**Indexed**, FK to Feeds, **Unique Indexes**: `idx_feed_url`, `idx_hash_feed_url`, `idx_hash_feed`).
- `url`: String (**Indexed**, **Unique Indexes**: `idx_feed_url`, `idx_hash_feed_url`).
- `analyzer_result`: JSONB (The processed AI content, **NOT NULL**).
- `content`: Text (The raw item content, **NOT NULL**).
- `min_hash`: String (Optional, Default: `null`).
- **Constraints**:
  - `FeedID` + `URL` is **UNIQUE**.
  - `URL` alone is **NOT UNIQUE** (One URL can appear in multiple Feeds).

**Table: `cached_feeds` (Cold Store)**
- `id`: UUID (PK & FK to Feed.ID).
- `created_at`: Timestamp (Default: `now()`)
- `updated_at`: Timestamp (Default: `now()`).
- `xml_header`: Text (**NOT NULL**).
- `item_refs`: Text Array (`text[]`, Default: `'{}'`, Stores item SHA256 refs).

### Cache (Valkey)

- **Key**: `feed_url:{url_coded_url}`
  - **Purpose**: Resolution of URL to UUID.
  - **Value**: JSON `FeedUrlToUUID` `{ "cache": int, "uuid": uuid }`.
- **Key**: `feed_uuid:{uuid}`
  - **Purpose**: Reverse lookup and Metadata.
  - **Value**: JSON `FeedInfo` `{ "cache": int, "base_domain": [], "url": string, "polling": bool }`.
- **Key**: `feed:{uuid}:final` -> Final XML String (Ready to serve).
- **Key**: `feed:{uuid}:pending` -> Integer (Counter of active AI tasks).
- **Key**: `queue:ingest` -> List of Feed UUIDs.
- **Key**: `queue:feed:{uuid}` -> List of Item SHA256s.

---

## 8. Scaling and Configuration

### Configuration Management
- **API Keys & Secrets**:
  - Injected via Environment Variables (`AI_API_KEY`, `AI_TYPE`, `DB_PASSWORD`).
  - Managed via Kubernetes Secrets or `.env` file for local dev.
- **Feed Subscription List**:
  - Managed via the **Database**.
  - Allows dynamic addition of feeds without restarting containers.
- **Worker Capacity**:
  - Workers track `current_jobs < MAX_CAPACITY` via local variables to determine availability.

### Scaling Strategy
- **Users (10k concurrent)**:
  - Handled by the **API Service**.
  - Scales on Memory/CPU usage.
  - Since 99% of requests hit Valkey (Hot Cache), a single small instance can handle thousands of req/sec.
- **New Networks/Feeds**:
  - Handled by the **Worker Service**.
  - Scales on `ProcessingQueue` (FeedQueue) length.
  - If AI API is slow or load increases, add more Worker Replicas.

### Deployment
- **Local**: `docker-compose` spins up 1 API, 1 Worker, 1 Valkey, 1 DB.
- **Kubernetes**:
  - `Deployment` for API (behind LoadBalancer).
  - `Deployment` for Worker (autoscaled via KEDA or HPA based on Queue depth).
  - `StatefulSet` for Valkey/DB (or Cloud Managed Services).

---

## 9. Implementation Notes & Hints

- **URL Ambiguity & Security**:
  - A single URL (e.g., `example.com/foo`) can legitimately appear in multiple feeds (Syndication).
  - **Risk**: A malicious RSS feed could theoretically "claim" popular URLs to inject bad data into the system.
  - **Mitigation**: The `EnforceFeedDomain` boolean (default `true`) ensures items are rejected if they do not match the Feed's domain base.
  - **TODO**: We might need a more granular allow-list (Regex/Glob) for items if we encounter feeds that legitimately host items on 3rd party domains.
  - **Lookup**: When querying by URL, the application must handle multiple results (e.g., sort by latest).
- **Trigger Mechanism**: The Webserver pushes the `Feed UUID` to a Redis List. The Worker uses a `SETNX` lock to ensure it doesn't process the same feed twice if the queue accumulates duplicates.
- **Future AI Scalability**: We intend to replace the local AI worker with an external decentralized network where multiple AIs process the item in parallel and "vote" on the best result before returning it.
- **Rate Limits**: **TOTAL TBD** (Needs further investigation). Limits should likely be applied at the Queue consumption level.
- **GORM**: The ORM should handle Item versioning and audit trails efficiently.
- **MinHash**: Calculation method (HTML vs Text DOM) is TBD.
- **Mock Implementation**: Use a mock for AI Worker development to simulate latency and avoid costs.
