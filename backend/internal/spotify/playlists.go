package spotify

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zmb3/spotify/v2"

	"github.com/daedal00/muse/backend/internal/models"
)

// PlaylistImportService handles importing playlists from Spotify
type PlaylistImportService struct {
	db          *pgxpool.Pool
	authService *AuthService
}

// NewPlaylistImportService creates a new playlist import service
func NewPlaylistImportService(db *pgxpool.Pool, authService *AuthService) *PlaylistImportService {
	return &PlaylistImportService{
		db:          db,
		authService: authService,
	}
}

// SpotifyPlaylistInfo represents a simplified Spotify playlist for import
type SpotifyPlaylistInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	TrackCount  int    `json:"track_count"`
	IsPublic    bool   `json:"is_public"`
}

// GetUserSpotifyPlaylists retrieves playlists from user's Spotify account with pagination
func (s *PlaylistImportService) GetUserSpotifyPlaylists(ctx context.Context, userID uuid.UUID, limit, offset int) ([]SpotifyPlaylistInfo, int, error) {
	spotifyClient, err := s.authService.GetUserSpotifyClient(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get Spotify client: %w", err)
	}

	// Set reasonable limits
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	// Get current user's playlists with pagination
	playlistPage, err := spotifyClient.CurrentUsersPlaylists(ctx, spotify.Limit(limit), spotify.Offset(offset))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get playlists: %w", err)
	}

	var playlists []SpotifyPlaylistInfo

	// Process current page
	for _, playlist := range playlistPage.Playlists {
		info := SpotifyPlaylistInfo{
			ID:          string(playlist.ID),
			Name:        playlist.Name,
			Description: playlist.Description,
			TrackCount:  int(playlist.Tracks.Total),
			IsPublic:    playlist.IsPublic,
		}

		if len(playlist.Images) > 0 {
			info.Image = playlist.Images[0].URL
		}

		playlists = append(playlists, info)
	}

	return playlists, int(playlistPage.Total), nil
}

// GetUserSpotifyPlaylistsLegacy retrieves all playlists from user's Spotify account (legacy method)
func (s *PlaylistImportService) GetUserSpotifyPlaylistsLegacy(ctx context.Context, userID uuid.UUID) ([]SpotifyPlaylistInfo, error) {
	spotifyClient, err := s.authService.GetUserSpotifyClient(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get Spotify client: %w", err)
	}

	// Get current user's playlists
	playlistPage, err := spotifyClient.CurrentUsersPlaylists(ctx, spotify.Limit(50))
	if err != nil {
		return nil, fmt.Errorf("failed to get playlists: %w", err)
	}

	var playlists []SpotifyPlaylistInfo

	// Process first page
	for _, playlist := range playlistPage.Playlists {
		info := SpotifyPlaylistInfo{
			ID:          string(playlist.ID),
			Name:        playlist.Name,
			Description: playlist.Description,
			TrackCount:  int(playlist.Tracks.Total),
			IsPublic:    playlist.IsPublic,
		}

		if len(playlist.Images) > 0 {
			info.Image = playlist.Images[0].URL
		}

		playlists = append(playlists, info)
	}

	// Handle pagination if needed
	for playlistPage.Next != "" {
		err = spotifyClient.NextPage(ctx, playlistPage)
		if err != nil {
			log.Printf("Error getting next page of playlists: %v", err)
			break
		}

		for _, playlist := range playlistPage.Playlists {
			info := SpotifyPlaylistInfo{
				ID:          string(playlist.ID),
				Name:        playlist.Name,
				Description: playlist.Description,
				TrackCount:  int(playlist.Tracks.Total),
				IsPublic:    playlist.IsPublic,
			}

			if len(playlist.Images) > 0 {
				info.Image = playlist.Images[0].URL
			}

			playlists = append(playlists, info)
		}
	}

	return playlists, nil
}

// ImportSpotifyPlaylist imports a specific Spotify playlist into Muse
func (s *PlaylistImportService) ImportSpotifyPlaylist(ctx context.Context, userID uuid.UUID, spotifyPlaylistID string) (*models.Playlist, error) {
	spotifyClient, err := s.authService.GetUserSpotifyClient(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get Spotify client: %w", err)
	}

	// Get playlist details from Spotify
	spotifyPlaylist, err := spotifyClient.GetPlaylist(ctx, spotify.ID(spotifyPlaylistID))
	if err != nil {
		return nil, fmt.Errorf("failed to get Spotify playlist: %w", err)
	}

	// Start transaction with proper context handling
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	// Ensure transaction is properly cleaned up
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err.Error() != "tx is closed" {
			log.Printf("Failed to rollback transaction: %v", err)
		}
	}()

	// Create playlist in Muse
	playlistID := uuid.New()
	var coverImage *string
	if len(spotifyPlaylist.Images) > 0 {
		coverImage = &spotifyPlaylist.Images[0].URL
	}

	insertPlaylistQuery := `
		INSERT INTO playlists (id, title, description, cover_image, creator_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, title, description, cover_image, creator_id, created_at, updated_at
	`

	var playlist models.Playlist
	row := tx.QueryRow(ctx, insertPlaylistQuery, playlistID, spotifyPlaylist.Name, spotifyPlaylist.Description, coverImage, userID)
	err = row.Scan(&playlist.ID, &playlist.Title, &playlist.Description, &playlist.CoverImage, &playlist.CreatorID, &playlist.CreatedAt, &playlist.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create playlist: %w", err)
	}

	// Get all tracks from the Spotify playlist
	tracks, err := s.getAllPlaylistTracks(ctx, spotifyClient, spotify.ID(spotifyPlaylistID))
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist tracks: %w", err)
	}

	// Import tracks and add to playlist with better error handling
	position := 1
	successfulImports := 0

	for _, track := range tracks {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("import cancelled: %w", ctx.Err())
		default:
		}

		if track.Track.Track == nil {
			log.Printf("Skipping null track at position %d", position)
			continue
		}

		// Import track to Muse database with individual error handling
		trackID, err := s.importTrackToMuse(ctx, tx, *track.Track.Track)
		if err != nil {
			log.Printf("Failed to import track %s: %v", track.Track.Track.Name, err)
			continue
		}

		// Add track to playlist
		_, err = tx.Exec(ctx, `
			INSERT INTO playlist_tracks (id, playlist_id, track_id, position, added_at)
			VALUES ($1, $2, $3, $4, NOW())
			ON CONFLICT (playlist_id, track_id) DO NOTHING
		`, uuid.New(), playlistID, trackID, position)

		if err != nil {
			log.Printf("Failed to add track to playlist: %v", err)
			continue
		}

		position++
		successfulImports++
	}

	log.Printf("Successfully imported %d out of %d tracks for playlist %s", successfulImports, len(tracks), spotifyPlaylist.Name)

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &playlist, nil
}

