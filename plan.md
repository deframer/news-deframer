## Stop Words Plan

### Goal

Add language-wide and feed-specific noun stop words so noun-based trend metrics exclude them before aggregation.

### Database Tables

Use one table matching the model you described.

#### `stop_words`

- `id uuid primary key default uuid_generate_v4()`
- `created_at timestamp not null default now()`
- `updated_at timestamp not null default now()`
- `language char(2) not null`
- `feed_id uuid null references feeds(id) on delete cascade on update cascade`
- `nouns text[] not null default '{}'`

Semantics:

- `feed_id IS NULL`: global stop words for a language
- `feed_id = <feed uuid>`: feed-specific stop words
- effective stop words for a feed or domain = global language row(s) union feed-specific row(s)

Suggested constraints and indexes:

- unique index on `(language, feed_id)` so one row owns one stop-word list for that scope
- btree index on `feed_id`
- btree index on `language`

Notes:

- Keep `nouns` as raw user input words in admin JSON and commands.
- The mining side should stem these words into the same form as `trends.noun_stems` before they are used by SQL.
- If we want to avoid storing both raw and stemmed forms in one array, add a second column later:
  - `noun_stems text[] not null default '{}'`
- For the first pass, the cleaner plan is to store both:
  - `nouns text[]`
  - `noun_stems text[]`

Recommended final table if we want SQL to stay simple:

#### `stop_words`

- `id uuid primary key default uuid_generate_v4()`
- `created_at timestamp not null default now()`
- `updated_at timestamp not null default now()`
- `language char(2) not null`
- `feed_id uuid null references feeds(id) on delete cascade on update cascade`
- `nouns text[] not null default '{}'`
- `noun_stems text[] not null default '{}'`

Why add `noun_stems`:

- admin/import can keep human-readable input
- SQL compares against already-stemmed values
- the view does not need to stem at query time

### Go Model and Migration Changes

Files to change:

- `pkg/database/models.go`
- `pkg/database/migrate.go`
- `tools/db-joiner/db_joiner.py`

Add a GORM model such as `StopWords` and include it in `AutoMigrate(...)`.

If `noun_stems` is added to the table, update `tools/db-joiner/db_joiner.py` so source/destination joins include the new table and all of its columns.

### Admin Commands

Current admin command style is Cobra under `cmd/admin/`.

Add a new command group:

- `admin stopwords`

Recommended subcommands:

- `admin stopwords list`
- `admin stopwords import --file dummy-stopwords.json`
- `admin stopwords export --file stopwords.json`
- `admin stopwords validate --file dummy-stopwords.json`
- `admin stopwords set-global <language> <noun1,noun2,...>`
- `admin stopwords delete-global <language>`
- `admin stopwords set-feed <uuid|url> <noun1,noun2,...>`
- `admin stopwords delete-feed <uuid|url>`

Files to add or change:

- `cmd/admin/main.go`
- `cmd/admin/feed.go` only if feed import gets an inline `stop_words` field
- `cmd/admin/importexport.go` if stop words are piggybacked onto feed import
- new file: `cmd/admin/stopwords.go`

Recommendation:

- keep stop words as their own admin command group
- do not overload existing feed import/export more than necessary
- optionally allow feed import JSON to contain a feed-local `stop_words` array for convenience

### Feed Import Changes

If you want feed-specific stop words directly in feed import, extend `ImportFeed` in `cmd/admin/importexport.go` with:

- `StopWords []string \`json:"stop_words,omitempty"\``

Then, after the feed is created or updated:

- stem the provided stop words
- upsert the feed-specific `stop_words` row for that feed

This covers your requirement that feed import can carry a feed-specific list.

### Global Stop Words Command

Yes, you also need a separate global stop-word workflow.

That should write rows with:

- `language = <lang>`
- `feed_id = NULL`

This is best exposed as:

- `admin stopwords set-global`
- `admin stopwords import`

### JSON Files To Add

Add example files at the repo root.

#### `dummy-stopwords.json`

Purpose:

- demo data for manual imports
- fixture for validating the import format

Suggested shape:

