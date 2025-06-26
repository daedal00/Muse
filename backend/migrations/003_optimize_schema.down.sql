-- Revert schema optimization migration
-- This will recreate the old schema structure

-- Drop new tables and columns
DROP TABLE IF EXISTS user_preferences CASCADE;
DROP TABLE IF EXISTS playlist_tracks CASCADE;

ALTER TABLE playlists DROP COLUMN IF EXISTS is_public;

-- Revert reviews table
ALTER TABLE reviews 
DROP COLUMN IF EXISTS spotify_id,
DROP COLUMN IF EXISTS spotify_type,
ADD COLUMN album_id UUID REFERENCES albums(id) ON DELETE CASCADE,
ADD COLUMN track_id UUID REFERENCES tracks(id) ON DELETE CASCADE,
ADD CONSTRAINT reviews_check CHECK ((album_id IS NOT NULL AND track_id IS NULL) OR (album_id IS NULL AND track_id IS NOT NULL));

-- Recreate old tables
CREATE TABLE artists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    spotify_id VARCHAR(255) UNIQUE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE albums (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    spotify_id VARCHAR(255) UNIQUE,
    title VARCHAR(255) NOT NULL,
    artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
    release_date DATE,
    cover_image VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE tracks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    spotify_id VARCHAR(255) UNIQUE,
    title VARCHAR(255) NOT NULL,
    album_id UUID NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
    duration_ms INTEGER,
    track_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE track_featured_artists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(track_id, artist_id)
);

CREATE TABLE playlist_tracks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    playlist_id UUID NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    position INTEGER NOT NULL,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(playlist_id, track_id)
);

-- Recreate old indexes
CREATE INDEX idx_albums_artist_id ON albums(artist_id);
CREATE INDEX idx_albums_spotify_id ON albums(spotify_id);
CREATE INDEX idx_tracks_album_id ON tracks(album_id);
CREATE INDEX idx_tracks_spotify_id ON tracks(spotify_id);
CREATE INDEX idx_reviews_album_id ON reviews(album_id);
CREATE INDEX idx_reviews_track_id ON reviews(track_id);
CREATE INDEX idx_track_featured_artists_track_id ON track_featured_artists(track_id);
CREATE INDEX idx_track_featured_artists_artist_id ON track_featured_artists(artist_id);
CREATE INDEX idx_playlist_tracks_playlist_id ON playlist_tracks(playlist_id);
CREATE INDEX idx_playlist_tracks_position ON playlist_tracks(playlist_id, position);

-- Recreate triggers
CREATE TRIGGER update_artists_updated_at BEFORE UPDATE ON artists FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_albums_updated_at BEFORE UPDATE ON albums FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tracks_updated_at BEFORE UPDATE ON tracks FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); 