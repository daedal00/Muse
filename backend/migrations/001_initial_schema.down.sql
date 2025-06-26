-- Drop triggers
DROP TRIGGER IF EXISTS update_playlists_updated_at ON playlists;
DROP TRIGGER IF EXISTS update_reviews_updated_at ON reviews;
DROP TRIGGER IF EXISTS update_tracks_updated_at ON tracks;
DROP TRIGGER IF EXISTS update_albums_updated_at ON albums;
DROP TRIGGER IF EXISTS update_artists_updated_at ON artists;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order (to handle foreign key constraints)
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS playlist_tracks;
DROP TABLE IF EXISTS playlists;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS tracks;
DROP TABLE IF EXISTS albums;
DROP TABLE IF EXISTS artists;
DROP TABLE IF EXISTS users;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp"; 