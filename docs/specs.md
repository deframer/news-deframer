# Specification: News Deframer System

---

## 1. News Deframer Architecture

The system is designed as a **Producer-Consumer** architecture (Reader-Writer). It decouples the HTTP handling (API) from the heavy processing (Worker) using the Database as the synchronization point.

### 1.1 Core Components
1.  **API Gateway (Service)**:
    - Written in **Golang**.
    - Stateless HTTP Server.
    - **Role**: "Dumb Reader". It only serves pre-calculated XML from the DB.
    - Reads from **PostgreSQL** (`cached_feeds`).
2.  **Deframer Worker (Service)**:
    - Written in **Golang**.
    - Background Monolith.
    - **Role**: "The Writer". Handles Download, AI Processing, and Feed Building using a locking queue mechanism.
3.  **PostgreSQL (Infrastructure)**:
    - **Persistent Storage**: Stores configuration (`Feeds`) and data (`Items`).
    - **Queue State**: Stores scheduling and locking info (`feed_schedules`).
    - **Cache**: Stores the ready-to-serve XML strings (`cached_feeds`).
4.  **Admin CLI (Tool)**:
    - Written in **Golang**.
    - Interface for operators to manage feeds, system configuration, and revive dead feeds.

---

## 2. The RSS Deframer Proxy (API)

### Purpose
To serve valid RSS 2.0 documents instantly. It acts as the "Read Model" of the system.

### Endpoints

#### A. The Feed Proxy
```bash
GET /rss?url=${ENCODED_URL}
```
**Behavior**:
1.  **Resolution**: Map `url` to `uuid` using PostgreSQL `feeds` table.
2.  **Access Check**: Check if `feeds.enabled` is `true`.
    -   **False**: Return `403 Forbidden` or `404 Not Found` (depending on preference).
3.  **Cache Hit**: Query PostgreSQL `cached_feeds` table.
    -   **Found**: Return `xml_content` immediately.
    -   **Not Found / Empty**: Return `404 Not Found`.
4.  **Header Logic**: The HTTP `Last-Modified` header and the RSS `<lastBuildDate>` are derived strictly from `cached_feeds.updated_at`.

#### B. The JSON Lookup
```bash
GET /api/item?url=${ARTICLE_URL}&max_score=0.5
```
**Behavior**:
- Used to check if a specific article URL has already been deframed.
- **Status Codes**:
    - `200 OK`: Single Item. Returns JSON object.
    - `404 Not Found`: Domain or Link unknown.

```bash
GET /api/site?root=${ROOT_DOMAIN}
```
**Behavior**:
- Used to return a list of the most recent items we deframed for the site logic. This will be used by the Web Browser plugin to replace the portal page.
- **Resolution**:
    1.  Find all `feeds` where `feeds.root_domain == root`.
    2.  Select items from these feeds.
- **Logic**: Returns distinct items (by URL) using a "Fair Distribution" algorithm (interleaving items from different feeds) to prevent high-volume feeds from dominating.
- **Status Codes**:
    - `200 OK`: Returns JSON object.
    - `404 Not Found`: Domain unknown.

---

## 3. The Background Worker

### Purpose
To handle the heavy lifting: Polling, Scraping, AI Processing, and XML Reconstruction.

### Routines & Locking

#### 1. The Queue Logic
The worker polls the `feed_schedules` table to find available jobs.

#### 2. The Locking Strategy (The Lease)
-   **Acquire Job**:
    -   Worker queries for a feed where:
        1.  `feed_schedules.next_run_at <= NOW()`
        2.  `feed_schedules.locked_until` is NULL (or expired).
        3.  **`feeds.enabled` is `true`**. (Worker ignores scheduled jobs if the feed is disabled).
    -   **Action**: Sets `locked_until = NOW() + 10 minutes`. This acts as a "Lease".

#### 3. The Processing Sequence (Monolith)
Once a lock is acquired, the worker executes the following sequential steps:
1.  **Fetch Feed**: Download upstream XML.
    -   **Domain Check**: If `feeds.enforce_feed_domain` is true, discard items where `item.Host != feed.Host`.
2.  **Identify New Items**:
    -   Calculate SHA256 Hash of the Item URL.
    -   Check if `(feed_id, hash)` exists in `items` table.
    -   *Constraint Check*: This check is always scoped to the `feed_id`.
