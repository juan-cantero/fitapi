-- Rollback: Drop workouts table
DROP TRIGGER IF EXISTS update_workouts_updated_at ON workouts;
DROP TABLE IF EXISTS workouts CASCADE;
