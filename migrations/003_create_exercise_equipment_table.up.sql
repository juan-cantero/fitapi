-- Create exercise_equipment junction table
-- Links exercises to equipment (many-to-many relationship)
CREATE TABLE IF NOT EXISTS exercise_equipment (
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    equipment_id UUID NOT NULL REFERENCES equipment(id) ON DELETE CASCADE,
    PRIMARY KEY (exercise_id, equipment_id)
);

-- Index for "what equipment does this exercise use?"
CREATE INDEX idx_exercise_equipment_exercise ON exercise_equipment(exercise_id);

-- Index for "what exercises use this equipment?"
CREATE INDEX idx_exercise_equipment_equipment ON exercise_equipment(equipment_id);
