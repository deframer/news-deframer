#!/bin/bash
set -e
set -o pipefail

# Defaults based on docker-compose.yml
DATABASE=${1:-deframer}
DB_USER=${POSTGRES_USER:-deframer}
DB_PASSWORD=${POSTGRES_PASSWORD:-deframer}

# Infrastructure settings
NETWORK="deframer-infra-env-net"
SERVER="postgres"
IMAGE="postgres:18"

OUTPUT_FILE="${DATABASE}.dump"

echo "Starting backup for database: ${DATABASE}..."

# Remove old dump if it exists
rm -f "${OUTPUT_FILE}"

# Run pg_dump in a transient container attached to the infra network
if docker run --rm --network ${NETWORK} \
  -e PGPASSWORD="${DB_PASSWORD}" \
  ${IMAGE} \
  pg_dump -h ${SERVER} -U ${DB_USER} ${DATABASE} > "${OUTPUT_FILE}"; then
    echo "Backup completed: ${OUTPUT_FILE}"
else
    echo "Backup failed!"
    rm -f "${OUTPUT_FILE}"
    exit 1
fi