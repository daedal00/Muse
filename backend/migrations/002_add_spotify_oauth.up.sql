-- Add Spotify OAuth columns to users table
ALTER TABLE users ADD COLUMN spotify_id VARCHAR(255) UNIQUE;
ALTER TABLE users ADD COLUMN spotify_access_token TEXT;
ALTER TABLE users ADD COLUMN spotify_refresh_token TEXT;
ALTER TABLE users ADD COLUMN spotify_token_expiry TIMESTAMP WITH TIME ZONE;

-- Create index for Spotify ID lookups
CREATE INDEX idx_users_spotify_id ON users(spotify_id); 