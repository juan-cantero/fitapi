-- Rollback: Drop exercises table
DROP TRIGGER IF EXISTS update_exercises_updated_at ON exercises;
DROP TABLE IF EXISTS exercises CASCADE;
