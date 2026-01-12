-- Run this migration to add tagged_blocs support to your database
-- Execute this SQL against your PostgreSQL database

ALTER TABLE events
ADD COLUMN IF NOT EXISTS tagged_blocs TEXT;

-- Verify the column was added
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'events' AND column_name = 'tagged_blocs';
