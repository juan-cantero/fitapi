-- Rollback: Drop equipment table and related objects
DROP TRIGGER IF EXISTS update_equipment_updated_at ON equipment;
DROP TABLE IF EXISTS equipment CASCADE;
-- Note: We don't drop the function as it will be reused by other tables
