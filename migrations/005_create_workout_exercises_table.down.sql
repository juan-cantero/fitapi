-- Rollback: Drop workout_exercises table
DROP TRIGGER IF EXISTS update_workout_exercises_updated_at ON workout_exercises;
DROP TABLE IF EXISTS workout_exercises CASCADE;
