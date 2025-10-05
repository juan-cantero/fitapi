-- Create exercise_logs table
-- Records individual exercise performances within workout sessions
CREATE TABLE IF NOT EXISTS exercise_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_session_id UUID NOT NULL REFERENCES workout_sessions(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    workout_exercise_id UUID REFERENCES workout_exercises(id) ON DELETE SET NULL,  -- Template reference (nullable)
    order_index INTEGER NOT NULL,

    -- Planned vs Actual
    sets_completed INTEGER DEFAULT 0,
    sets_planned INTEGER DEFAULT 1,
    reps_completed INTEGER,
    reps_planned INTEGER,

    -- Performance metrics
    weight_kg REAL,
    duration_seconds INTEGER,
    distance_meters REAL,
    rest_time_seconds INTEGER,
    intensity_percentage REAL,

    -- Subjective feedback
    rpe INTEGER CHECK (rpe BETWEEN 1 AND 10),  -- Actual Rate of Perceived Exertion
    form_rating INTEGER CHECK (form_rating BETWEEN 1 AND 5),  -- How good was your form?

    -- Context
    equipment_used TEXT,  -- JSON array of equipment IDs actually used
    notes TEXT,

    -- Personal Records
    is_personal_record BOOLEAN DEFAULT FALSE,
    previous_best_weight REAL,
    previous_best_reps INTEGER,
    previous_best_duration INTEGER,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for session's exercises
CREATE INDEX idx_exercise_logs_session ON exercise_logs(workout_session_id);

-- Index for exercise history
CREATE INDEX idx_exercise_logs_exercise ON exercise_logs(exercise_id);

-- Index for finding personal records
CREATE INDEX idx_exercise_logs_pr ON exercise_logs(exercise_id, is_personal_record) WHERE is_personal_record = TRUE;

-- Composite index for progression tracking
CREATE INDEX idx_exercise_logs_progression ON exercise_logs(exercise_id, created_at DESC);

-- Auto-update updated_at timestamp
CREATE TRIGGER update_exercise_logs_updated_at
    BEFORE UPDATE ON exercise_logs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
