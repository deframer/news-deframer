#!/bin/bash

echo "Waiting for ScyllaDB to be ready..."

# Loop until cqlsh can connect to the scylla-node on port 9042
until cqlsh scylla-node 9042 -e "DESCRIBE KEYSPACES"; do
  echo "ScyllaDB is unavailable - sleeping"
  sleep 2
done

echo "ScyllaDB is up! Applying schema..."
cqlsh scylla-node 9042 -f /data/schema.cql

echo "Schema applied successfully."
