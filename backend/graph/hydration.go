package graph

import (
	"context"
	"fmt"

	"github.com/daedal00/muse/backend/internal/models"
	"github.com/google/uuid"
)

func (r *Resolver) hydrateAlbum(ctx context.Context, album *models.Album, artistCache map[uuid.UUID]*models.Artist) error {
	if album == nil || album.Artist != nil {
		return nil
	}

	if artistCache != nil {
		if cached, ok := artistCache[album.ArtistID]; ok {
			album.Artist = cached
			return nil
		}
	}

	artist, err := r.repos.Artist.GetByID(ctx, album.ArtistID)
	if err != nil {
		return fmt.Errorf("artist not found: %w", err)
	}

	album.Artist = artist
	if artistCache != nil {
		artistCache[album.ArtistID] = artist
	}

	return nil
}

func (r *Resolver) hydrateTrack(ctx context.Context, track *models.Track, albumCache map[uuid.UUID]*models.Album, artistCache map[uuid.UUID]*models.Artist) error {
	if track == nil {
		return nil
	}

	if track.Album == nil {
		if albumCache != nil {
			if cached, ok := albumCache[track.AlbumID]; ok {
				track.Album = cached
			}
		}

		if track.Album == nil {
			album, err := r.repos.Album.GetByID(ctx, track.AlbumID)
			if err != nil {
				return fmt.Errorf("album not found: %w", err)
			}
			track.Album = album
			if albumCache != nil {
				albumCache[track.AlbumID] = album
			}
		}
	}

	return r.hydrateAlbum(ctx, track.Album, artistCache)
}

func (r *Resolver) hydrateReview(ctx context.Context, review *models.Review, userCache map[uuid.UUID]*models.User, albumCache map[uuid.UUID]*models.Album, artistCache map[uuid.UUID]*models.Artist) error {
	if review == nil {
		return nil
	}

	if review.User == nil {
		if userCache != nil {
			if cached, ok := userCache[review.UserID]; ok {
				review.User = cached
			}
		}

		if review.User == nil {
			user, err := r.repos.User.GetByID(ctx, review.UserID)
			if err != nil {
				return fmt.Errorf("user not found: %w", err)
			}
			review.User = user
			if userCache != nil {
				userCache[review.UserID] = user
			}
		}
	}

	if review.Album == nil {
		if albumCache != nil {
			if cached, ok := albumCache[review.AlbumID]; ok {
				review.Album = cached
			}
		}

		if review.Album == nil {
			album, err := r.repos.Album.GetByID(ctx, review.AlbumID)
			if err != nil {
				return fmt.Errorf("album not found: %w", err)
			}
			review.Album = album
			if albumCache != nil {
				albumCache[review.AlbumID] = album
			}
		}
	}

	return r.hydrateAlbum(ctx, review.Album, artistCache)
}

func (r *Resolver) hydratePlaylist(ctx context.Context, playlist *models.Playlist, userCache map[uuid.UUID]*models.User) error {
	if playlist == nil || playlist.Creator != nil {
		return nil
	}

	if userCache != nil {
		if cached, ok := userCache[playlist.CreatorID]; ok {
			playlist.Creator = cached
			return nil
		}
	}

	user, err := r.repos.User.GetByID(ctx, playlist.CreatorID)
	if err != nil {
		return fmt.Errorf("creator not found: %w", err)
	}

	playlist.Creator = user
	if userCache != nil {
		userCache[playlist.CreatorID] = user
	}

	return nil
}
