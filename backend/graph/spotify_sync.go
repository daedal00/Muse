package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/daedal00/muse/backend/graph/model"
	"github.com/daedal00/muse/backend/internal/models"
	"github.com/daedal00/muse/backend/internal/spotify"
	"github.com/google/uuid"
	spotifyapi "github.com/zmb3/spotify/v2"
)

func (r *Resolver) importSpotifyTrack(ctx context.Context, spotifyTrackID, spotifyAlbumID string) (*models.Track, bool, bool, error) {
	if spotifyTrackID == "" {
		return nil, false, false, fmt.Errorf("spotify track ID is required")
	}

	if existing, err := r.repos.Track.GetBySpotifyID(ctx, spotifyTrackID); err == nil {
		return existing, false, false, nil
	} else if !isNotFoundError(err) {
		return nil, false, false, err
	}

	albumCreated := false
	if spotifyAlbumID != "" {
		if _, err := r.repos.Album.GetBySpotifyID(ctx, spotifyAlbumID); err != nil {
			if isNotFoundError(err) {
				albumCreated = true
			} else {
				return nil, false, false, err
			}
		}
	}

	track, err := r.importTrackFromSpotify(ctx, spotifyTrackID)
	if err != nil {
		return nil, false, false, err
	}

	return track, true, albumCreated, nil
}

func (r *Resolver) importSpotifyPlaylist(
	ctx context.Context,
	services *spotify.Services,
	spotifyPlaylistID string,
	userID uuid.UUID,
) (*models.Playlist, *model.ImportSummary, error) {
	if spotifyPlaylistID == "" {
		return nil, nil, fmt.Errorf("spotify playlist ID is required")
	}

	if existing, err := r.repos.Playlist.GetBySpotifyID(ctx, spotifyPlaylistID); err == nil {
		return existing, &model.ImportSummary{}, nil
	} else if !isNotFoundError(err) {
		return nil, nil, err
	}

	playlist, err := services.Playlist.GetPlaylist(ctx, spotifyapi.ID(spotifyPlaylistID))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch spotify playlist: %w", err)
	}

	var description *string
	if playlist.Description != "" {
		desc := playlist.Description
		description = &desc
	}

	var coverImage *string
	if len(playlist.Images) > 0 {
		coverImage = &playlist.Images[0].URL
	}

	now := time.Now()
	dbPlaylist := &models.Playlist{
		ID:          uuid.New(),
		SpotifyID:   &spotifyPlaylistID,
		Title:       playlist.Name,
		Description: description,
		CoverImage:  coverImage,
		CreatorID:   userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := r.repos.Playlist.Create(ctx, dbPlaylist); err != nil {
		return nil, nil, fmt.Errorf("failed to create playlist: %w", err)
	}

	summary := &model.ImportSummary{
		ImportedPlaylists: 1,
	}

	items, err := services.Playlist.GetPlaylistItems(ctx, spotifyapi.ID(spotifyPlaylistID), spotifyapi.Limit(50))
	if err != nil {
		return dbPlaylist, summary, fmt.Errorf("failed to fetch playlist items: %w", err)
	}

	for {
		for _, item := range items.Items {
			if item.Track.Track == nil {
				continue
			}

			track := item.Track.Track
			trackID := string(track.ID)
			albumID := ""
			if track.Album.ID != "" {
				albumID = string(track.Album.ID)
			}

			importedTrack, trackCreated, albumCreated, err := r.importSpotifyTrack(ctx, trackID, albumID)
			if err != nil {
				return dbPlaylist, summary, err
			}

			if trackCreated {
				summary.ImportedTracks++
			}
			if albumCreated {
				summary.ImportedAlbums++
			}

			if importedTrack != nil {
				if err := r.repos.Playlist.AddTrack(ctx, dbPlaylist.ID, importedTrack.ID, 0); err != nil {
					return dbPlaylist, summary, fmt.Errorf("failed to add track to playlist: %w", err)
				}
			}
		}

		if err := services.NextPlaylistPage(ctx, items); err != nil {
			if err == spotifyapi.ErrNoMorePages {
				break
			}
			return dbPlaylist, summary, fmt.Errorf("failed to fetch next playlist page: %w", err)
		}
	}

	return dbPlaylist, summary, nil
}
