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

type artistRepository struct {
	db *database.PostgresDB
}

func NewArtistRepository(db *database.PostgresDB) repository.ArtistRepository {
	return &artistRepository{db: db}
}

func (r *artistRepository) Create(ctx context.Context, artist *models.Artist) error {
	query := `
		INSERT INTO artists (id, spotify_id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	
	_, err := r.db.Pool.Exec(ctx, query,
		artist.ID, artist.SpotifyID, artist.Name, artist.CreatedAt, artist.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create artist: %w", err)
	}
	
	return nil
}

func (r *artistRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Artist, error) {
	query := `
		SELECT id, spotify_id, name, created_at, updated_at
		FROM artists 
		WHERE id = $1
	`
	
	artist := &models.Artist{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&artist.ID, &artist.SpotifyID, &artist.Name, &artist.CreatedAt, &artist.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("artist not found")
		}
		return nil, fmt.Errorf("failed to get artist: %w", err)
	}
	
	return artist, nil
}

func (r *artistRepository) GetBySpotifyID(ctx context.Context, spotifyID string) (*models.Artist, error) {
	query := `
		SELECT id, spotify_id, name, created_at, updated_at
		FROM artists 
		WHERE spotify_id = $1
	`
	
	artist := &models.Artist{}
	err := r.db.Pool.QueryRow(ctx, query, spotifyID).Scan(
		&artist.ID, &artist.SpotifyID, &artist.Name, &artist.CreatedAt, &artist.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("artist not found")
		}
		return nil, fmt.Errorf("failed to get artist: %w", err)
	}
	
	return artist, nil
}

func (r *artistRepository) Update(ctx context.Context, artist *models.Artist) error {
	query := `
		UPDATE artists 
		SET spotify_id = $2, name = $3, updated_at = NOW()
		WHERE id = $1
	`
	
	result, err := r.db.Pool.Exec(ctx, query,
		artist.ID, artist.SpotifyID, artist.Name,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update artist: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("artist not found")
	}
	
	return nil
}

func (r *artistRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM artists WHERE id = $1`
	
	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete artist: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("artist not found")
	}
	
	return nil
}

func (r *artistRepository) List(ctx context.Context, limit, offset int) ([]*models.Artist, error) {
	query := `
		SELECT id, spotify_id, name, created_at, updated_at
		FROM artists 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list artists: %w", err)
	}
	defer rows.Close()
	
	var artists []*models.Artist
	for rows.Next() {
		artist := &models.Artist{}
		err := rows.Scan(
			&artist.ID, &artist.SpotifyID, &artist.Name, &artist.CreatedAt, &artist.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan artist: %w", err)
		}
		artists = append(artists, artist)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating artists: %w", err)
	}
	
	return artists, nil
} 