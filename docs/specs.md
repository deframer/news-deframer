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
    - Handles User Traffic (10k concurrent users).
    - Reads from **Valkey** (Cache) and **PostgreSQL** (Source of Truth).
2.  **Deframer Worker (Service)**:
    - Written in **Golang**.
    - Background processor.
    - Handles fetching upstream HTML, running AI algorithms, and updating the database.
    - Scales horizontally based on Queue Depth.
3.  **Valkey (Infrastructure)**:
    - **Hot Cache**: Stores the ready-to-serve XML strings.
    - **Message Queue**: Handles tasks (`IngestQueue`, `ProcessingQueue`).
4.  **PostgreSQL (Infrastructure)**:
    - **Persistent Storage**: Stores configuration (Feeds) and processed data (Items).
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
**Behavior**:
1.  **Cache Hit**: If Valkey contains `feed:{hash}`, return XML immediately (< 50ms).
2.  **Cache Miss (Unknown Feed)**:
    - Register feed in DB.
    - Push to `IngestQueue`.
    - **Fallback**: Proxy the original Upstream XML directly to the user (Pass-through) to prevent blocking.
3.  **Cache Miss (Known Feed / Reboot)**:
    - Query DB for existing items.
    - Perform **Hybrid Merge** (mix of processed and pending items).
    - Update Valkey and Serve.

#### B. The JSON Lookup (New Feature)
```bash
GET /api/lookup?link=${ARTICLE_URL}
```
**Behavior**:
- Used to check if a specific article URL has already been deframed.
- **Returns**: JSON Object representing the enhanced item (see Data Schema).
- **Status Codes**:
    - `200 OK`: Item found and deframed.
    - `202 Accepted`: Item found but currently processing.
    - `404 Not Found`: Domain or Link unknown.

---

## 4. The Background Worker

### Purpose
To handle the heavy lifting: Polling, Scraping, and AI Processing.

### Routines
1.  **The Poller (Cron)**:
    - Runs every 15 minutes.
    - Scans DB for feeds where `last_polled > 15m`.
    - Pushes feeds to `IngestQueue`.
2.  **The Ingester**:
    - Consumes `IngestQueue`.
    - Fetches upstream XML.
    - Identifies **new** items (deduplication via GUID/Link).
    - Saves new items to DB with status `PENDING`.
    - **Trigger**: Immediately calls `RegenerateCache()` so users see new headlines instantly (even if content isn't deframed yet).
    - Pushes new items to `ProcessingQueue`.
3.  **The Processor**:
    - Consumes `ProcessingQueue`.
    - Executes the **Deframing Algorithm**.
    - Updates DB to `COMPLETED`.
    - **Trigger**: Calls `RegenerateCache()` to upgrade the feed from "Hybrid" to "Deframed".

---

## 5. The Hybrid-State Algorithm

This algorithm ensures the system handles Bootstraps, Reboots, and Updates without ever serving an empty response.

### Phase 1: Ingest (The Hybrid View)
When new items arrive from the upstream RSS:
1.  We store them as `PENDING`.
2.  We generate the Output XML immediately.
    - **Old Items**: Rendered with full AI content (Summary, Score).
    - **New Items**: Rendered with original upstream content (Description only).
3.  **Result**: The user sees the news *now*.

### Phase 2: Processing (The Upgrade)
As the Worker processes the queue:
1.  Item status changes from `PENDING` -> `COMPLETED`.
2.  Valkey Cache is updated.
3.  **Result**: On the next refresh, the user sees the "New" items have transformed into "Deframed" items.

### Phase 3: Bootstrap (Cold Start)
If the Docker container restarts and Valkey is empty:
1.  API checks DB.
2.  DB returns the last 30 items.
3.  API generates the XML from DB data.
4.  Valkey is repopulated.
5.  **Result**: Zero external network calls required to recover service.

---

## 6. The Deframing Algorithm (AI Service)

This component is an abstraction layer over Large Language Models (LLMs). It handles the extraction and enhancement of content.

### Interface
- **Input**: `feeds.RssItem` (from `gorilla.feed` package).
- **Output**: `JsonDocument` (Structured extraction).

### Implementation Details
- **AI Provider**: The system is agnostic to the AI backend.
    - Supported types (e.g., OpenAI, Anthropic, Local LLM) are injected via `AI_TYPE` and `AI_API_KEY` environment variables.
- **Null Implementation**: A specific "Null" provider exists for debugging/local dev. It returns hardcoded mock data immediately without making external API calls, saving costs and latency during UI testing.
- **Prompts**:
    - The system embeds specific system prompts to guide the AI (e.g., "Extract the main content, ignore ads, summarize in 3 bullet points, assign a score").
    - Prompts are configurable but part of the application binary/config.
- **Behavior**:
    - The worker passes the `feeds.RssItem` (containing the Link/Description) to the AI Service.
    - The AI Service fetches the content (if necessary) and applies the embedded prompt.
    - The result is returned as a structured JSON.

---

## 7. Data Schema (Conceptual)

### Database (PostgreSQL)

**Table: `feeds`**
- `id`: UUID
- `url`: String (Upstream RSS URL)
- `config`: JSON (Interval, specific parsing rules)
- `last_polled`: Timestamp

**Table: `items`**
- `hash`: String (SHA256 of Link - Primary Key)
- `feed_id`: FK
- `upstream_xml`: Text (The raw `<item>` blob)
- `deframed_json`: JSONB (The AI result)
- `status`: Enum (`PENDING`, `COMPLETED`, `FAILED`)
- `created_at`: Timestamp

### Cache (Valkey)

- **Key**: `feed:{hash_of_feed_url}`
  - **Value**: String (Complete XML Document)
  - **TTL**: 15 minutes (Reset on update)

---

## 8. Scaling and Configuration

### Configuration Management
- **API Keys & Secrets**:
  - Injected via Environment Variables (`AI_API_KEY`, `AI_TYPE`, `DB_PASSWORD`).
  - Managed via Kubernetes Secrets or `.env` file for local dev.
- **Feed Subscription List**:
  - Managed via the **Database**.
  - Allows dynamic addition of feeds without restarting containers.

### Scaling Strategy
- **Users (10k concurrent)**:
  - Handled by the **API Service**.
  - Scales on Memory/CPU usage.
  - Since 99% of requests hit Valkey, a single small instance can handle thousands of req/sec.
- **New Networks/Feeds**:
  - Handled by the **Worker Service**.
  - Scales on `ProcessingQueue` length.
  - If AI API is slow, add more Worker Replicas.

### Deployment
- **Local**: `docker-compose` spins up 1 API, 1 Worker, 1 Valkey, 1 DB.
- **Kubernetes**:
  - `Deployment` for API (behind LoadBalancer).
  - `Deployment` for Worker (autoscaled via KEDA or HPA).
  - `StatefulSet` for Valkey/DB (or Cloud Managed Services).
