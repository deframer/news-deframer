# Repository Refactoring Plan

## 1. Goal
Reorganize the project to support a **Polyglot Monorepo** structure, accommodating the new Python Trend Mining service while maintaining the existing Go backend.

## 2. Decision: Monorepo vs. Multirepo
**Verdict:** **Polyglot Monorepo**

**Reasoning:**
*   **Tight Coupling:** The Go Service (Syncer) and Python Service (Miner) share the exact same database schema (`items`, `feed_schedules`).
*   **Schema Evolution:** Changes to the database (e.g., adding `mining_done_at`) must be reflected in both services simultaneously.
*   **Operational Simplicity:** A single `docker-compose up` orchestrates the entire platform.
*   **AI Context:** Facilitates comprehensive code analysis and refactoring across language boundaries.

## 3. Alternative: Polyglot Multirepo

While the Monorepo (Option A) is recommended, the Multirepo (Option B) is a valid strategy if:
*   You want zero disruption to the existing Go codebase.
*   You prefer standard Go tooling (module at root).
*   Teams are completely independent.

### Structure
1.  **Repo 1: `news-deframer` (The Core)**
    *   Contains: Go Backend (`pkg`, `cmd`), `docker-compose.yml`, `docs/`.
    *   **Pros:** No file moves, no import breakage, familiar structure.
    *   **Cons:** Python service is invisible here.

2.  **Repo 2: `news-deframer-trends` (The Miner)**
    *   Contains: Python Miner (`src/`, `pyproject.toml`), Dockerfile.
    *   **Pros:** Clean Python environment, independent CI/CD.
    *   **Cons:** Must duplicate DB schema knowledge (or submodule/share protos).

### Integration
To run them together locally, you would:
1.  Run `docker-compose up` in Repo 1 (Core).
2.  Run `docker run --network news-deframer_default ...` in Repo 2 (Miner) to connect to the shared database.

## 4. Directory Structure (Monorepo Plan)

**Current:**
```
/
├── cmd/               (Go binaries)
├── pkg/               (Go libraries)
├── chrome-extension/  (Legacy/Separate)
├── infra-env/         (Docker/Env)
├── docs/
├── go.mod
└── docker-compose.yml
```

**Target:**
```
/
├── services/
│   ├── backend/       (Go Service - Moved from root)
│   │   ├── cmd/
│   │   ├── pkg/
│   │   ├── go.mod
│   │   └── Dockerfile
│   └── miner/         (New Python Service)
│       ├── src/
│       ├── pyproject.toml (uv)
│       └── Dockerfile
├── infra/             (Shared Infrastructure)
│   ├── docker-compose.yml
│   └── env/
├── docs/              (Documentation)
└── Makefile           (Orchestration)
```

## 4. Migration Steps

### Phase 1: Preparation & Cleanup
1.  **Chrome Extension:** Move `chrome-extension/` to a separate repository (or delete if unused).
2.  **Infrastructure:** Move `infra-env/` to `infra/`.

### Phase 2: Go Backend Relocation
1.  Create `services/backend/`.
2.  Move `cmd/`, `pkg/`, `go.mod`, `go.sum` into `services/backend/`.
3.  **Critical:** Update Go import paths.
    *   The module name in `go.mod` is `github.com/deframer/news-deframer`.
    *   Moving files to a subdirectory *does not* require changing import paths if the `go.mod` remains at the root of the *module* (which is now `services/backend/`).
    *   *However*, external references (Dockerfiles, CI) will need path updates.

### Phase 3: Python Miner Initialization
1.  Create `services/miner/`.
2.  Initialize with `uv init`.
3.  Add dependencies (`duckdb`, `psycopg`, `scikit-learn`, `pandas`).

### Phase 4: Orchestration Update
1.  Update root `docker-compose.yml` to point to the new build contexts:
    *   `backend`: `./services/backend`
    *   `miner`: `./services/miner`
2.  Create a root `Makefile` to simplify commands (e.g., `make up`, `make build-backend`).

## 5. Risks & Mitigation
*   **Import Breakage:** Moving `go.mod` might break local tooling (gopls).
    *   *Mitigation:* Open the editor at the root, but ensure it detects the nested module (Go Workspaces `go.work` can help here).
*   **Docker Contexts:** The `COPY` commands in Dockerfiles will need adjustment relative to the new build context.

