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

type trackRepository struct {
	db *database.PostgresDB
}

func NewTrackRepository(db *database.PostgresDB) repository.TrackRepository {
	return &trackRepository{db: db}
}

func (r *trackRepository) Create(ctx context.Context, track *models.Track) error {
	query := `
		INSERT INTO tracks (id, spotify_id, title, album_id, duration_ms, track_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		track.ID, track.SpotifyID, track.Title, track.AlbumID,
		track.DurationMs, track.TrackNumber, track.CreatedAt, track.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create track: %w", err)
	}

	return nil
}

func (r *trackRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Track, error) {
	query := `
		SELECT id, spotify_id, title, album_id, duration_ms, track_number, created_at, updated_at
		FROM tracks 
		WHERE id = $1
	`

	track := &models.Track{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&track.ID, &track.SpotifyID, &track.Title, &track.AlbumID,
		&track.DurationMs, &track.TrackNumber, &track.CreatedAt, &track.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("track not found")
		}
		return nil, fmt.Errorf("failed to get track: %w", err)
	}

	return track, nil
}

func (r *trackRepository) GetBySpotifyID(ctx context.Context, spotifyID string) (*models.Track, error) {
	query := `
		SELECT id, spotify_id, title, album_id, duration_ms, track_number, created_at, updated_at
		FROM tracks 
		WHERE spotify_id = $1
	`

	track := &models.Track{}
	err := r.db.Pool.QueryRow(ctx, query, spotifyID).Scan(
		&track.ID, &track.SpotifyID, &track.Title, &track.AlbumID,
		&track.DurationMs, &track.TrackNumber, &track.CreatedAt, &track.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("track not found")
		}
		return nil, fmt.Errorf("failed to get track: %w", err)
	}

	return track, nil
}

func (r *trackRepository) GetByAlbumID(ctx context.Context, albumID uuid.UUID, limit, offset int) ([]*models.Track, error) {
	query := `
		SELECT id, spotify_id, title, album_id, duration_ms, track_number, created_at, updated_at
		FROM tracks 
		WHERE album_id = $1
		ORDER BY track_number ASC, created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, albumID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tracks by album: %w", err)
	}
	defer rows.Close()

	var tracks []*models.Track
	for rows.Next() {
		track := &models.Track{}
		err := rows.Scan(
			&track.ID, &track.SpotifyID, &track.Title, &track.AlbumID,
			&track.DurationMs, &track.TrackNumber, &track.CreatedAt, &track.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan track: %w", err)
		}
		tracks = append(tracks, track)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tracks: %w", err)
	}

	return tracks, nil
}

func (r *trackRepository) Update(ctx context.Context, track *models.Track) error {
	query := `
		UPDATE tracks 
		SET spotify_id = $2, title = $3, album_id = $4, duration_ms = $5, track_number = $6, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query,
		track.ID, track.SpotifyID, track.Title, track.AlbumID, track.DurationMs, track.TrackNumber,
	)

	if err != nil {
		return fmt.Errorf("failed to update track: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("track not found")
	}

	return nil
}

func (r *trackRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tracks WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete track: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("track not found")
	}

	return nil
}

func (r *trackRepository) List(ctx context.Context, limit, offset int) ([]*models.Track, error) {
	query := `
		SELECT id, spotify_id, title, album_id, duration_ms, track_number, created_at, updated_at
		FROM tracks 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tracks: %w", err)
	}
	defer rows.Close()

	var tracks []*models.Track
	for rows.Next() {
		track := &models.Track{}
		err := rows.Scan(
			&track.ID, &track.SpotifyID, &track.Title, &track.AlbumID,
			&track.DurationMs, &track.TrackNumber, &track.CreatedAt, &track.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan track: %w", err)
		}
		tracks = append(tracks, track)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tracks: %w", err)
	}

	return tracks, nil
}
