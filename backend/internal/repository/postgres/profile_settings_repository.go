package postgres

import (
	"context"
	"fmt"

	"github.com/daedal00/muse/backend/internal/database"
	"github.com/daedal00/muse/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

// ProfileSettingsRepository handles database operations for profile settings
type ProfileSettingsRepository struct {
	db *database.PostgresDB
}

// NewProfileSettingsRepository creates a new profile settings repository
func NewProfileSettingsRepository(db *database.PostgresDB) *ProfileSettingsRepository {
	return &ProfileSettingsRepository{db: db}
}

// GetByUserID retrieves profile settings for a user, returning defaults if not found
func (r *ProfileSettingsRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.ProfileSettings, error) {
	query := `
		SELECT 
			user_id, layout, primary_color, accent_color,
			background_style, background_value,
			pinned_album_ids, pinned_track_ids, featured_artist_ids,
			sections_order, show_spotify_stats, show_listening_history,
			bio_style, custom_tags, created_at, updated_at
		FROM profile_settings
		WHERE user_id = $1
	`

	settings := &models.ProfileSettings{}
	var (
		pinnedAlbumIDs    []string
		pinnedTrackIDs    []string
		featuredArtistIDs []string
		sectionsOrder     []string
		customTags        []string
		bgStyle           *string
		bioStyle          *string
	)

	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&settings.UserID,
		&settings.Layout,
		&settings.PrimaryColor,
		&settings.AccentColor,
		&bgStyle,
		&settings.BackgroundValue,
		&pinnedAlbumIDs,
		&pinnedTrackIDs,
		&featuredArtistIDs,
		&sectionsOrder,
		&settings.ShowSpotifyStats,
		&settings.ShowListeningHistory,
		&bioStyle,
		&customTags,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		// Return default settings if none exist
		return models.DefaultProfileSettings(userID), nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get profile settings: %w", err)
	}

	// Convert string arrays to UUID arrays
	settings.PinnedAlbumIDs = parseUUIDs(pinnedAlbumIDs)
	settings.PinnedTrackIDs = parseUUIDs(pinnedTrackIDs)
	settings.FeaturedArtistIDs = parseUUIDs(featuredArtistIDs)
	settings.SectionsOrder = sectionsOrder
	settings.CustomTags = customTags

	// Convert string to enum types
	if bgStyle != nil {
		bs := models.BackgroundStyle(*bgStyle)
		settings.BackgroundStyle = &bs
	}
	if bioStyle != nil {
		bs := models.BioStyle(*bioStyle)
		settings.BioStyle = &bs
	}

	return settings, nil
}

// Upsert creates or updates profile settings for a user
func (r *ProfileSettingsRepository) Upsert(ctx context.Context, settings *models.ProfileSettings) error {
	query := `
		INSERT INTO profile_settings (
			user_id, layout, primary_color, accent_color,
			background_style, background_value,
			pinned_album_ids, pinned_track_ids, featured_artist_ids,
			sections_order, show_spotify_stats, show_listening_history,
			bio_style, custom_tags
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (user_id) DO UPDATE SET
			layout = EXCLUDED.layout,
			primary_color = EXCLUDED.primary_color,
			accent_color = EXCLUDED.accent_color,
			background_style = EXCLUDED.background_style,
			background_value = EXCLUDED.background_value,
			pinned_album_ids = EXCLUDED.pinned_album_ids,
			pinned_track_ids = EXCLUDED.pinned_track_ids,
			featured_artist_ids = EXCLUDED.featured_artist_ids,
			sections_order = EXCLUDED.sections_order,
			show_spotify_stats = EXCLUDED.show_spotify_stats,
			show_listening_history = EXCLUDED.show_listening_history,
			bio_style = EXCLUDED.bio_style,
			custom_tags = EXCLUDED.custom_tags,
			updated_at = NOW()
	`

	// Convert UUID arrays to string arrays for PostgreSQL
	pinnedAlbumIDs := uuidsToStrings(settings.PinnedAlbumIDs)
	pinnedTrackIDs := uuidsToStrings(settings.PinnedTrackIDs)
	featuredArtistIDs := uuidsToStrings(settings.FeaturedArtistIDs)

	// Convert enum types to strings
	var bgStyle, bioStyle *string
	if settings.BackgroundStyle != nil {
		s := string(*settings.BackgroundStyle)
		bgStyle = &s
	}
	if settings.BioStyle != nil {
		s := string(*settings.BioStyle)
		bioStyle = &s
	}

	_, err := r.db.Pool.Exec(ctx, query,
		settings.UserID,
		string(settings.Layout),
		settings.PrimaryColor,
		settings.AccentColor,
		bgStyle,
		settings.BackgroundValue,
		pq.Array(pinnedAlbumIDs),
		pq.Array(pinnedTrackIDs),
		pq.Array(featuredArtistIDs),
		pq.Array(settings.SectionsOrder),
		settings.ShowSpotifyStats,
		settings.ShowListeningHistory,
		bioStyle,
		pq.Array(settings.CustomTags),
	)

	if err != nil {
		return fmt.Errorf("failed to upsert profile settings: %w", err)
	}

	return nil
}

// Delete removes profile settings for a user
func (r *ProfileSettingsRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM profile_settings WHERE user_id = $1`

	result, err := r.db.Pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete profile settings: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("profile settings not found for user: %s", userID)
	}

	return nil
}

// Helper functions

func parseUUIDs(strs []string) []uuid.UUID {
	result := make([]uuid.UUID, 0, len(strs))
	for _, s := range strs {
		if id, err := uuid.Parse(s); err == nil {
			result = append(result, id)
		}
	}
	return result
}

func uuidsToStrings(uuids []uuid.UUID) []string {
	result := make([]string, len(uuids))
	for i, id := range uuids {
		result[i] = id.String()
	}
	return result
}