// Helper function to get all tracks from a playlist (handling pagination)
func (s *PlaylistImportService) getAllPlaylistTracks(ctx context.Context, client *spotify.Client, playlistID spotify.ID) ([]spotify.PlaylistItem, error) {
	var allTracks []spotify.PlaylistItem

	// Get first page
	trackPage, err := client.GetPlaylistItems(ctx, playlistID, spotify.Limit(100))
	if err != nil {
		return nil, err
	}

	allTracks = append(allTracks, trackPage.Items...)

	// Handle pagination
	for trackPage.Next != "" {
		err = client.NextPage(ctx, trackPage)
		if err != nil {
			log.Printf("Error getting next page of tracks: %v", err)
			break
		}
		allTracks = append(allTracks, trackPage.Items...)
	}

	return allTracks, nil
}

// Helper function to import a track to Muse database
func (s *PlaylistImportService) importTrackToMuse(ctx context.Context, tx pgx.Tx, spotifyTrack spotify.FullTrack) (uuid.UUID, error) {
	// First, ensure artist exists
	artistID, err := s.ensureArtistExists(ctx, tx, spotifyTrack.Artists[0])
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to ensure artist exists: %w", err)
	}

	// Then, ensure album exists
	albumID, err := s.ensureAlbumExists(ctx, tx, spotifyTrack.Album, artistID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to ensure album exists: %w", err)
	}

	// Finally, ensure track exists
	trackID, err := s.ensureTrackExists(ctx, tx, spotifyTrack, albumID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to ensure track exists: %w", err)
	}

	return trackID, nil
}

func (s *PlaylistImportService) ensureArtistExists(ctx context.Context, tx pgx.Tx, spotifyArtist spotify.SimpleArtist) (uuid.UUID, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return uuid.Nil, ctx.Err()
	default:
	}

	// Check if artist already exists
	var artistID uuid.UUID
	err := tx.QueryRow(ctx, "SELECT id FROM artists WHERE spotify_id = $1", spotifyArtist.ID).Scan(&artistID)
	if err == nil {
		return artistID, nil
	}

	// Create new artist
	artistID = uuid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO artists (id, spotify_id, name, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (spotify_id) DO UPDATE SET name = EXCLUDED.name, updated_at = NOW()
		RETURNING id
	`, artistID, spotifyArtist.ID, spotifyArtist.Name)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create artist: %w", err)
	}

	return artistID, nil
}

func (s *PlaylistImportService) ensureAlbumExists(ctx context.Context, tx pgx.Tx, spotifyAlbum spotify.SimpleAlbum, artistID uuid.UUID) (uuid.UUID, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return uuid.Nil, ctx.Err()
	default:
	}

	// Check if album already exists
	var albumID uuid.UUID
	err := tx.QueryRow(ctx, "SELECT id FROM albums WHERE spotify_id = $1", spotifyAlbum.ID).Scan(&albumID)
	if err == nil {
		return albumID, nil
	}

	// Create new album
	albumID = uuid.New()
	var coverImage *string
	if len(spotifyAlbum.Images) > 0 {
		coverImage = &spotifyAlbum.Images[0].URL
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO albums (id, spotify_id, title, artist_id, release_date, cover_image, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		ON CONFLICT (spotify_id) DO UPDATE SET 
			title = EXCLUDED.title,
			cover_image = EXCLUDED.cover_image,
			updated_at = NOW()
	`, albumID, spotifyAlbum.ID, spotifyAlbum.Name, artistID, spotifyAlbum.ReleaseDate, coverImage)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create album: %w", err)
	}

	return albumID, nil
}

func (s *PlaylistImportService) ensureTrackExists(ctx context.Context, tx pgx.Tx, spotifyTrack spotify.FullTrack, albumID uuid.UUID) (uuid.UUID, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return uuid.Nil, ctx.Err()
	default:
	}

	// Check if track already exists
	var trackID uuid.UUID
	err := tx.QueryRow(ctx, "SELECT id FROM tracks WHERE spotify_id = $1", spotifyTrack.ID).Scan(&trackID)
	if err == nil {
		return trackID, nil
	}

	// Create new track
	trackID = uuid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO tracks (id, spotify_id, title, album_id, duration_ms, track_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		ON CONFLICT (spotify_id) DO UPDATE SET 
			title = EXCLUDED.title,
			duration_ms = EXCLUDED.duration_ms,
			track_number = EXCLUDED.track_number,
			updated_at = NOW()
	`, trackID, spotifyTrack.ID, spotifyTrack.Name, albumID, spotifyTrack.Duration, spotifyTrack.TrackNumber)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create track: %w", err)
	}

	return trackID, nil
}
