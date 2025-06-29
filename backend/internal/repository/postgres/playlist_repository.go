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

type playlistRepository struct {
	db *database.PostgresDB
}

func NewPlaylistRepository(db *database.PostgresDB) repository.PlaylistRepository {
	return &playlistRepository{db: db}
}

func (r *playlistRepository) Create(ctx context.Context, playlist *models.Playlist) error {
	query := `
		INSERT INTO playlists (id, title, description, cover_image, creator_id, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		playlist.ID, playlist.Title, playlist.Description, playlist.CoverImage,
		playlist.CreatorID, playlist.IsPublic, playlist.CreatedAt, playlist.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create playlist: %w", err)
	}

	return nil
}

func (r *playlistRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Playlist, error) {
	query := `
		SELECT p.id, p.title, p.description, p.cover_image, p.creator_id, p.is_public, p.created_at, p.updated_at,
		       u.id, u.name, u.email, u.bio, u.avatar, u.created_at, u.updated_at
		FROM playlists p
		LEFT JOIN users u ON p.creator_id = u.id
		WHERE p.id = $1
	`

	playlist := &models.Playlist{}
	creator := &models.User{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&playlist.ID, &playlist.Title, &playlist.Description, &playlist.CoverImage,
		&playlist.CreatorID, &playlist.IsPublic, &playlist.CreatedAt, &playlist.UpdatedAt,
		&creator.ID, &creator.Name, &creator.Email, &creator.Bio, &creator.Avatar,
		&creator.CreatedAt, &creator.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("playlist not found")
		}
		return nil, fmt.Errorf("failed to get playlist: %w", err)
	}

	playlist.Creator = creator
	return playlist, nil
}

func (r *playlistRepository) GetByCreatorID(ctx context.Context, creatorID uuid.UUID, limit, offset int) ([]*models.Playlist, error) {
	query := `
		SELECT p.id, p.title, p.description, p.cover_image, p.creator_id, p.is_public, p.created_at, p.updated_at,
		       u.id, u.name, u.email, u.bio, u.avatar, u.created_at, u.updated_at
		FROM playlists p
		LEFT JOIN users u ON p.creator_id = u.id
		WHERE p.creator_id = $1
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, creatorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list playlists by creator: %w", err)
	}
	defer rows.Close()

	var playlists []*models.Playlist
	for rows.Next() {
		playlist := &models.Playlist{}
		creator := &models.User{}
		err := rows.Scan(
			&playlist.ID, &playlist.Title, &playlist.Description, &playlist.CoverImage,
			&playlist.CreatorID, &playlist.IsPublic, &playlist.CreatedAt, &playlist.UpdatedAt,
			&creator.ID, &creator.Name, &creator.Email, &creator.Bio, &creator.Avatar,
			&creator.CreatedAt, &creator.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan playlist: %w", err)
		}
		playlist.Creator = creator
		playlists = append(playlists, playlist)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating playlists: %w", err)
	}

	return playlists, nil
}

func (r *playlistRepository) Update(ctx context.Context, playlist *models.Playlist) error {
	query := `
		UPDATE playlists 
		SET title = $2, description = $3, cover_image = $4, is_public = $5, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query,
		playlist.ID, playlist.Title, playlist.Description, playlist.CoverImage, playlist.IsPublic,
	)

	if err != nil {
		return fmt.Errorf("failed to update playlist: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("playlist not found")
	}

	return nil
}

func (r *playlistRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Start a transaction to delete playlist tracks first, then the playlist
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Delete all playlist tracks first
	_, err = tx.Exec(ctx, `DELETE FROM playlist_tracks WHERE playlist_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete playlist tracks: %w", err)
	}

	// Delete the playlist
	result, err := tx.Exec(ctx, `DELETE FROM playlists WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete playlist: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("playlist not found")
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *playlistRepository) List(ctx context.Context, limit, offset int) ([]*models.Playlist, error) {
	query := `
		SELECT p.id, p.title, p.description, p.cover_image, p.creator_id, p.is_public, p.created_at, p.updated_at,
		       u.id, u.name, u.email, u.bio, u.avatar, u.created_at, u.updated_at
		FROM playlists p
		LEFT JOIN users u ON p.creator_id = u.id
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list playlists: %w", err)
	}
	defer rows.Close()

	var playlists []*models.Playlist
	for rows.Next() {
		playlist := &models.Playlist{}
		creator := &models.User{}
		err := rows.Scan(
			&playlist.ID, &playlist.Title, &playlist.Description, &playlist.CoverImage,
			&playlist.CreatorID, &playlist.IsPublic, &playlist.CreatedAt, &playlist.UpdatedAt,
			&creator.ID, &creator.Name, &creator.Email, &creator.Bio, &creator.Avatar,
			&creator.CreatedAt, &creator.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan playlist: %w", err)
		}
		playlist.Creator = creator
		playlists = append(playlists, playlist)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating playlists: %w", err)
	}

	return playlists, nil
}

