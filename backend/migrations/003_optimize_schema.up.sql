-- Migration to optimize database schema for Spotify-based app
-- Remove unnecessary tables that duplicate Spotify data
-- Keep only user-generated data

-- Drop unnecessary tables and their dependencies
DROP TABLE IF EXISTS track_featured_artists CASCADE;
DROP TABLE IF EXISTS playlist_tracks CASCADE;
DROP TABLE IF EXISTS tracks CASCADE;
DROP TABLE IF EXISTS albums CASCADE;
DROP TABLE IF EXISTS artists CASCADE;

-- Update reviews table to work with Spotify IDs instead of local references
ALTER TABLE reviews 
DROP COLUMN IF EXISTS album_id,
DROP COLUMN IF EXISTS track_id,
ADD COLUMN spotify_id VARCHAR(255) NOT NULL,
ADD COLUMN spotify_type VARCHAR(20) NOT NULL CHECK (spotify_type IN ('album', 'track'));

-- Create optimized playlist_tracks table with Spotify IDs
CREATE TABLE playlist_tracks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    playlist_id UUID NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    spotify_id VARCHAR(255) NOT NULL, -- Spotify track ID
    position INTEGER NOT NULL,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    added_by_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(playlist_id, spotify_id)
);

-- Add public/private flag to playlists
ALTER TABLE playlists 
ADD COLUMN IF NOT EXISTS is_public BOOLEAN DEFAULT true;

-- Create user preferences table for storing user settings and preferences
CREATE TABLE user_preferences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    preferred_genres JSONB DEFAULT '[]'::jsonb, -- Array of genre names
    favorite_artist_ids JSONB DEFAULT '[]'::jsonb, -- Array of Spotify artist IDs
    notification_settings JSONB DEFAULT '{}'::jsonb, -- Notification preferences
    privacy_settings JSONB DEFAULT '{}'::jsonb, -- Privacy settings
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id)
);

-- Create indexes for better performance
CREATE INDEX idx_reviews_spotify_id ON reviews(spotify_id);
CREATE INDEX idx_reviews_spotify_type ON reviews(spotify_type);
CREATE INDEX idx_reviews_user_id_spotify ON reviews(user_id, spotify_id);
CREATE INDEX idx_playlist_tracks_playlist_id ON playlist_tracks(playlist_id);
CREATE INDEX idx_playlist_tracks_position ON playlist_tracks(playlist_id, position);
CREATE INDEX idx_playlist_tracks_spotify_id ON playlist_tracks(spotify_id);
CREATE INDEX idx_playlist_tracks_added_by ON playlist_tracks(added_by_user_id);
CREATE INDEX idx_user_preferences_user_id ON user_preferences(user_id);
CREATE INDEX idx_playlists_public ON playlists(is_public) WHERE is_public = true;
CREATE INDEX idx_playlists_creator_public ON playlists(creator_id, is_public);

-- Add updated_at trigger for new tables
CREATE TRIGGER update_user_preferences_updated_at 
BEFORE UPDATE ON user_preferences 
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Update existing indexes that may have been affected
DROP INDEX IF EXISTS idx_albums_artist_id;
DROP INDEX IF EXISTS idx_albums_spotify_id;
DROP INDEX IF EXISTS idx_tracks_album_id;
DROP INDEX IF EXISTS idx_tracks_spotify_id;
DROP INDEX IF EXISTS idx_reviews_album_id;
DROP INDEX IF EXISTS idx_reviews_track_id;
DROP INDEX IF EXISTS idx_track_featured_artists_track_id;
DROP INDEX IF EXISTS idx_track_featured_artists_artist_id;

-- Remove old playlist_tracks index that's no longer relevant
DROP INDEX IF EXISTS idx_playlist_tracks_playlist_id;
DROP INDEX IF EXISTS idx_playlist_tracks_position; 