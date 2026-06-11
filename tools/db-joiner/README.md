# db-joiner

Merge source -> destination for `feeds`, `items`, and `trends`.

## Behavior

- Process only source feeds where `deleted_at IS NULL`, `enabled IS TRUE`, and `url` is not empty.
- Match feeds by `feeds.url`; keep destination `feeds.id` if URL exists.
- If a matching destination feed does not exist, insert the source feed.
- Match items by `(feed_id, url)` in destination.
- By default, skip replacing a destination item if it already has a `think_result` and at least one `trend`.
- If an existing item has no `think_result` and no `trend`, delete it and insert the source item as new.
- Set `FORCE_REPLACE=1` to force the old delete-and-reinsert behavior.
- If item is eligible for replacement: delete destination trend + item, then insert source item fresh.
- During `make apply`, print a light feed-based ETA/progress line to stderr.
- Insert trends after items using final destination `feed_id` and `item_id`.
- Use one destination transaction per feed (atomic per-feed merge).
- Item/trend work is set-based per feed using temporary staging tables.

## Quickstart

```bash
uv venv
uv sync
```

Example DSNs:

```bash
SOURCE_DSN="postgresql://postgres:sourcepass@source-db.example.com:5432/news_deframer"
DEST_DSN="postgresql://postgres:destpass@dest-db.example.com:5432/news_deframer"
```

## Commands

```bash
# schema check only
make validate SOURCE_DSN="..." DEST_DSN="..."

# preview merge (no writes)
make dry-run SOURCE_DSN="..." DEST_DSN="..."

# execute merge writes
make apply SOURCE_DSN="..." DEST_DSN="..."

# execute join for one feed URL only
make join-feed SOURCE_DSN="..." DEST_DSN="..." FEED_URL="https://example.com/rss"

# force replace even when item already has think_result + trend
FORCE_REPLACE=1 make dry-run SOURCE_DSN="..." DEST_DSN="..."
```

Direct script form for one feed:

```bash
uv run db_joiner.py --source-dsn "..." --dest-dsn "..." --join --feed-url "https://example.com/rss"

# force replace even when item already has think_result + trend
FORCE_REPLACE=1 uv run db_joiner.py --source-dsn "..." --dest-dsn "..."
```

## After Join

After a successful join/apply run, trigger scheduler work:

```bash
./bin/admin feed sync-all
./bin/admin feed mine-all
```
