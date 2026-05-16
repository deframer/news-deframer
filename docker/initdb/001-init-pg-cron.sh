#!/usr/bin/env bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<'SQL'
SELECT 'CREATE DATABASE cron'
WHERE NOT EXISTS (
  SELECT FROM pg_database WHERE datname = 'cron'
)\gexec
SQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "cron" <<'SQL'
CREATE EXTENSION IF NOT EXISTS pg_cron;
SQL
