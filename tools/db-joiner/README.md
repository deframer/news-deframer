# db-joiner

Merge source -> destination for `feeds`, `items`, and `trends`.

## Behavior

- Process only source feeds where `deleted_at IS NULL`.
- Match feeds by `feeds.url`; keep destination `feeds.id` if URL exists.
- Match items by `(feed_id, url)` in destination.
- If item exists: delete destination trend + item, then insert source item fresh.
- Insert trends after items using final destination `feed_id` and `item_id`.
- Use one destination transaction per feed (atomic per-feed merge).

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
```

Direct script form for one feed:

```bash
uv run db_joiner.py --source-dsn "..." --dest-dsn "..." --join --feed-url "https://example.com/rss"
```
