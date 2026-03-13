DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'trends'
          AND column_name = 'sentiments'
    ) THEN
        ALTER TABLE trends
        ADD COLUMN sentiments jsonb NOT NULL DEFAULT '{}';
    END IF;
END $$;