// Playlist track operations

func (r *playlistRepository) AddTrack(ctx context.Context, playlistID uuid.UUID, spotifyID string, position int, addedByUserID uuid.UUID) error {
	// If position is 0 or negative, append to the end
	if position <= 0 {
		var maxPosition int
		query := `SELECT COALESCE(MAX(position), 0) FROM playlist_tracks WHERE playlist_id = $1`
		err := r.db.Pool.QueryRow(ctx, query, playlistID).Scan(&maxPosition)
		if err != nil {
			return fmt.Errorf("failed to get max position: %w", err)
		}
		position = maxPosition + 1
	} else {
		// Shift existing tracks to make room
		shiftQuery := `UPDATE playlist_tracks SET position = position + 1 WHERE playlist_id = $1 AND position >= $2`
		_, err := r.db.Pool.Exec(ctx, shiftQuery, playlistID, position)
		if err != nil {
			return fmt.Errorf("failed to shift track positions: %w", err)
		}
	}

	insertQuery := `
		INSERT INTO playlist_tracks (id, playlist_id, spotify_id, position, added_at, added_by_user_id)
		VALUES ($1, $2, $3, $4, NOW(), $5)
	`

	_, err := r.db.Pool.Exec(ctx, insertQuery, uuid.New(), playlistID, spotifyID, position, addedByUserID)
	if err != nil {
		return fmt.Errorf("failed to add track to playlist: %w", err)
	}

	return nil
}