```json
[
  {
    "language": "de",
    "feed_url": null,
    "nouns": ["montag", "dienstag", "mittwoch", "donnerstag", "freitag", "samstag", "sonntag", "horoskop", "wetter"]
  },
  {
    "language": "en",
    "feed_url": null,
    "nouns": ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday", "horoscope", "weather"]
  },
  {
    "language": "da",
    "feed_url": null,
    "nouns": ["mandag", "tirsdag", "onsdag", "torsdag", "fredag", "lordag", "sondag", "horoskop", "vejr"]
  },
  {
    "language": "de",
    "feed_url": "https://example.com/feed.xml",
    "nouns": ["bundesliga", "promi"]
  }
]
```

Notes:

- use `feed_url` in JSON for admin convenience, then resolve it to `feed_id`
- export can include both `feed_id` and `feed_url` if useful

Optional additional file:

- `stopwords.json` as the real managed dataset

### news-deframer-mining Changes

The miner must provide a stop-word stemming path using the exact same stemming logic used for trends.

Required capability:

- input: language + raw nouns
- output: normalized/stemmed nouns in the same form as `trends.noun_stems`

Recommended addition in `news-deframer-mining`:

- a small CLI or helper module to stem stop-word lists
- example command shape:
  - `python -m news_deframer_mining.stopwords --input dummy-stopwords.json --output stopwords.stemmed.json`

The admin import path can then either:

- call the miner tool before importing, or
- require the JSON to already contain stemmed values, or
- later expose a service/API from the miner

Recommended first pass:

- admin imports raw `nouns`
- miner helper produces `noun_stems`
- database row stores both arrays

### View Changes

Files to change:

- `pkg/database/sql/views/create_view_trend_metrics.sql`
- `pkg/database/sql/views/create_view_trend_metrics_by_domain.sql`

Conceptual change:

1. keep verb and adjective branches as they are
2. only change the noun branch
3. after `unnest(noun_stems)`, exclude stop words before aggregation

To build effective stop words for a trend row:

- include global rows where `stop_words.language = trends.language` and `stop_words.feed_id IS NULL`
- include feed-specific rows where `stop_words.feed_id` belongs to a feed with the same `root_domain` as the trend row

Because `stop_words` stores arrays, the view logic will need to expand them too.

Conceptually:

- build a CTE of effective stop-word stems by `language + root_domain`
- unnest `trends.noun_stems`
- anti-join the unnested noun stems against the effective stop-word stems
- aggregate the remaining nouns

High-level view shape:

```sql
WITH effective_stop_words AS (
  ... global language rows ...
  UNION
  ... feed-specific rows joined through feeds to root_domain ...
),
raw_unrolled AS (
  ... noun branch with unnest(noun_stems) minus effective_stop_words ...
  UNION ALL
  ... verb branch unchanged ...
  UNION ALL
  ... adjective branch unchanged ...
)
...
```

### Direct Query Follow-Up

These queries do not use the views and will still need a stop-word guard if you want consistent behavior for direct term lookups:

- `pkg/database/sql/statement/articles_by_trend.sql`
- `pkg/database/sql/statement/sentiments_by_trend.sql`
- `pkg/database/sql/statement/context_by_domain.sql`

That can be done later, but it should be part of the full rollout.

### Repository Additions

Repository methods likely needed:

- `UpsertStopWords(language string, feedID *uuid.UUID, nouns []string, nounStems []string) error`
- `DeleteStopWords(language string, feedID *uuid.UUID) error`
- `ListStopWords(...) ([]StopWords, error)`
- optional resolver helper from `uuid|url` to feed

Files to change:

- `pkg/database/repository.go`
- `pkg/database/repository_test.go`

### Suggested Rollout Order

1. add `StopWords` model and migration
2. add repository methods and tests
3. add `admin stopwords` commands
4. add `dummy-stopwords.json`
5. add miner stemming helper and decide import flow for `noun_stems`
6. update both trend views to exclude effective stop words
7. update direct noun-based statements for consistency
8. test with one global German list and one feed-specific list

### Minimal First Slice

If you want the smallest useful first delivery:

1. add `stop_words(language, feed_id, nouns, noun_stems)`
2. add `admin stopwords import`
3. add `dummy-stopwords.json`
4. update only the two trend views

That gives you cleaned noun trend metrics quickly, while direct article and sentiment lookups can follow after.
