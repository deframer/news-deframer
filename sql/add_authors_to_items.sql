DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'items'
          AND column_name = 'authors'
    ) THEN
        ALTER TABLE items
        ADD COLUMN authors text[] NOT NULL DEFAULT '{}';
    END IF;
END $$;
