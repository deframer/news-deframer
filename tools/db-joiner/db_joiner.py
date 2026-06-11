#!/usr/bin/env python3

from __future__ import annotations

import argparse
import os
import sys
import time
from dataclasses import dataclass
from typing import Any

import psycopg
from psycopg.rows import dict_row
from psycopg.types.json import Jsonb


@dataclass
class Stats:
    feeds_inserted: int = 0
    feeds_updated: int = 0
    items_skipped: int = 0
    items_replaced: int = 0
    items_inserted: int = 0
    trends_skipped: int = 0
    trends_inserted: int = 0
    trends_replaced: int = 0


TABLES = ("feeds", "items", "trends")


FEED_UPDATE_SQL = """
UPDATE feeds
SET
  created_at = %(created_at)s,
  updated_at = %(updated_at)s,
  deleted_at = NULL,
  url = %(url)s,
  root_domain = %(root_domain)s,
  portal_url = %(portal_url)s,
  language = %(language)s,
  country = %(country)s,
  enforce_feed_domain = %(enforce_feed_domain)s,
  enabled = %(enabled)s,
  polling = %(polling)s,
  mining = %(mining)s,
  resolve_item_url = %(resolve_item_url)s,
  last_synced_at = %(last_synced_at)s,
  last_error = %(last_error)s,
  categories = %(categories)s,
  tags = %(tags)s
WHERE id = %(dest_id)s
"""


FEED_INSERT_SQL = """
INSERT INTO feeds (
  id, created_at, updated_at, deleted_at, url, root_domain, portal_url, language,
  country, enforce_feed_domain, enabled, polling, mining, resolve_item_url,
  last_synced_at, last_error, categories, tags
)
VALUES (
  %(id)s, %(created_at)s, %(updated_at)s, NULL, %(url)s, %(root_domain)s, %(portal_url)s, %(language)s,
  %(country)s, %(enforce_feed_domain)s, %(enabled)s, %(polling)s, %(mining)s, %(resolve_item_url)s,
  %(last_synced_at)s, %(last_error)s, %(categories)s, %(tags)s
)
"""


ITEM_INSERT_SQL = """
INSERT INTO items (
  id, created_at, updated_at, hash, feed_id, url, language, content, pub_date,
  media_content, think_result, think_error, think_error_count, think_rating,
  categories, authors
)
VALUES (
  %(id)s, %(created_at)s, %(updated_at)s, %(hash)s, %(feed_id)s, %(url)s,
  %(language)s, %(content)s, %(pub_date)s, %(media_content)s, %(think_result)s,
  %(think_error)s, %(think_error_count)s, %(think_rating)s, %(categories)s, %(authors)s
)
"""


TREND_INSERT_SQL = """
INSERT INTO trends (
  item_id, feed_id, language, pub_date, category_stems, noun_stems,
  verb_stems, adjective_stems, root_domain, sentiments, sentiments_deframed
)
VALUES (
  %(item_id)s, %(feed_id)s, %(language)s, %(pub_date)s, %(category_stems)s,
  %(noun_stems)s, %(verb_stems)s, %(adjective_stems)s, %(root_domain)s,
  %(sentiments)s, %(sentiments_deframed)s
)
"""


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Merge source DB into destination DB")
    parser.add_argument("--source-dsn", required=True, help="PostgreSQL DSN for source")
    parser.add_argument(
        "--dest-dsn", required=True, help="PostgreSQL DSN for destination"
    )
    parser.add_argument(
        "--apply", action="store_true", help="Apply join changes (default: dry-run)"
    )
    parser.add_argument("--join", action="store_true", help="Alias for --apply")
    parser.add_argument(
        "--validate-schema",
        action="store_true",
        help="Compare source/destination schema for feeds/items/trends and exit",
    )
    parser.add_argument("--feed-url", help="Only process one feed URL")
    parser.add_argument("--limit-feeds", type=int, help="Process first N source feeds")
    return parser.parse_args()


def fetch_table_columns(conn: psycopg.Connection):
    with conn.cursor(row_factory=dict_row) as cur:
        cur.execute(
            """
            SELECT
              table_name,
              column_name,
              is_nullable,
              data_type,
              udt_name,
              column_default
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = ANY(%s)
            ORDER BY table_name, column_name
            """,
            (list(TABLES),),
        )
        return cur.fetchall()


def fetch_indexes(conn: psycopg.Connection):
    with conn.cursor(row_factory=dict_row) as cur:
        cur.execute(
            """
            SELECT schemaname, tablename, indexname, indexdef
            FROM pg_indexes
            WHERE schemaname = 'public'
              AND tablename = ANY(%s)
            ORDER BY tablename, indexname
            """,
            (list(TABLES),),
        )
        return cur.fetchall()


