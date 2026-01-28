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

type spotifyTokenRepository struct {
	db *database.PostgresDB
}

func NewSpotifyTokenRepository(db *database.PostgresDB) repository.SpotifyTokenRepository {
	return &spotifyTokenRepository{db: db}
}

func (r *spotifyTokenRepository) Upsert(ctx context.Context, token *models.SpotifyToken) error {
	query := `
		INSERT INTO spotify_tokens (
			user_id, spotify_user_id, access_token, refresh_token, token_type, scope, expires_at, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id) DO UPDATE SET
			spotify_user_id = EXCLUDED.spotify_user_id,
			access_token = EXCLUDED.access_token,
			refresh_token = EXCLUDED.refresh_token,
			token_type = EXCLUDED.token_type,
			scope = EXCLUDED.scope,
			expires_at = EXCLUDED.expires_at,
			updated_at = NOW()
	`

	_, err := r.db.Pool.Exec(ctx, query,
		token.UserID, token.SpotifyUserID, token.AccessToken, token.RefreshToken,
		token.TokenType, token.Scope, token.ExpiresAt, token.CreatedAt, token.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert spotify token: %w", err)
	}

	return nil
}

func (r *spotifyTokenRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.SpotifyToken, error) {
	query := `
		SELECT user_id, spotify_user_id, access_token, refresh_token, token_type, scope, expires_at, created_at, updated_at
		FROM spotify_tokens
		WHERE user_id = $1
	`

	token := &models.SpotifyToken{}
	if err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&token.UserID, &token.SpotifyUserID, &token.AccessToken, &token.RefreshToken,
		&token.TokenType, &token.Scope, &token.ExpiresAt, &token.CreatedAt, &token.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("spotify token not found")
		}
		return nil, fmt.Errorf("failed to get spotify token: %w", err)
	}

	return token, nil
}

func (r *spotifyTokenRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM spotify_tokens WHERE user_id = $1`

	_, err := r.db.Pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete spotify token: %w", err)
	}

	return nil
}
