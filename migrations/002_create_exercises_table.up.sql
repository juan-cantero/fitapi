-- Create exercises table
-- Stores exercise definitions with public/private visibility
CREATE TABLE IF NOT EXISTS exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    image_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for user's exercises
CREATE INDEX idx_exercises_user_id ON exercises(user_id);

-- Index for public exercises
CREATE INDEX idx_exercises_is_public ON exercises(is_public);

-- Composite index for "show me public exercises AND my private ones"
CREATE INDEX idx_exercises_public_user ON exercises(is_public, user_id);

-- Auto-update updated_at timestamp
CREATE TRIGGER update_exercises_updated_at
    BEFORE UPDATE ON exercises
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