def fetch_constraints(conn: psycopg.Connection):
    with conn.cursor(row_factory=dict_row) as cur:
        cur.execute(
            """
            SELECT
              tc.table_name,
              tc.constraint_name,
              tc.constraint_type,
              pg_get_constraintdef(c.oid) AS definition
            FROM information_schema.table_constraints tc
            JOIN pg_constraint c ON c.conname = tc.constraint_name
            JOIN pg_namespace n ON n.oid = c.connamespace
            WHERE tc.table_schema = 'public'
              AND n.nspname = 'public'
              AND tc.table_name = ANY(%s)
            ORDER BY tc.table_name, tc.constraint_name
            """,
            (list(TABLES),),
        )
        return cur.fetchall()


def normalize_column_rows(rows):
    return {
        (
            r["table_name"],
            r["column_name"],
            r["is_nullable"],
            r["data_type"],
            r["udt_name"],
            r["column_default"],
        )
        for r in rows
    }


def normalize_index_rows(rows):
    return {
        (
            r["tablename"],
            r["indexname"],
            r["indexdef"],
        )
        for r in rows
    }


def normalize_constraint_rows(rows):
    return {
        (
            r["table_name"],
            r["constraint_name"],
            r["constraint_type"],
            r["definition"],
        )
        for r in rows
    }


def print_diff(label: str, source_set, dest_set) -> bool:
    source_only = sorted(source_set - dest_set)
    dest_only = sorted(dest_set - source_set)
    if not source_only and not dest_only:
        print(f"{label}: equivalent")
        return True

    print(f"{label}: DIFFERENT")
    if source_only:
        print("  source-only:")
        for entry in source_only:
            print(f"    - {entry}")
    if dest_only:
        print("  destination-only:")
        for entry in dest_only:
            print(f"    - {entry}")
    return False


def validate_schema(source_dsn: str, dest_dsn: str) -> bool:
    with (
        psycopg.connect(source_dsn, autocommit=True) as src_conn,
        psycopg.connect(dest_dsn, autocommit=True) as dest_conn,
    ):
        src_cols = normalize_column_rows(fetch_table_columns(src_conn))
        dst_cols = normalize_column_rows(fetch_table_columns(dest_conn))
        src_idx = normalize_index_rows(fetch_indexes(src_conn))
        dst_idx = normalize_index_rows(fetch_indexes(dest_conn))
        src_cons = normalize_constraint_rows(fetch_constraints(src_conn))
        dst_cons = normalize_constraint_rows(fetch_constraints(dest_conn))

        ok_columns = print_diff("columns", src_cols, dst_cols)
        ok_indexes = print_diff("indexes", src_idx, dst_idx)
        ok_constraints = print_diff("constraints", src_cons, dst_cons)

    return ok_columns and ok_indexes and ok_constraints


def fetch_source_feeds(
    src_conn: psycopg.Connection, feed_url: str | None, limit: int | None
) -> list[dict[str, Any]]:
    sql = """
    SELECT id, created_at, updated_at, deleted_at, url, root_domain, portal_url, language,
           country, enforce_feed_domain, enabled, polling, mining, resolve_item_url,
           last_synced_at, last_error, categories, tags
    FROM feeds
    WHERE deleted_at IS NULL
      AND enabled IS TRUE
      AND url IS NOT NULL
      AND url <> ''
    """
    params: list[Any] = []
    if feed_url:
        sql += " AND url = %s"
        params.append(feed_url)
    sql += " ORDER BY created_at ASC"
    if limit:
        sql += " LIMIT %s"
        params.append(limit)
    with src_conn.cursor(row_factory=dict_row) as cur:
        cur.execute(sql, params)
        return cur.fetchall()


def find_dest_feed_by_url(
    dest_cur: psycopg.Cursor[Any], url: str
) -> dict[str, Any] | None:
    dest_cur.execute(
        """
        SELECT id
        FROM feeds
        WHERE url = %s AND deleted_at IS NULL
        ORDER BY updated_at DESC, created_at DESC
        LIMIT 1
        """,
        (url,),
    )
    return dest_cur.fetchone()


def source_items_for_feed(
    src_conn: psycopg.Connection, source_feed_id: Any
) -> list[dict[str, Any]]:
    with src_conn.cursor(row_factory=dict_row) as cur:
        cur.execute(
            """
            SELECT id, created_at, updated_at, hash, feed_id, url, language, content, pub_date,
                   media_content, think_result, think_error, think_error_count, think_rating,
                   categories, authors
            FROM items
            WHERE feed_id = %s
            ORDER BY pub_date ASC, created_at ASC
            """,
            (source_feed_id,),
        )
        return cur.fetchall()


