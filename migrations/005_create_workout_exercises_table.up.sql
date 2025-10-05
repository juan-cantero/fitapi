-- Create workout_exercises junction table
-- Links workouts to exercises with detailed parameters (sets, reps, intensity, etc.)
CREATE TABLE IF NOT EXISTS workout_exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_id UUID NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    order_index INTEGER NOT NULL,

    -- Basic parameters
    sets INTEGER DEFAULT 1,
    reps INTEGER,
    weight_kg REAL,
    duration_seconds INTEGER,
    distance_meters REAL,
    rest_time_seconds INTEGER DEFAULT 60,

    -- Advanced parameters
    intensity_percentage REAL,  -- % of 1RM (one-rep max)
    tempo TEXT,                  -- e.g., "3-1-2-0" (eccentric-pause-concentric-pause)
    notes TEXT,

    -- Superset/dropset tracking
    is_superset BOOLEAN DEFAULT FALSE,
    superset_group_id UUID,      -- Groups exercises performed back-to-back
    is_dropset BOOLEAN DEFAULT FALSE,

    -- Workout phase flags
    is_warmup BOOLEAN DEFAULT FALSE,
    is_cooldown BOOLEAN DEFAULT FALSE,

    -- Target intensity
    target_rpe INTEGER CHECK (target_rpe BETWEEN 1 AND 10),  -- Rate of Perceived Exertion

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for getting all exercises in a workout
CREATE INDEX idx_workout_exercises_workout ON workout_exercises(workout_id);

-- Index for ordered exercise list
CREATE INDEX idx_workout_exercises_order ON workout_exercises(workout_id, order_index);

-- Index for finding supersets
CREATE INDEX idx_workout_exercises_superset ON workout_exercises(superset_group_id) WHERE superset_group_id IS NOT NULL;

-- Auto-update updated_at timestamp
CREATE TRIGGER update_workout_exercises_updated_at
    BEFORE UPDATE ON workout_exercises
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
