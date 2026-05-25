# Database Merge Plan (Source -> Destination)

## Scope

Only these tables are in scope:

- `feeds`
- `items`
- `trends`

`feed_schedules` is explicitly out of scope.

## Terms

- **Source DB**: the database copied from
- **Destination DB**: the database merged into

## Feed Matching Rule

- Use `feeds.url` as the alternative primary key for identity/matching.
- Only consider non-deleted feeds (`deleted_at IS NULL`) during matching/migration.
- If source feed URL is not present in destination: insert as new feed.
- If source feed URL already exists in destination: keep destination `feeds.id` and treat that as canonical.

## Merge Direction

- Default policy: source overwrites destination.
- For an existing feed match by URL, destination row is updated from source data, while keeping destination `id`.

## Dependent Table Mapping

After feed URL matching is resolved:

- Merge `items` using the mapped destination `feed_id`.
  - Item identity key: `(feed_id, url)`.
  - Only migrate non-deleted source feeds and their items.
  - If `(feed_id, url)` exists in destination: delete destination item, delete related trend row(s), then insert source item as fresh data.
  - If `(feed_id, url)` does not exist: insert new item.
- Merge `trends` using mapped destination `feed_id` and mapped destination `item_id`.
  - Trends are inserted for the final destination item IDs after item merge/reinsert.

## High-Level Order

1. Resolve feed URL matches and build source->destination feed ID mapping.
2. Upsert `feeds` for non-deleted source feeds (preserve destination IDs on URL matches).
3. For each source item mapped by `(feed_id, url)`, delete conflicting destination item+trend and insert source item.
4. Insert `trends` using mapped destination feed and final destination item IDs.

## Tooling Setup

- Create standalone merge tool under `tools/db-joiner/`.
- Use `uv` for environment and dependency management.
- Initialize virtual environment with `uv venv` in `tools/db-joiner/`.
- Track dependencies via `pyproject.toml` (or `requirements.txt` if preferred).
- Add `.gitignore` entries for Python/venv artifacts:
  - `.venv/`
  - `__pycache__/`
  - `*.pyc`
  - local logs/output files
