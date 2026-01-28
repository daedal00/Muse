-- Migration: Add profile_settings table for customizable user profiles
-- This enables users to personalize their public profile appearance

CREATE TABLE IF NOT EXISTS profile_settings (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    
    -- Layout configuration
    layout TEXT NOT NULL DEFAULT 'GRID' CHECK (layout IN ('GRID', 'LIST', 'BENTO')),
    
    -- Theme colors
    primary_color TEXT DEFAULT '#1DB954',  -- Default Spotify green
    accent_color TEXT DEFAULT '#191414',   -- Default Spotify black
    
    -- Background customization
    background_style TEXT DEFAULT 'SOLID' CHECK (background_style IN ('SOLID', 'GRADIENT', 'IMAGE')),
    background_value TEXT DEFAULT '#0f172a',  -- slate-900
    
    -- Pinned content (arrays of IDs)
    pinned_album_ids UUID[] DEFAULT ARRAY[]::UUID[],
    pinned_track_ids UUID[] DEFAULT ARRAY[]::UUID[],
    featured_artist_ids UUID[] DEFAULT ARRAY[]::UUID[],
    
    -- Section ordering
    sections_order TEXT[] DEFAULT ARRAY['featuredArtists', 'topTracks', 'pinnedAlbums', 'recentReviews', 'playlists']::TEXT[],
    
    -- Visibility toggles
    show_spotify_stats BOOLEAN DEFAULT true,
    show_listening_history BOOLEAN DEFAULT true,
    
    -- Bio style
    bio_style TEXT DEFAULT 'MINIMAL' CHECK (bio_style IN ('MINIMAL', 'DETAILED', 'QUOTE')),
    
    -- Custom genre/vibe tags
    custom_tags TEXT[] DEFAULT ARRAY[]::TEXT[],
    
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_profile_settings_user_id ON profile_settings(user_id);

-- Add trigger to auto-update updated_at
CREATE OR REPLACE FUNCTION update_profile_settings_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_profile_settings_updated_at
    BEFORE UPDATE ON profile_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_profile_settings_updated_at();

-- Add comment for documentation
COMMENT ON TABLE profile_settings IS 'User profile customization settings for public profile display';
