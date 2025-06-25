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

type albumRepository struct {
	db *database.PostgresDB
}

func NewAlbumRepository(db *database.PostgresDB) repository.AlbumRepository {
	return &albumRepository{db: db}
}

func (r *albumRepository) Create(ctx context.Context, album *models.Album) error {
	query := `
		INSERT INTO albums (id, spotify_id, title, artist_id, release_date, cover_image, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err := r.db.Pool.Exec(ctx, query,
		album.ID, album.SpotifyID, album.Title, album.ArtistID,
		album.ReleaseDate, album.CoverImage, album.CreatedAt, album.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create album: %w", err)
	}
	
	return nil
}

func (r *albumRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Album, error) {
	query := `
		SELECT id, spotify_id, title, artist_id, release_date, cover_image, created_at, updated_at
		FROM albums 
		WHERE id = $1
	`
	
	album := &models.Album{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&album.ID, &album.SpotifyID, &album.Title, &album.ArtistID,
		&album.ReleaseDate, &album.CoverImage, &album.CreatedAt, &album.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("album not found")
		}
		return nil, fmt.Errorf("failed to get album: %w", err)
	}
	
	return album, nil
}

func (r *albumRepository) GetBySpotifyID(ctx context.Context, spotifyID string) (*models.Album, error) {
	query := `
		SELECT id, spotify_id, title, artist_id, release_date, cover_image, created_at, updated_at
		FROM albums 
		WHERE spotify_id = $1
	`
	
	album := &models.Album{}
	err := r.db.Pool.QueryRow(ctx, query, spotifyID).Scan(
		&album.ID, &album.SpotifyID, &album.Title, &album.ArtistID,
		&album.ReleaseDate, &album.CoverImage, &album.CreatedAt, &album.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("album not found")
		}
		return nil, fmt.Errorf("failed to get album: %w", err)
	}
	
	return album, nil
}

func (r *albumRepository) GetByArtistID(ctx context.Context, artistID uuid.UUID, limit, offset int) ([]*models.Album, error) {
	query := `
		SELECT id, spotify_id, title, artist_id, release_date, cover_image, created_at, updated_at
		FROM albums 
		WHERE artist_id = $1
		ORDER BY release_date DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Pool.Query(ctx, query, artistID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list albums by artist: %w", err)
	}
	defer rows.Close()
	
	var albums []*models.Album
	for rows.Next() {
		album := &models.Album{}
		err := rows.Scan(
			&album.ID, &album.SpotifyID, &album.Title, &album.ArtistID,
			&album.ReleaseDate, &album.CoverImage, &album.CreatedAt, &album.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan album: %w", err)
		}
		albums = append(albums, album)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating albums: %w", err)
	}
	
	return albums, nil
}

func (r *albumRepository) Update(ctx context.Context, album *models.Album) error {
	query := `
		UPDATE albums 
		SET spotify_id = $2, title = $3, artist_id = $4, release_date = $5, cover_image = $6, updated_at = NOW()
		WHERE id = $1
	`
	
	result, err := r.db.Pool.Exec(ctx, query,
		album.ID, album.SpotifyID, album.Title, album.ArtistID, album.ReleaseDate, album.CoverImage,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update album: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("album not found")
	}
	
	return nil
}

func (r *albumRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM albums WHERE id = $1`
	
	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete album: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("album not found")
	}
	
	return nil
}

func (r *albumRepository) List(ctx context.Context, limit, offset int) ([]*models.Album, error) {
	query := `
		SELECT id, spotify_id, title, artist_id, release_date, cover_image, created_at, updated_at
		FROM albums 
		ORDER BY release_date DESC, created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list albums: %w", err)
	}
	defer rows.Close()
	
	var albums []*models.Album
	for rows.Next() {
		album := &models.Album{}
		err := rows.Scan(
			&album.ID, &album.SpotifyID, &album.Title, &album.ArtistID,
			&album.ReleaseDate, &album.CoverImage, &album.CreatedAt, &album.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan album: %w", err)
		}
		albums = append(albums, album)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating albums: %w", err)
	}
	
	return albums, nil
} 