package graph

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/daedal00/muse/backend/internal/models"
	"github.com/google/uuid"
	spotifyapi "github.com/zmb3/spotify/v2"
)

func isNotFoundError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "not found")
}

func (r *Resolver) getOrCreateArtist(ctx context.Context, spotifyID, name string) (*models.Artist, error) {
	existing, err := r.repos.Artist.GetBySpotifyID(ctx, spotifyID)
	if err == nil {
		return existing, nil
	}
	if !isNotFoundError(err) {
		return nil, err
	}

	now := time.Now()
	dbArtist := &models.Artist{
		ID:        uuid.New(),
		SpotifyID: &spotifyID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := r.repos.Artist.Create(ctx, dbArtist); err != nil {
		return nil, fmt.Errorf("failed to create artist: %w", err)
	}

	return dbArtist, nil
}

func (r *Resolver) importArtistFromSpotify(ctx context.Context, spotifyArtistID string) (*models.Artist, error) {
	if spotifyArtistID == "" {
		return nil, fmt.Errorf("spotify artist ID is required")
	}
	if r.spotifyServices == nil {
		return nil, fmt.Errorf("spotify service not available")
	}

	existing, err := r.repos.Artist.GetBySpotifyID(ctx, spotifyArtistID)
	if err == nil {
		return existing, nil
	}
	if !isNotFoundError(err) {
		return nil, err
	}

	artist, err := r.spotifyServices.Artist.GetArtist(ctx, spotifyapi.ID(spotifyArtistID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch artist from spotify: %w", err)
	}

	return r.getOrCreateArtist(ctx, spotifyArtistID, artist.Name)
}

func (r *Resolver) importAlbumFromSpotify(ctx context.Context, spotifyAlbumID string) (*models.Album, error) {
	if spotifyAlbumID == "" {
		return nil, fmt.Errorf("spotify album ID is required")
	}
	if r.spotifyServices == nil {
		return nil, fmt.Errorf("spotify service not available")
	}

	existing, err := r.repos.Album.GetBySpotifyID(ctx, spotifyAlbumID)
	if err == nil {
		return existing, nil
	}
	if !isNotFoundError(err) {
		return nil, err
	}

	album, err := r.spotifyServices.Album.GetAlbum(ctx, spotifyapi.ID(spotifyAlbumID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch album from spotify: %w", err)
	}
	if len(album.Artists) == 0 {
		return nil, fmt.Errorf("spotify album has no artists")
	}

	primaryArtist := album.Artists[0]
	dbArtist, err := r.getOrCreateArtist(ctx, string(primaryArtist.ID), primaryArtist.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure artist: %w", err)
	}

	var releaseDate *time.Time
	if album.ReleaseDate != "" {
		parsed := album.ReleaseDateTime()
		releaseDate = &parsed
	}

	var coverImage *string
	if len(album.Images) > 0 {
		coverImage = &album.Images[0].URL
	}

	now := time.Now()
	dbAlbum := &models.Album{
		ID:          uuid.New(),
		SpotifyID:   &spotifyAlbumID,
		Title:       album.Name,
		ArtistID:    dbArtist.ID,
		ReleaseDate: releaseDate,
		CoverImage:  coverImage,
		CreatedAt:   now,
		UpdatedAt:   now,
		Artist:      dbArtist,
	}

	if err := r.repos.Album.Create(ctx, dbAlbum); err != nil {
		return nil, fmt.Errorf("failed to create album: %w", err)
	}

	if err := r.importAlbumTracks(ctx, dbAlbum.ID, spotifyAlbumID); err != nil {
		return nil, err
	}

	return dbAlbum, nil
}

func (r *Resolver) importAlbumTracks(ctx context.Context, albumID uuid.UUID, spotifyAlbumID string) error {
	page, err := r.spotifyServices.Album.GetAlbumTracks(ctx, spotifyapi.ID(spotifyAlbumID), spotifyapi.Limit(50))
	if err != nil {
		return fmt.Errorf("failed to fetch album tracks: %w", err)
	}

	for {
		if err := r.upsertAlbumTracks(ctx, albumID, page.Tracks); err != nil {
			return err
		}

		if err := r.spotifyServices.NextTrackPage(ctx, page); err != nil {
			if err == spotifyapi.ErrNoMorePages {
				break
			}
			return fmt.Errorf("failed to fetch next track page: %w", err)
		}
	}

	return nil
}

func (r *Resolver) upsertAlbumTracks(ctx context.Context, albumID uuid.UUID, tracks []spotifyapi.SimpleTrack) error {
	for _, track := range tracks {
		spotifyID := string(track.ID)
		if spotifyID == "" {
			continue
		}

		if _, err := r.repos.Track.GetBySpotifyID(ctx, spotifyID); err == nil {
			continue
		} else if !isNotFoundError(err) {
			return fmt.Errorf("failed to check track: %w", err)
		}

		duration := int(track.Duration)
		trackNumber := int(track.TrackNumber)
		now := time.Now()
		dbTrack := &models.Track{
			ID:          uuid.New(),
			SpotifyID:   &spotifyID,
			Title:       track.Name,
			AlbumID:     albumID,
			DurationMs:  &duration,
			TrackNumber: &trackNumber,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		if err := r.repos.Track.Create(ctx, dbTrack); err != nil {
			return fmt.Errorf("failed to create track: %w", err)
		}
	}

	return nil
}

func (r *Resolver) importTrackFromSpotify(ctx context.Context, spotifyTrackID string) (*models.Track, error) {
	if spotifyTrackID == "" {
		return nil, fmt.Errorf("spotify track ID is required")
	}
	if r.spotifyServices == nil {
		return nil, fmt.Errorf("spotify service not available")
	}

	existing, err := r.repos.Track.GetBySpotifyID(ctx, spotifyTrackID)
	if err == nil {
		return existing, nil
	}
	if !isNotFoundError(err) {
		return nil, err
	}

	track, err := r.spotifyServices.Track.GetTrack(ctx, spotifyapi.ID(spotifyTrackID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch track from spotify: %w", err)
	}

	if track.Album.ID == "" {
		return nil, fmt.Errorf("spotify track missing album")
	}

	album, err := r.importAlbumFromSpotify(ctx, string(track.Album.ID))
	if err != nil {
		return nil, err
	}

	existing, err = r.repos.Track.GetBySpotifyID(ctx, spotifyTrackID)
	if err == nil {
		return existing, nil
	}
	if !isNotFoundError(err) {
		return nil, err
	}

	duration := int(track.Duration)
	trackNumber := int(track.TrackNumber)
	now := time.Now()
	dbTrack := &models.Track{
		ID:          uuid.New(),
		SpotifyID:   &spotifyTrackID,
		Title:       track.Name,
		AlbumID:     album.ID,
		DurationMs:  &duration,
		TrackNumber: &trackNumber,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := r.repos.Track.Create(ctx, dbTrack); err != nil {
		return nil, fmt.Errorf("failed to create track: %w", err)
	}

	return dbTrack, nil
}
