DROP TRIGGER IF EXISTS update_spotify_tokens_updated_at ON spotify_tokens;
DROP INDEX IF EXISTS idx_spotify_tokens_expires_at;
DROP TABLE IF EXISTS spotify_tokens;
