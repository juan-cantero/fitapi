-- Create equipment table
-- Stores gym equipment definitions (dumbbells, barbells, etc.)
CREATE TABLE IF NOT EXISTS equipment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create index for fast user lookup
CREATE INDEX idx_equipment_user_id ON equipment(user_id);

-- Create trigger to auto-update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_equipment_updated_at
    BEFORE UPDATE ON equipment
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
