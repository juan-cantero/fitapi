-- Create workout_sessions table
-- Records actual workout performances
CREATE TABLE IF NOT EXISTS workout_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    workout_id UUID REFERENCES workouts(id) ON DELETE SET NULL,  -- Template used (nullable for ad-hoc workouts)
    name TEXT,  -- Optional custom name for this session

    -- Timing
    started_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    duration_minutes INTEGER,

    -- Session status
    status TEXT NOT NULL DEFAULT 'planned' CHECK (status IN ('planned', 'in_progress', 'completed', 'cancelled', 'paused')),

    -- Environmental context
    location TEXT,
    weather_conditions TEXT,

    -- Subjective metrics
    energy_level_start INTEGER CHECK (energy_level_start BETWEEN 1 AND 10),
    energy_level_end INTEGER CHECK (energy_level_end BETWEEN 1 AND 10),
    perceived_exertion INTEGER CHECK (perceived_exertion BETWEEN 1 AND 10),
    mood_before TEXT,
    mood_after TEXT,

    -- Performance metrics
    calories_burned INTEGER,
    heart_rate_avg INTEGER,
    heart_rate_max INTEGER,

    -- Session feedback
    notes TEXT,
    workout_rating INTEGER CHECK (workout_rating BETWEEN 1 AND 5),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for user's workout history
CREATE INDEX idx_workout_sessions_user ON workout_sessions(user_id);

-- Index for chronological history
CREATE INDEX idx_workout_sessions_user_date ON workout_sessions(user_id, started_at DESC);

-- Index for active workouts
CREATE INDEX idx_workout_sessions_status ON workout_sessions(status) WHERE status IN ('in_progress', 'paused');

-- Auto-update updated_at timestamp
CREATE TRIGGER update_workout_sessions_updated_at
    BEFORE UPDATE ON workout_sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