def source_trends_for_feed(
    src_conn: psycopg.Connection, source_feed_id: Any
) -> list[dict[str, Any]]:
    with src_conn.cursor(row_factory=dict_row) as cur:
        cur.execute(
            """
            SELECT t.item_id AS source_item_id,
                   t.feed_id,
                   t.language,
                   t.pub_date,
                   t.category_stems,
                   t.noun_stems,
                   t.verb_stems,
                   t.adjective_stems,
                   t.root_domain,
                   t.sentiments,
                   t.sentiments_deframed
            FROM trends t
            JOIN items i ON i.id = t.item_id
            WHERE i.feed_id = %s
            """,
            (source_feed_id,),
        )
        return cur.fetchall()


def prepare_item_rows(rows: list[dict[str, Any]]) -> list[dict[str, Any]]:
    prepared: list[dict[str, Any]] = []
    for row in rows:
        prepared.append(
            {
                **row,
                "media_content": Jsonb(row["media_content"])
                if row["media_content"] is not None
                else None,
                "think_result": Jsonb(row["think_result"])
                if row["think_result"] is not None
                else None,
            }
        )
    return prepared


def prepare_trend_rows(rows: list[dict[str, Any]]) -> list[dict[str, Any]]:
    prepared: list[dict[str, Any]] = []
    for row in rows:
        prepared.append(
            {
                **row,
                "sentiments": Jsonb(row["sentiments"])
                if row["sentiments"] is not None
                else None,
                "sentiments_deframed": Jsonb(row["sentiments_deframed"])
                if row["sentiments_deframed"] is not None
                else None,
            }
        )
    return prepared


def log_progress(message: str) -> None:
    print(message, file=sys.stderr, flush=True)


def format_duration(seconds: float) -> str:
    total_seconds = max(0, int(round(seconds)))
    hours, remainder = divmod(total_seconds, 3600)
    minutes, secs = divmod(remainder, 60)
    if hours:
        return f"{hours}h{minutes:02d}m{secs:02d}s"
    if minutes:
        return f"{minutes}m{secs:02d}s"
    return f"{secs}s"


def log_apply_eta(started_at: float, completed: int, total: int) -> None:
    elapsed = time.monotonic() - started_at
    if completed < total:
        eta = elapsed / completed * (total - completed)
        log_progress(
            f"  apply progress: {completed}/{total} feeds, elapsed {format_duration(elapsed)}, eta {format_duration(eta)}"
        )
    else:
        log_progress(
            f"  apply progress: {completed}/{total} feeds, elapsed {format_duration(elapsed)}, done"
        )


def env_flag(name: str) -> bool:
    value = os.getenv(name)
    return value is not None and value.lower() not in {"", "0", "false", "no"}


def find_locked_item_urls(
    dest_cur: psycopg.Cursor[Any], dest_feed_id: Any, urls: list[str]
) -> set[str]:
    if not urls:
        return set()

    dest_cur.execute(
        """
        SELECT DISTINCT i.url
        FROM items i
        JOIN trends t ON t.item_id = i.id
        WHERE i.feed_id = %s
          AND i.think_result IS NOT NULL
          AND i.url = ANY(%s)
        """,
        (dest_feed_id, urls),
    )
    return {row["url"] for row in dest_cur.fetchall()}


