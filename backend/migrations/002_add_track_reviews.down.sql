DROP INDEX IF EXISTS idx_playlists_spotify_id;
ALTER TABLE playlists DROP COLUMN IF EXISTS spotify_id;

DROP TRIGGER IF EXISTS update_track_reviews_updated_at ON track_reviews;
DROP INDEX IF EXISTS idx_track_reviews_track_id;
DROP INDEX IF EXISTS idx_track_reviews_user_id;
DROP TABLE IF EXISTS track_reviews;
