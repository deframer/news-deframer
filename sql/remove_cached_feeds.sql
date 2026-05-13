DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_schema = 'public'
      AND table_name = 'cached_feeds'
  ) THEN
    DROP TABLE cached_feeds CASCADE;
  END IF;
END $$;
