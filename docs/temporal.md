# Architecture Analysis: Introducing Temporal.io for LLM Orchestration

## Context
The current architecture processes feeds sequentially within a worker. A worker performs the following steps:
1. Locks a feed.
2. Fetches RSS items.
3. **Synchronously processes** each item via the LLM ("Thinking").
4. Saves to the database.
5. Unlocks the feed.

Since LLM calls are slow (2–30 seconds per item), this effectively holds the worker "hostage," preventing it from syncing other feeds. The proposal was to offload the "Thinking" phase to a [Temporal.io](https://temporal.io) workflow to unblock the feed syncer.

## Trade-off Analysis

### What we would GAIN (Pros)
1.  **Massive Parallelism (Unblocking the Syncer):**
    The Syncer becomes a fast "producer," dispatching tasks to Temporal and immediately releasing the feed lock. This allows a single Syncer to handle thousands of feeds while the heavy lifting happens asynchronously in the background.
2.  **Durable Retries:**
    Temporal provides out-of-the-box durability. If the LLM provider (OpenAI/Gemini) is down or returning 503s, Temporal queues the tasks and retries with exponential backoff for hours or days if necessary, without custom logic.
3.  **Global Rate Limiting:**
    Temporal allows configuring global rate limits (e.g., "Max 500 requests/minute to OpenAI") shared across all distributed workers. This is difficult to achieve with independent Go binaries without a separate distributed store like Redis.
4.  **Observability:**
    Instant visibility into the state of every single item (Pending, Retrying, Failed) via the Temporal UI.

### What we would LOSE (Cons)
1.  **Operational Complexity (The "Tax"):**
    Introducing Temporal requires running a Temporal Server, a persistence database (Cassandra/Postgres), and potentially ElasticSearch for visibility. For a self-hosted project, this is significant infrastructure overhead.
2.  **Development Complexity:**
    Requires adopting the Temporal SDK patterns (Activities, Workflows) and managing "Workflow Versioning" when prompt logic changes.
3.  **Data Consistency:**
    State is now split between the application database (Postgres) and the Temporal history. We must ensure reliable callbacks/signaling to update the primary database once the workflow completes.

## Verdict
**Decision:** We will **NOT** implement Temporal at this stage.

**Reasoning:**
*   **Scale:** For the current scale (< 1,000 feeds), horizontal scaling of simple Go workers is sufficient and cost-effective.
*   **Simplicity:** Keeping the architecture as a single deployable binary is a priority.
*   **Alternative Path:** If the blocking nature of the LLM becomes a bottleneck, we will first implement an **internal decoupled pattern**:
    1. Syncer fetches items and saves them with a `NEEDS_THINKING` status.
    2. Syncer releases the lock immediately.
    3. A separate internal goroutine pool (or separate "Thinker" binary) polls for `NEEDS_THINKING` items and processes them.
    This provides the parallelism benefits without the infrastructure weight of Temporal.
