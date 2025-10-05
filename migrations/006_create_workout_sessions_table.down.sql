-- Rollback: Drop workout_sessions table
DROP TRIGGER IF EXISTS update_workout_sessions_updated_at ON workout_sessions;
DROP TABLE IF EXISTS workout_sessions CASCADE;
