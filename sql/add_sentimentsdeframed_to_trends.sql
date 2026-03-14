DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'trends'
          AND column_name = 'sentiments_deframed'
    ) THEN
        ALTER TABLE trends
        ADD COLUMN sentiments_deframed jsonb NOT NULL DEFAULT '{}';
    END IF;
END $$;
