-- Rollback: Drop exercise_logs table
DROP TRIGGER IF EXISTS update_exercise_logs_updated_at ON exercise_logs;
DROP TABLE IF EXISTS exercise_logs CASCADE;
