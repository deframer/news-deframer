#!/bin/bash
set -e
set -o pipefail

# Defaults based on docker-compose.yml
DATABASE=${1:-deframer}
DB_USER=${POSTGRES_USER:-deframer}
DB_PASSWORD=${POSTGRES_PASSWORD:-deframer}

OUTPUT_FILE="${DATABASE}-$(date +%Y-%m-%d-%H-%M).dump"

echo "Starting backup for database: ${DATABASE}..."

# Use docker compose exec to stream dump to host file
# -T is required to disable TTY so binary output isn't corrupted
if docker compose exec -T -e PGPASSWORD="${DB_PASSWORD}" postgres \
  pg_dump -U ${DB_USER} ${DATABASE} > "${OUTPUT_FILE}"; then
    echo "Backup completed: ${OUTPUT_FILE}"
else
    echo "Backup failed!"
    rm -f "${OUTPUT_FILE}"
    exit 1
fi