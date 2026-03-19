#!/usr/bin/env bash
set -Eeuo pipefail

if [[ -f .env ]]; then
  set -a
  # shellcheck disable=SC1091
  source .env
  set +a
fi

: "${DEST_HOST:?DEST_HOST is required}"
: "${DEST_DIR:?DEST_DIR is required}"
: "${DEST_USER:?DEST_USER is required}"
: "${SSH_KEY:?SSH_KEY is required}"
: "${DEST_POSTGRES_USER:?DEST_POSTGRES_USER is required}"
: "${DEST_POSTGRES_PASSWORD:?DEST_POSTGRES_PASSWORD is required}"
: "${DEST_POSTGRES_DB:?DEST_POSTGRES_DB is required}"

EXPORT_FILE="${EXPORT_FILE:-deframer-export.dump}"
EXPORT_SCRIPT="${EXPORT_SCRIPT:-./export.sh}"
EXPORT_BASENAME="$(basename "$EXPORT_FILE")"
RESTORE_DB="${DEST_POSTGRES_DB}-restore"

if [[ ! -x "$EXPORT_SCRIPT" ]]; then
  echo "Export script is not executable: $EXPORT_SCRIPT" >&2
  exit 1
fi

"$EXPORT_SCRIPT"

if [[ ! -f "$EXPORT_FILE" ]]; then
  echo "Export file not found: $EXPORT_FILE" >&2
  exit 1
fi

if [[ "$(dd if="$EXPORT_FILE" bs=5 count=1 2>/dev/null)" == "PGDMP" ]]; then
  DUMP_FORMAT="custom"
else
  DUMP_FORMAT="plain"
fi

ssh -i "$SSH_KEY" \
  -o BatchMode=yes \
  -o StrictHostKeyChecking=accept-new \
  "${DEST_USER}@${DEST_HOST}" "mkdir -p '$DEST_DIR'"

rsync -avz --progress --partial --inplace \
  -e "ssh -i $SSH_KEY -o BatchMode=yes -o StrictHostKeyChecking=accept-new" \
  "$EXPORT_FILE" \
  "${DEST_USER}@${DEST_HOST}:${DEST_DIR}/"

echo "Restoring database as restore_db"
echo "starting remote restore"

ssh -i "$SSH_KEY" \
  -o BatchMode=yes \
  -o StrictHostKeyChecking=accept-new \
  "${DEST_USER}@${DEST_HOST}" \
  bash -s -- "$DEST_DIR" "$EXPORT_BASENAME" "$RESTORE_DB" "$DEST_POSTGRES_USER" "$DEST_POSTGRES_PASSWORD" "$DUMP_FORMAT" <<'REMOTE'
set -Eeuo pipefail

echo "starting restore on remote host"

DEST_DIR="$1"
EXPORT_BASENAME="$2"
RESTORE_DB="$3"
DEST_POSTGRES_USER="$4"
DEST_POSTGRES_PASSWORD="$5"
DUMP_FORMAT="$6"

cd "$DEST_DIR"

if [[ ! -f docker-compose.yml ]]; then
  echo "docker-compose.yml not found in $DEST_DIR" >&2
  exit 1
fi

if [[ ! -f "$EXPORT_BASENAME" ]]; then
  echo "Remote export file not found: $DEST_DIR/$EXPORT_BASENAME" >&2
  exit 1
fi

docker compose -f docker-compose.yml exec -T \
  -e PGPASSWORD="$DEST_POSTGRES_PASSWORD" \
  postgres \
  psql -U "$DEST_POSTGRES_USER" -d postgres -v ON_ERROR_STOP=1 \
  -v restore_db="$RESTORE_DB" \
  -v db_user="$DEST_POSTGRES_USER" <<'SQL'
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE datname = :'restore_db'
  AND pid <> pg_backend_pid();

DROP DATABASE IF EXISTS :"restore_db";
CREATE DATABASE :"restore_db" WITH OWNER :"db_user";
SQL

if [[ "$DUMP_FORMAT" == "custom" ]]; then
  echo "about to run restore"
  cat "$EXPORT_BASENAME" | docker compose -f docker-compose.yml exec -T \
    -e PGPASSWORD="$DEST_POSTGRES_PASSWORD" \
    postgres \
    pg_restore -U "$DEST_POSTGRES_USER" -d "$RESTORE_DB" --no-owner --no-privileges --verbose
else
  echo "about to run restore"
  cat "$EXPORT_BASENAME" | docker compose -f docker-compose.yml exec -T \
    -e PGPASSWORD="$DEST_POSTGRES_PASSWORD" \
    postgres \
    psql -U "$DEST_POSTGRES_USER" -d "$RESTORE_DB" -v ON_ERROR_STOP=1
fi

echo "restore command finished"

echo "about to verify items count"

if ! COUNT_RAW="$({
  docker compose -f docker-compose.yml exec -T \
    -e PGPASSWORD="$DEST_POSTGRES_PASSWORD" \
    postgres \
    psql -U "$DEST_POSTGRES_USER" -d "$RESTORE_DB" -Atqc "SELECT count(*) FROM public.items;" \
    </dev/null
} 2>&1)"; then
  echo "items count query failed:" >&2
  printf '%s\n' "$COUNT_RAW" >&2
  exit 1
fi

COUNT="$(printf '%s' "$COUNT_RAW" | tr -d '[:space:]')"

echo "count query finished"

if [[ ! "$COUNT" =~ ^[0-9]+$ ]]; then
  echo "Verification failed: invalid items count: $COUNT_RAW" >&2
  exit 1
fi

echo "items count: $COUNT"

if (( COUNT <= 0 )); then
  echo "Verification failed: items count must be > 0" >&2
  exit 1
fi

echo "Restore completed into database: $RESTORE_DB"
REMOTE

echo "starting database promotion"

ssh -i "$SSH_KEY" \
  -o BatchMode=yes \
  -o StrictHostKeyChecking=accept-new \
  "${DEST_USER}@${DEST_HOST}" \
  bash -s -- "$DEST_DIR" "$RESTORE_DB" "$DEST_POSTGRES_DB" "$DEST_POSTGRES_USER" "$DEST_POSTGRES_PASSWORD" <<'PROMOTE'
set -Eeuo pipefail

DEST_DIR="$1"
RESTORE_DB="$2"
DEST_POSTGRES_DB="$3"
DEST_POSTGRES_USER="$4"
DEST_POSTGRES_PASSWORD="$5"

cd "$DEST_DIR"

if [[ ! -f docker-compose.yml ]]; then
  echo "docker-compose.yml not found in $DEST_DIR" >&2
  exit 1
fi

if [[ "$RESTORE_DB" == "$DEST_POSTGRES_DB" ]]; then
  echo "Promotion skipped: restore database already has target name $DEST_POSTGRES_DB" >&2
  exit 1
fi

docker compose -f docker-compose.yml exec -T \
  -e PGPASSWORD="$DEST_POSTGRES_PASSWORD" \
  postgres \
  psql -U "$DEST_POSTGRES_USER" -d postgres -v ON_ERROR_STOP=1 \
  -v restore_db="$RESTORE_DB" \
  -v dest_db="$DEST_POSTGRES_DB" <<'SQL'
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE datname = :'dest_db'
  AND pid <> pg_backend_pid();

DROP DATABASE IF EXISTS :"dest_db";
ALTER DATABASE :"restore_db" RENAME TO :"dest_db";
SQL

echo "Promotion completed: $RESTORE_DB -> $DEST_POSTGRES_DB"
PROMOTE

rm -f "$EXPORT_FILE"
echo "Removed local export file: $EXPORT_FILE"