func (r *playlistRepository) RemoveTrack(ctx context.Context, playlistID uuid.UUID, spotifyID string) error {
	// Get the position of the track being removed
	var position int
	getPositionQuery := `SELECT position FROM playlist_tracks WHERE playlist_id = $1 AND spotify_id = $2`
	err := r.db.Pool.QueryRow(ctx, getPositionQuery, playlistID, spotifyID).Scan(&position)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("track not found in playlist")
		}
		return fmt.Errorf("failed to get track position: %w", err)
	}

	// Start a transaction to remove track and shift positions
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Remove the track
	removeQuery := `DELETE FROM playlist_tracks WHERE playlist_id = $1 AND spotify_id = $2`
	result, err := tx.Exec(ctx, removeQuery, playlistID, spotifyID)
	if err != nil {
		return fmt.Errorf("failed to remove track from playlist: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("track not found in playlist")
	}

	// Shift remaining tracks down
	shiftQuery := `UPDATE playlist_tracks SET position = position - 1 WHERE playlist_id = $1 AND position > $2`
	_, err = tx.Exec(ctx, shiftQuery, playlistID, position)
	if err != nil {
		return fmt.Errorf("failed to shift track positions: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *playlistRepository) GetTrackCount(ctx context.Context, playlistID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM playlist_tracks WHERE playlist_id = $1`

	var count int
	err := r.db.Pool.QueryRow(ctx, query, playlistID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get playlist track count: %w", err)
	}

	return count, nil
}

func (r *playlistRepository) GetTrackSpotifyIDs(ctx context.Context, playlistID uuid.UUID, limit, offset int) ([]string, error) {
	query := `
		SELECT spotify_id
		FROM playlist_tracks
		WHERE playlist_id = $1
		ORDER BY position ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, playlistID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist track IDs: %w", err)
	}
	defer rows.Close()

	var spotifyIDs []string
	for rows.Next() {
		var spotifyID string
		err := rows.Scan(&spotifyID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan spotify ID: %w", err)
		}
		spotifyIDs = append(spotifyIDs, spotifyID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating spotify IDs: %w", err)
	}

	return spotifyIDs, nil
}

func (r *playlistRepository) GetPlaylistTracks(ctx context.Context, playlistID uuid.UUID, limit, offset int) ([]*models.PlaylistTrack, error) {
	query := `
		SELECT pt.id, pt.playlist_id, pt.spotify_id, pt.position, pt.added_at, pt.added_by_user_id,
		       u.id, u.name, u.email, u.bio, u.avatar, u.created_at, u.updated_at
		FROM playlist_tracks pt
		LEFT JOIN users u ON pt.added_by_user_id = u.id
		WHERE pt.playlist_id = $1
		ORDER BY pt.position ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, playlistID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist tracks: %w", err)
	}
	defer rows.Close()

	var playlistTracks []*models.PlaylistTrack
	for rows.Next() {
		track := &models.PlaylistTrack{}
		user := &models.User{}
		err := rows.Scan(
			&track.ID, &track.PlaylistID, &track.SpotifyID, &track.Position, &track.AddedAt, &track.AddedByUserID,
			&user.ID, &user.Name, &user.Email, &user.Bio, &user.Avatar, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan playlist track: %w", err)
		}
		track.AddedByUser = user
		playlistTracks = append(playlistTracks, track)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating playlist tracks: %w", err)
	}

	return playlistTracks, nil
}

func (r *playlistRepository) ReorderTracks(ctx context.Context, playlistID uuid.UUID, trackPositions map[string]int) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for spotifyID, position := range trackPositions {
		updateQuery := `UPDATE playlist_tracks SET position = $1 WHERE playlist_id = $2 AND spotify_id = $3`
		_, err := tx.Exec(ctx, updateQuery, position, playlistID, spotifyID)
		if err != nil {
			return fmt.Errorf("failed to update track position: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *playlistRepository) GetPublicPlaylists(ctx context.Context, limit, offset int) ([]*models.Playlist, error) {
	query := `
		SELECT p.id, p.title, p.description, p.cover_image, p.creator_id, p.is_public, p.created_at, p.updated_at,
		       u.id, u.name, u.email, u.bio, u.avatar, u.created_at, u.updated_at
		FROM playlists p
		LEFT JOIN users u ON p.creator_id = u.id
		WHERE p.is_public = true
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get public playlists: %w", err)
	}
	defer rows.Close()

	var playlists []*models.Playlist
	for rows.Next() {
		playlist := &models.Playlist{}
		creator := &models.User{}
		err := rows.Scan(
			&playlist.ID, &playlist.Title, &playlist.Description, &playlist.CoverImage,
			&playlist.CreatorID, &playlist.IsPublic, &playlist.CreatedAt, &playlist.UpdatedAt,
			&creator.ID, &creator.Name, &creator.Email, &creator.Bio, &creator.Avatar,
			&creator.CreatedAt, &creator.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan playlist: %w", err)
		}
		playlist.Creator = creator
		playlists = append(playlists, playlist)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating playlists: %w", err)
	}

	return playlists, nil
}

func (r *playlistRepository) GetTracksByPlaylistID(ctx context.Context, playlistID uuid.UUID, limit, offset int) ([]*models.PlaylistTrack, error) {
	query := `
		SELECT pt.id, pt.playlist_id, pt.spotify_id, pt.position, pt.added_at, pt.added_by_user_id
		FROM playlist_tracks pt
		WHERE pt.playlist_id = $1
		ORDER BY pt.position ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, playlistID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist tracks: %w", err)
	}
	defer rows.Close()

	var tracks []*models.PlaylistTrack
	for rows.Next() {
		track := &models.PlaylistTrack{}
		err := rows.Scan(
			&track.ID, &track.PlaylistID, &track.SpotifyID, &track.Position,
			&track.AddedAt, &track.AddedByUserID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan playlist track: %w", err)
		}
		tracks = append(tracks, track)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating playlist tracks: %w", err)
	}

	return tracks, nil
}