3.  **AI Processing Loop**:
    -   For each *new* item (where `feed_id` + `hash` is new):
        -   Download HTML content.
        -   **Run AI Deframer** (Extract facts/summary).
        -   Insert result into `items` table.
4.  **Feed Build**:
    -   Select last $N$ items for this `feed_id`.
    -   Generate final RSS XML.
    -   Collect the list of Item Hashes into `item_refs`.
    -   **Commit**: UPDATE/INSERT `cached_feeds` with `xml_content` and `item_refs`.

#### 4. Completion & Rescheduling
After the feed is built, the worker determines if it should run again.

-   **On Success**:
    -   **If `feeds.polling == true`**: Set `next_run_at = NOW() + Interval` (Re-arm the timer).
    -   **If `feeds.polling == false`**: Set `next_run_at = NULL` (Do not re-arm. This was a one-time sync).
    -   Clear `locked_until`.
    -   Clear `last_error`.

-   **On Error (Strike One)**:
    -   Worker writes the failure reason to `last_error` (e.g., "Timeout", "404").
    -   Worker sets `next_run_at` to **NULL**.
    -   Clear `locked_until`.
    -   **Result**: The feed is immediately ejected from the schedule regardless of the `polling` flag. It is considered "Dead" until manual intervention.

---

## 4. The Hybrid-State Algorithm

This algorithm ensures the system handles Bootstraps and Updates efficiently using the Database as a Cold Cache.

### Phase 1: Ingest
When new items arrive, they are processed and stored in the `items` table. This table allows duplication of content *across* feeds (same URL in two different feeds = two different rows), but enforces uniqueness *within* a feed.

### Phase 2: Completion
Once the Worker completes the AI processing, it updates `cached_feeds`. This table is the **Single Source of Truth** for the API regarding *what* is served and *when* it was last updated.

### Phase 3: Bootstrap (Cold Start)
If a feed is newly added:
1.  API returns 404/Empty.
2.  Worker picks up the feed via the Schedule (if triggered via Sync or Polling).
3.  Worker populates `cached_feeds`.
4.  API begins serving content.

---

## 5. The Deframing Algorithm (AI Service)

This component is an abstraction layer over Large Language Models (LLMs).

### Interface
- **Input**: `feeds.RssItem` (from `gorilla.feed` package).
- **Output**: `JsonDocument` (Structured extraction/`AnalyzerResult`).

### Implementation Details
- **Mock Mode**: For development, a Mock Implementation simulates latency and returns dummy data.
- **Prompts**: Configurable system prompts embedded in the application.

---

## 6. Data Schema

### Database (PostgreSQL)

**Table: `feeds` (Configuration)**
*Embeds `Base` (ID, CreatedAt, UpdatedAt, DeletedAt)*
- `id`: UUID (PK, default `uuid_generate_v4()`)
- `url`: String (**Indexed**)
- `root_domain`: String (**Indexed**).
  - *Logic*: Grouping identifier (e.g., `example.com`). Automatically populated from `url` on creation, or manually set via Admin.
- `enforce_feed_domain`: Boolean (Default: `true`). Software enforcement flag).
- `polling`: Boolean (Default: `false`).
  - *Logic*: Controls re-arming. If `true`, the worker schedules the next run after success. If `false`, the worker runs once and sets `next_run_at = NULL`.
- `enabled`: Boolean (Default: `false`, **Indexed**).
  - *Logic*: Gatekeeper. If `false`, the API denies access and the Worker ignores the feed (even if scheduled).
- `deleted_at`: Timestamp.

**Table: `cached_feeds` (Source of Truth for API)**
*Sidecar table, 1:1 with Feeds*
- `id`: UUID (PK & FK to Feed.ID).
- `created_at`: Timestamp (Default: `now()`)
- `updated_at`: Timestamp (Default: `now()`).
- `xml_header`: Text (**NOT NULL**).
- `item_refs`: Text Array (`text[]`, Stores list of Item SHA256 hashes currently in the XML).
- `xml_content`: Text (The final ready-to-serve RSS XML string).

**Table: `feed_schedules` (Operational State)**
*Sidecar table, 1:1 with Feeds*
- `id`: UUID (PK, FK to Feeds).
- `next_run_at`: Timestamp (**Indexed**, Nullable).
  - *Logic*: If `NULL`, the feed is Dead/Paused. If `Timestamp <= NOW()`, it is due.
