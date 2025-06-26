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

type userRepository struct {
	db *database.PostgresDB
}

func NewUserRepository(db *database.PostgresDB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, name, email, password_hash, spotify_id, spotify_access_token, spotify_refresh_token, spotify_token_expiry, bio, avatar, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		user.ID, user.Name, user.Email, user.PasswordHash, user.SpotifyID, user.SpotifyAccessToken, user.SpotifyRefreshToken, user.SpotifyTokenExpiry,
		user.Bio, user.Avatar, user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, name, email, password_hash, spotify_id, spotify_access_token, spotify_refresh_token, spotify_token_expiry, bio, avatar, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.SpotifyID, &user.SpotifyAccessToken, &user.SpotifyRefreshToken, &user.SpotifyTokenExpiry,
		&user.Bio, &user.Avatar, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, name, email, password_hash, spotify_id, spotify_access_token, spotify_refresh_token, spotify_token_expiry, bio, avatar, created_at, updated_at
		FROM users 
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.SpotifyID, &user.SpotifyAccessToken, &user.SpotifyRefreshToken, &user.SpotifyTokenExpiry,
		&user.Bio, &user.Avatar, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetBySpotifyID(ctx context.Context, spotifyID string) (*models.User, error) {
	query := `
		SELECT id, name, email, password_hash, spotify_id, spotify_access_token, spotify_refresh_token, spotify_token_expiry, bio, avatar, created_at, updated_at
		FROM users 
		WHERE spotify_id = $1
	`

	user := &models.User{}
	err := r.db.Pool.QueryRow(ctx, query, spotifyID).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.SpotifyID, &user.SpotifyAccessToken, &user.SpotifyRefreshToken, &user.SpotifyTokenExpiry,
		&user.Bio, &user.Avatar, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by spotify ID: %w", err)
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET name = $2, email = $3, password_hash = $4, spotify_id = $5, spotify_access_token = $6, spotify_refresh_token = $7, spotify_token_expiry = $8, bio = $9, avatar = $10, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query,
		user.ID, user.Name, user.Email, user.PasswordHash, user.SpotifyID, user.SpotifyAccessToken, user.SpotifyRefreshToken, user.SpotifyTokenExpiry, user.Bio, user.Avatar,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	query := `
		SELECT id, name, email, password_hash, spotify_id, spotify_access_token, spotify_refresh_token, spotify_token_expiry, bio, avatar, created_at, updated_at
		FROM users 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.SpotifyID, &user.SpotifyAccessToken, &user.SpotifyRefreshToken, &user.SpotifyTokenExpiry,
			&user.Bio, &user.Avatar, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}
