-- Create a keyspace (database) with replication factor 1 for local dev
CREATE KEYSPACE IF NOT EXISTS key_value_store
WITH REPLICATION = {
    'class' : 'SimpleStrategy',
    'replication_factor' : 1
};

USE key_value_store;

-- Create the table
-- compaction = 'TimeWindowCompactionStrategy' is BEST for time-series/TTL data
CREATE TABLE IF NOT EXISTS items (
    key text PRIMARY KEY,
    value text,
    created_at timestamp
) WITH default_time_to_live = 2592000 -- 30 days
  AND compaction = {'class': 'TimeWindowCompactionStrategy',
                    'compaction_window_unit': 'DAYS',
                    'compaction_window_size': 1};
