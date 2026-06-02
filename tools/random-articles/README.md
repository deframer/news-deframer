# random-articles

Fetch analyzed items from enabled, non-deleted feeds and print JSON.

## Behavior

- Uses one DSN (`--dsn`) to read from PostgreSQL.
- Reads only feeds where `enabled = true` and `deleted_at IS NULL`.
- Returns up to 20 recent analyzed items per feed (override with `--limit-per-feed`).
- Includes `root_domain` in each returned item.
- Excludes `associated_items` entirely.

## Output fields

Each item includes:

- `id`
- `url`
- `root_domain`
- `thinker_result`
- `sentiments`
- `sentiments_deframed`
- `media`
- `authors`
- `pubDate`

## Quickstart

```bash
uv venv
uv sync
```

Example DSN:

```bash
DSN="postgresql://postgres:password@localhost:5432/news_deframer"
```

## Commands

```bash
# compact JSON output
make run DSN="${DSN}"

# pretty JSON output
make pretty DSN="${DSN}"
```

Direct script usage:

```bash
uv run random_articles.py --dsn "${DSN}" --limit-per-feed 20 --pretty
```
