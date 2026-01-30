#!/bin/bash
set -e
set -o pipefail

# Defaults based on docker-compose.yml
DATABASE=${1:-deframer}
DB_USER=${POSTGRES_USER:-deframer}
DB_PASSWORD=${POSTGRES_PASSWORD:-deframer}

INPUT_FILE="${DATABASE}.dump"
VERIFY_DB="${DATABASE}_verify_$(date +%s)"

if [ ! -f "${INPUT_FILE}" ]; then
    echo "Error: Dump file '${INPUT_FILE}' not found."
    exit 1
fi

echo "Verifying backup '${INPUT_FILE}'..."

# Ensure cleanup happens on exit (success or failure)
trap 'echo "Cleaning up..."; docker compose exec -T -e PGPASSWORD="${DB_PASSWORD}" postgres psql -U ${DB_USER} -d postgres -c "DROP DATABASE IF EXISTS ${VERIFY_DB};" > /dev/null' EXIT

# Create temporary database
echo "Creating temporary database: ${VERIFY_DB}"
docker compose exec -T -e PGPASSWORD="${DB_PASSWORD}" postgres \
  psql -U ${DB_USER} -d postgres -c "CREATE DATABASE ${VERIFY_DB} WITH OWNER ${DB_USER};"

# Restore dump to temporary database
# -v ON_ERROR_STOP=1 ensures psql exits with error code if any SQL command fails
echo "Restoring dump to ${VERIFY_DB}..."
cat "${INPUT_FILE}" | docker compose exec -T -e PGPASSWORD="${DB_PASSWORD}" postgres \
  psql -U ${DB_USER} -v ON_ERROR_STOP=1 ${VERIFY_DB} > /dev/null

echo "Verification successful!"