def merge(args: argparse.Namespace) -> Stats:
    stats = Stats()
    force_replace = env_flag("FORCE_REPLACE")
    apply_started_at = time.monotonic() if args.apply else None
    with (
        psycopg.connect(args.source_dsn, autocommit=False) as src_conn,
        psycopg.connect(args.dest_dsn, autocommit=False) as dest_conn,
    ):
        src_feeds = fetch_source_feeds(src_conn, args.feed_url, args.limit_feeds)
        total_feeds = len(src_feeds)
        log_progress(f"merge: processing {total_feeds} feed(s)")
        if force_replace:
            log_progress("merge: FORCE_REPLACE enabled")
        for feed_index, src_feed in enumerate(src_feeds, start=1):
            # One transaction per feed so feed update + item replacements + trend inserts
            # are atomic for that feed.
            log_progress(f"[{feed_index}/{total_feeds}] {src_feed['url']}")
            with dest_conn.transaction():
                with dest_conn.cursor(row_factory=dict_row) as dest_cur:
                    dest_feed = find_dest_feed_by_url(dest_cur, src_feed["url"])

                    if dest_feed is None:
                        if args.apply:
                            dest_cur.execute(FEED_INSERT_SQL, src_feed)
                        dest_feed_id = src_feed["id"]
                        stats.feeds_inserted += 1
                    else:
                        feed_params = dict(src_feed)
                        feed_params["dest_id"] = dest_feed["id"]
                        if args.apply:
                            dest_cur.execute(FEED_UPDATE_SQL, feed_params)
                        dest_feed_id = dest_feed["id"]
                        stats.feeds_updated += 1

                    src_items = source_items_for_feed(src_conn, src_feed["id"])
                    src_trends = source_trends_for_feed(src_conn, src_feed["id"])
                    locked_urls = (
                        set()
                        if force_replace
                        else find_locked_item_urls(
                            dest_cur,
                            dest_feed_id,
                            [row["url"] for row in src_items],
                        )
                    )
                    work_items = [row for row in src_items if row["url"] not in locked_urls]
                    skipped_item_ids = {row["id"] for row in src_items if row["url"] in locked_urls}
                    work_trends = [
                        row for row in src_trends if row["source_item_id"] not in skipped_item_ids
                    ]
                    stg_items = prepare_item_rows(work_items)
                    stg_trends = prepare_trend_rows(work_trends)

                    log_progress(
                        f"  source items: {len(src_items)}, source trends: {len(src_trends)}"
                    )
                    if locked_urls:
                        log_progress(f"  skipping {len(locked_urls)} item(s) already published")
                    log_progress(
                        f"  work items: {len(work_items)}, work trends: {len(work_trends)}"
                    )

                    stats.items_skipped += len(locked_urls)
                    stats.trends_skipped += len(src_trends) - len(work_trends)

                    if not work_items:
                        log_progress("  no remaining items, skipping")
                        if apply_started_at is not None:
                            log_apply_eta(apply_started_at, feed_index, total_feeds)
                        continue

                    dest_cur.execute(
                        """
                        CREATE TEMP TABLE stg_items (
                          source_item_id uuid,
                          created_at timestamptz,
                          updated_at timestamptz,
                          hash char(64),
                          url text,
                          language char(2),
                          content text,
                          pub_date timestamptz,
                          media_content jsonb,
                          think_result jsonb,
                          think_error text,
                          think_error_count integer,
                          think_rating double precision,
                          categories text[],
                          authors text[]
                        ) ON COMMIT DROP
                        """
                    )
                    dest_cur.execute(
                        """
                        CREATE TEMP TABLE stg_trends (
                          source_item_id uuid,
                          language text,
                          pub_date timestamptz,
                          category_stems text[],
                          noun_stems text[],
                          verb_stems text[],
                          adjective_stems text[],
                          root_domain text,
                          sentiments jsonb,
                          sentiments_deframed jsonb
                        ) ON COMMIT DROP
                        """
                    )

                    dest_cur.executemany(
                        """
                        INSERT INTO stg_items (
                          source_item_id, created_at, updated_at, hash, url, language, content,
                          pub_date, media_content, think_result, think_error, think_error_count,
                          think_rating, categories, authors
                        )
                        VALUES (
                          %(id)s, %(created_at)s, %(updated_at)s, %(hash)s, %(url)s, %(language)s,
                          %(content)s, %(pub_date)s, %(media_content)s, %(think_result)s,
                          %(think_error)s, %(think_error_count)s, %(think_rating)s,
                          %(categories)s, %(authors)s
                        )
                        """,
                        stg_items,
                    )

                    if stg_trends:
                        dest_cur.executemany(
                            """
                            INSERT INTO stg_trends (
                              source_item_id, language, pub_date, category_stems, noun_stems,
                              verb_stems, adjective_stems, root_domain, sentiments, sentiments_deframed
                            )
                            VALUES (
                              %(source_item_id)s, %(language)s, %(pub_date)s, %(category_stems)s,
                              %(noun_stems)s, %(verb_stems)s, %(adjective_stems)s, %(root_domain)s,
                              %(sentiments)s, %(sentiments_deframed)s
                            )
                            """,
                            stg_trends,
                        )

                    dest_cur.execute(
                        """
                        SELECT COUNT(*)
                        FROM stg_items s
                        JOIN items i ON i.feed_id = %s AND i.url = s.url
                        """,
                        (dest_feed_id,),
                    )
                    replaced_row = dest_cur.fetchone()
                    replaced_count = (
                        int(replaced_row["count"]) if replaced_row is not None else 0
                    )
                    stats.items_replaced += replaced_count
                    stats.items_inserted += len(work_items) - replaced_count
                    log_progress(
                        f"  items: skip {len(locked_urls)}, replace {replaced_count}, insert {len(work_items) - replaced_count}"
                    )

                    if args.apply:
                        dest_cur.execute(
                            """
                            CREATE TEMP TABLE item_id_map (
                              source_item_id uuid PRIMARY KEY,
                              new_item_id uuid NOT NULL
                            ) ON COMMIT DROP
                            """
                        )

                        dest_cur.execute(
                            """
                            CREATE TEMP TABLE replaced_items AS
                            SELECT i.id
                            FROM stg_items s
                            JOIN items i ON i.feed_id = %s AND i.url = s.url
                            """,
                            (dest_feed_id,),
                        )
                        dest_cur.execute(
                            """
                            SELECT COUNT(*)
                            FROM trends t
                            JOIN replaced_items r ON r.id = t.item_id
                            """
                        )
                        replaced_trends_row = dest_cur.fetchone()
                        stats.trends_replaced += (
                            int(replaced_trends_row["count"])
                            if replaced_trends_row is not None
                            else 0
                        )

                        dest_cur.execute(
                            "DELETE FROM trends WHERE item_id IN (SELECT id FROM replaced_items)"
                        )
                        dest_cur.execute(
                            "DELETE FROM items WHERE id IN (SELECT id FROM replaced_items)"
                        )

                        dest_cur.execute(
                            """
                            INSERT INTO item_id_map(source_item_id, new_item_id)
                            SELECT s.source_item_id,
                                   CASE
                                     WHEN EXISTS (SELECT 1 FROM items i WHERE i.id = s.source_item_id)
                                       THEN uuid_generate_v4()
                                     ELSE s.source_item_id
                                   END AS new_item_id
                            FROM stg_items s
                            """
                        )

                        dest_cur.execute(
                            """
                            INSERT INTO items (
                              id, created_at, updated_at, hash, feed_id, url, language, content,
                              pub_date, media_content, think_result, think_error, think_error_count,
                              think_rating, categories, authors
                            )
                            SELECT m.new_item_id, s.created_at, s.updated_at, s.hash, %s, s.url,
                                   s.language, s.content, s.pub_date, s.media_content, s.think_result,
                                   s.think_error, s.think_error_count, s.think_rating, s.categories, s.authors
                            FROM stg_items s
                            JOIN item_id_map m ON m.source_item_id = s.source_item_id
                            """,
                            (dest_feed_id,),
                        )

                        if stg_trends:
                            dest_cur.execute(
                                """
                                INSERT INTO trends (
                                  item_id, feed_id, language, pub_date, category_stems, noun_stems,
                                  verb_stems, adjective_stems, root_domain, sentiments, sentiments_deframed
                                )
                                SELECT m.new_item_id, %s, t.language, t.pub_date, t.category_stems,
                                       t.noun_stems, t.verb_stems, t.adjective_stems, t.root_domain,
                                       t.sentiments, t.sentiments_deframed
                                FROM stg_trends t
                                JOIN item_id_map m ON m.source_item_id = t.source_item_id
                                """,
                                (dest_feed_id,),
                            )
                            stats.trends_inserted += len(work_trends)

                    if apply_started_at is not None:
                        log_apply_eta(apply_started_at, feed_index, total_feeds)

        if not args.apply:
            dest_conn.rollback()
    return stats


