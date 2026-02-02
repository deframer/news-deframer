#!/bin/bash
set -e
set -o pipefail

if [ -z "$1" ] || [ "$1" == "-h" ] || [ "$1" == "--help" ]; then
    echo "Usage: $0 <input_file> [database_name]"
    exit 1
fi

INPUT_FILE="$1"
DATABASE=${2:-deframer}
DB_USER=${POSTGRES_USER:-deframer}
DB_PASSWORD=${POSTGRES_PASSWORD:-deframer}

if [ ! -f "${INPUT_FILE}" ]; then
    echo "Error: Dump file '${INPUT_FILE}' not found."
    exit 1
fi

echo "Starting restore for database: ${DATABASE} from ${INPUT_FILE}..."

# Drop and recreate database
# We connect to the 'postgres' maintenance database to perform operations on the target database
# We also terminate any existing connections to the target database to ensure DROP works
echo "Recreating database..."

# Terminate connections
docker compose exec -T -e PGPASSWORD="${DB_PASSWORD}" postgres \
  psql -U ${DB_USER} -d postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '${DATABASE}';" || true

docker compose exec -T -e PGPASSWORD="${DB_PASSWORD}" postgres \
  psql -U ${DB_USER} -d postgres -c "DROP DATABASE IF EXISTS ${DATABASE};" || { echo "Error: Failed to drop database."; exit 1; }

docker compose exec -T -e PGPASSWORD="${DB_PASSWORD}" postgres \
  psql -U ${DB_USER} -d postgres -c "CREATE DATABASE ${DATABASE} WITH OWNER ${DB_USER};" || { echo "Error: Failed to create database."; exit 1; }

# Restore dump
echo "Restoring data..."
cat "${INPUT_FILE}" | docker compose exec -T -e PGPASSWORD="${DB_PASSWORD}" postgres \
  psql -U ${DB_USER} ${DATABASE}

echo "Restore completed."