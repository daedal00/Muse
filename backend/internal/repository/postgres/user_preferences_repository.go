package postgres

import (
	"context"
	"fmt"

	"github.com/daedal00/muse/backend/internal/database"
	"github.com/daedal00/muse/backend/internal/models"
	"github.com/daedal00/muse/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type userPreferencesRepository struct {
	db *database.PostgresDB
}

func NewUserPreferencesRepository(db *database.PostgresDB) repository.UserPreferencesRepository {
	return &userPreferencesRepository{db: db}
}

func (r *userPreferencesRepository) Create(ctx context.Context, preferences *models.UserPreferences) error {
	preferences.ID = uuid.New()

	query := `
		INSERT INTO user_preferences (id, user_id, preferred_genres, favorite_artist_ids, notification_settings, privacy_settings, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`

	_, err := r.db.Pool.Exec(ctx, query,
		preferences.ID,
		preferences.UserID,
		preferences.PreferredGenres,
		preferences.FavoriteArtistIDs,
		preferences.NotificationSettings,
		preferences.PrivacySettings,
	)
	if err != nil {
		return fmt.Errorf("failed to create user preferences: %w", err)
	}

	return nil
}

func (r *userPreferencesRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserPreferences, error) {
	query := `
		SELECT id, user_id, preferred_genres, favorite_artist_ids, notification_settings, privacy_settings, created_at, updated_at
		FROM user_preferences
		WHERE user_id = $1
	`

	preferences := &models.UserPreferences{}
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&preferences.ID,
		&preferences.UserID,
		&preferences.PreferredGenres,
		&preferences.FavoriteArtistIDs,
		&preferences.NotificationSettings,
		&preferences.PrivacySettings,
		&preferences.CreatedAt,
		&preferences.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No preferences found
		}
		return nil, fmt.Errorf("failed to get user preferences: %w", err)
	}

	return preferences, nil
}

func (r *userPreferencesRepository) Update(ctx context.Context, preferences *models.UserPreferences) error {
	query := `
		UPDATE user_preferences 
		SET preferred_genres = $1, favorite_artist_ids = $2, notification_settings = $3, privacy_settings = $4, updated_at = NOW()
		WHERE user_id = $5
	`

	result, err := r.db.Pool.Exec(ctx, query,
		preferences.PreferredGenres,
		preferences.FavoriteArtistIDs,
		preferences.NotificationSettings,
		preferences.PrivacySettings,
		preferences.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user preferences: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user preferences not found")
	}

	return nil
}

func (r *userPreferencesRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_preferences WHERE user_id = $1`

	result, err := r.db.Pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user preferences: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user preferences not found")
	}

	return nil
}

func (r *userPreferencesRepository) AddFavoriteArtist(ctx context.Context, userID uuid.UUID, spotifyArtistID string) error {
	// Get current preferences or create if not exists
	preferences, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if preferences == nil {
		// Create new preferences with this artist
		preferences = &models.UserPreferences{
			UserID:               userID,
			PreferredGenres:      []string{},
			FavoriteArtistIDs:    []string{spotifyArtistID},
			NotificationSettings: map[string]bool{},
			PrivacySettings:      map[string]bool{},
		}
		return r.Create(ctx, preferences)
	}

	// Check if artist is already in favorites
	for _, artistID := range preferences.FavoriteArtistIDs {
		if artistID == spotifyArtistID {
			return nil // Already a favorite
		}
	}

	// Add artist to favorites
	preferences.FavoriteArtistIDs = append(preferences.FavoriteArtistIDs, spotifyArtistID)
	return r.Update(ctx, preferences)
}

func (r *userPreferencesRepository) RemoveFavoriteArtist(ctx context.Context, userID uuid.UUID, spotifyArtistID string) error {
	preferences, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if preferences == nil {
		return fmt.Errorf("user preferences not found")
	}

	// Remove artist from favorites
	var newFavorites []string
	for _, artistID := range preferences.FavoriteArtistIDs {
		if artistID != spotifyArtistID {
			newFavorites = append(newFavorites, artistID)
		}
	}

	preferences.FavoriteArtistIDs = newFavorites
	return r.Update(ctx, preferences)
}

func (r *userPreferencesRepository) UpdatePreferredGenres(ctx context.Context, userID uuid.UUID, genres []string) error {
	preferences, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if preferences == nil {
		// Create new preferences with these genres
		preferences = &models.UserPreferences{
			UserID:               userID,
			PreferredGenres:      genres,
			FavoriteArtistIDs:    []string{},
			NotificationSettings: map[string]bool{},
			PrivacySettings:      map[string]bool{},
		}
		return r.Create(ctx, preferences)
	}

	preferences.PreferredGenres = genres
	return r.Update(ctx, preferences)
}