- `locked_until`: Timestamp (Nullable).
  - *Logic*: The Lease. If `Timestamp > NOW()`, a worker is busy.
- `last_error`: Text (Nullable).
  - *Logic*: Contains the reason for the "One Strike" death.

**Table: `items` (Content)**
*Does NOT embed `Base`*
- `id`: UUID (PK)
- `created_at`: Timestamp (Default: `now()`)
- `updated_at`: Timestamp (Default: `now()`)
- `hash`: String (64 chars, SHA256 of URL).
- `feed_id`: UUID (**Indexed**, FK to Feeds).
- `url`: String (**Indexed**).
- `analyzer_result`: JSONB (The processed AI content, **NOT NULL**).
- `content`: Text (The raw item content, **NOT NULL**).
- `pub_date`: Timestamp (pubDate of the XML, Default: `now()`).
- **Constraints**:
  - **Unique Index**: `idx_feed_id_hash` (`feed_id`, `hash`).
  - **Important Note**: The `url` (and `hash`) is **NOT unique globally**. A URL may appear in multiple feeds. It is only unique relative to the `feed_id`.

---

## 7. Scaling and Configuration

### Configuration Management
- **API Keys & Secrets**:
  - Injected via Environment Variables (`AI_API_KEY`, `AI_TYPE`, `DB_PASSWORD`).
- **Feed Subscription List**:
  - Managed via the **Database**.

### Scaling Strategy
- **Users**:
  - Handled by the **API Service** reading from Postgres.
- **New Networks/Feeds**:
  - Handled by the **Worker Service**. Multiple workers can run in parallel if they lock different Feeds.

### Deployment
- **Local**: `docker-compose` spins up 1 API, 1 Worker, 1 DB.

---

## 8. Implementation Notes

- **Optimization**: Ensure `cached_feeds` is TOASTed properly by Postgres as `xml_content` can be large.
- **Syndication**: Since `items` are scoped to `feed_id`, if two feeds syndicate the same article, the AI will process it twice (once for each feed). This is intentional to allow per-feed prompts or domain enforcement contexts.
- **Admin Sync**: Use the CLI to force updates during development.
- **Endpoints**: Probably also support `max_score` for the RSS endpoint. This is much harder as we need dynamic XML generating.

---

## 9. Admin CLI Tool

A command-line interface for operators to manage the system state and feed configurations without direct database manipulation.

**Usage**: `admin feed [command]`

-   **`add`**: Registers a new Feed URL. Creates the `feed` and `feed_schedule` entries.
    -   *Logic*: Automatically extracts the domain to populate `root_domain`.
    -   *Flag*: `--no-root-domain` prevents this automatic extraction, leaving the field empty.
-   **`delete`**: Soft-deletes a feed from Postgres.
-   **`enable`**: Sets `feeds.enabled=true`.
    -   *Logic*: Allows the API to serve content and the Worker to pick up the feed (if scheduled). Does *not* automatically schedule a run.
-   **`disable`**: Sets `feeds.enabled=false`.
    -   *Logic*: Immediately stops the API from serving content. Prevents the Worker from picking up the feed (even if `next_run_at` is pending). Does not clear the schedule, but renders it inert.
-   **`polling`**: Toggles `feeds.polling`.
    -   *Logic*: Controls re-arming. Enabling this ensures that *future* successful runs will schedule a follow-up run.
-   **`list`**: Displays a table of all feeds, including `enabled` status, `polling` status, `root_domain`, `next_run_at`, and `last_error`.
-   **`sync` (Resurrection)**:
    -   Forces a feed to run immediately.
    -   **Pre-requisite**: **Fails** if `feeds.enabled` is `false` or `feeds.deleted_at` is set.
    -   **Action**: Sets `next_run_at = NOW()`, clears `last_error`, and clears `locked_until`.
    -   *Note*: If `polling` is false, this acts as a one-time fetch. If `polling` is true, this kickstarts the continuous cycle.

**Usage**: `admin root-domain [command]`

-   **`list`**: Displays a list of all unique `root_domain` values currently registered in the system.
-   **`set`**: Manually sets the `root_domain` for a specific feed (by ID or URL).
-   **`clear`**: Clears the `root_domain` for a specific feed (by ID or URL).    -   **`set`**: Manually sets the `root_domain` for a specific feed (by ID or URL).
