CREATE TABLE spotify_tokens (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    spotify_user_id VARCHAR(255),
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    token_type VARCHAR(50),
    scope TEXT,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_spotify_tokens_expires_at ON spotify_tokens(expires_at);

CREATE TRIGGER update_spotify_tokens_updated_at BEFORE UPDATE ON spotify_tokens FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
