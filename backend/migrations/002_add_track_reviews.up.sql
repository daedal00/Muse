-- Track reviews table
CREATE TABLE track_reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    review_text TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, track_id)
);

CREATE INDEX idx_track_reviews_user_id ON track_reviews(user_id);
CREATE INDEX idx_track_reviews_track_id ON track_reviews(track_id);

CREATE TRIGGER update_track_reviews_updated_at BEFORE UPDATE ON track_reviews FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Spotify playlists support
ALTER TABLE playlists ADD COLUMN spotify_id VARCHAR(255);
CREATE UNIQUE INDEX idx_playlists_spotify_id ON playlists(spotify_id) WHERE spotify_id IS NOT NULL;
