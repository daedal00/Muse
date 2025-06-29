-- Remove Spotify OAuth columns from users table
DROP INDEX IF EXISTS idx_users_spotify_id;
ALTER TABLE users DROP COLUMN IF EXISTS spotify_token_expiry;
ALTER TABLE users DROP COLUMN IF EXISTS spotify_refresh_token;
ALTER TABLE users DROP COLUMN IF EXISTS spotify_access_token;
ALTER TABLE users DROP COLUMN IF EXISTS spotify_id; 