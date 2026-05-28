#!/usr/bin/env python3

from __future__ import annotations

import argparse
import json
import sys
from datetime import datetime, timezone
from typing import Any

import psycopg
from psycopg.rows import dict_row


NUMERIC_RE = r"^-?[0-9]+(\.[0-9]+)?$"


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Fetch recent analyzed items from enabled feeds"
    )
    parser.add_argument("--dsn", required=True, help="PostgreSQL DSN")
    parser.add_argument(
        "--limit-per-feed",
        type=int,
        default=20,
        help="Maximum recent items per feed (default: 20)",
    )
    parser.add_argument(
        "--pretty", action="store_true", help="Pretty-print JSON output"
    )
    return parser.parse_args()


def to_rfc3339_utc(value: datetime) -> str:
    if value.tzinfo is None:
        value = value.replace(tzinfo=timezone.utc)
    else:
        value = value.astimezone(timezone.utc)
    return value.isoformat().replace("+00:00", "Z")


def fetch_items(conn: psycopg.Connection[Any], limit_per_feed: int) -> list[dict[str, Any]]:
    sql = f"""
    WITH ranked AS (
      SELECT
        i.id::text AS id,
        i.url,
        f.root_domain,
        i.pub_date,
        i.media_content,
        i.authors,
        i.think_result,
        ROW_NUMBER() OVER (
          PARTITION BY i.feed_id
          ORDER BY i.pub_date DESC, i.created_at DESC
        ) AS rn
      FROM items i
      JOIN feeds f ON f.id = i.feed_id
      WHERE f.enabled = TRUE
        AND f.deleted_at IS NULL
        AND COALESCE(f.root_domain, '') <> ''
        AND i.think_result IS NOT NULL
        AND COALESCE(i.think_result->>'category', '') <> ''
        AND i.think_error IS NULL
        AND i.think_error_count = 0
        AND i.think_result ? 'overall'
        AND (i.think_result->>'overall') ~ '{NUMERIC_RE}'
    )
    SELECT
      id,
      url,
      root_domain,
      think_result->>'title_corrected' AS title_corrected,
      think_result->>'description_corrected' AS description_corrected,
      (think_result->>'overall')::double precision AS overall,
      think_result->>'overall_reason' AS overall_reason,
      think_result->>'category' AS category,
      media_content,
      authors,
      pub_date
    FROM ranked
    WHERE rn <= %s
    ORDER BY pub_date DESC, id ASC
    """

    with conn.cursor(row_factory=dict_row) as cur:
        cur.execute(sql, (limit_per_feed,))
        rows = cur.fetchall()

    result: list[dict[str, Any]] = []
    for row in rows:
        pub_date = row["pub_date"]
        if not isinstance(pub_date, datetime):
            raise ValueError("invalid pub_date type from database")

        result.append(
            {
                "id": row["id"],
                "url": row["url"],
                "root_domain": row["root_domain"],
                "title_corrected": row["title_corrected"],
                "description_corrected": row["description_corrected"],
                "overall": row["overall"],
                "overall_reason": row["overall_reason"],
                "category": row["category"],
                "media": row["media_content"],
                "authors": row["authors"] or [],
                "pubDate": to_rfc3339_utc(pub_date),
            }
        )

    return result


def main() -> int:
    args = parse_args()

    if args.limit_per_feed <= 0:
        print("--limit-per-feed must be > 0", file=sys.stderr)
        return 2

    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")

    try:
        with psycopg.connect(args.dsn, autocommit=True) as conn:
            items = fetch_items(conn, args.limit_per_feed)
    except Exception as exc:
        print(f"query failed: {exc}", file=sys.stderr)
        return 1

    json_kwargs: dict[str, Any] = {"ensure_ascii": False}
    if args.pretty:
        json_kwargs["indent"] = 2

    print(json.dumps(items, **json_kwargs))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
