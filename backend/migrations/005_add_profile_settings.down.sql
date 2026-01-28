-- Rollback: Remove profile_settings table

DROP TRIGGER IF EXISTS trigger_profile_settings_updated_at ON profile_settings;
DROP FUNCTION IF EXISTS update_profile_settings_updated_at();
DROP INDEX IF EXISTS idx_profile_settings_user_id;
DROP TABLE IF EXISTS profile_settings;