def main() -> int:
    args = parse_args()
    args.apply = args.apply or args.join

    if args.validate_schema:
        try:
            equivalent = validate_schema(args.source_dsn, args.dest_dsn)
        except Exception as exc:
            print(f"schema validation failed: {exc}", file=sys.stderr)
            return 1
        if equivalent:
            print("schema check complete: source and destination are equivalent")
            return 0
        print("schema check complete: source and destination differ")
        return 2

    mode = "APPLY" if args.apply else "DRY-RUN"
    print(f"db-joiner mode: {mode}")

    try:
        stats = merge(args)
    except Exception as exc:
        print(f"merge failed: {exc}", file=sys.stderr)
        return 1

    print("merge complete")
    print(f"feeds inserted: {stats.feeds_inserted}")
    print(f"feeds updated: {stats.feeds_updated}")
    print(f"items skipped: {stats.items_skipped}")
    print(f"items replaced: {stats.items_replaced}")
    print(f"items inserted: {stats.items_inserted}")
    print(f"trends skipped: {stats.trends_skipped}")
    print(f"trends replaced: {stats.trends_replaced}")
    print(f"trends inserted: {stats.trends_inserted}")

    if args.apply:
        print("post-join actions:")
        print("  ./bin/admin feed sync-all")
        print("  ./bin/admin feed mine-all")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
