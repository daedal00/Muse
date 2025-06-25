package graph

import (
	"time"

	"github.com/daedal00/muse/backend/graph/model"
	"github.com/daedal00/muse/backend/internal/models"
)

// Helper functions to convert between database models and GraphQL models

func dbUserToGraphQL(dbUser *models.User) *model.User {
	if dbUser == nil {
		return nil
	}
	
	return &model.User{
		ID:     dbUser.ID.String(),
		Name:   dbUser.Name,
		Email:  dbUser.Email,
		Bio:    dbUser.Bio,
		Avatar: dbUser.Avatar,
	}
}

func dbArtistToGraphQL(dbArtist *models.Artist) *model.Artist {
	if dbArtist == nil {
		return nil
	}
	
	return &model.Artist{
		ID:        dbArtist.ID.String(),
		SpotifyID: dbArtist.SpotifyID,
		Name:      dbArtist.Name,
	}
}

func dbAlbumToGraphQL(dbAlbum *models.Album) *model.Album {
	if dbAlbum == nil {
		return nil
	}
	
	var releaseDate *string
	if dbAlbum.ReleaseDate != nil {
		formatted := dbAlbum.ReleaseDate.Format(time.RFC3339)
		releaseDate = &formatted
	}
	
	return &model.Album{
		ID:          dbAlbum.ID.String(),
		SpotifyID:   dbAlbum.SpotifyID,
		Title:       dbAlbum.Title,
		ReleaseDate: releaseDate,
		CoverImage:  dbAlbum.CoverImage,
		Artist:      dbArtistToGraphQL(dbAlbum.Artist),
	}
}

func dbTrackToGraphQL(dbTrack *models.Track) *model.Track {
	if dbTrack == nil {
		return nil
	}
	
	var duration *int32
	if dbTrack.DurationMs != nil {
		val := int32(*dbTrack.DurationMs)
		duration = &val
	}
	
	var trackNumber *int32
	if dbTrack.TrackNumber != nil {
		val := int32(*dbTrack.TrackNumber)
		trackNumber = &val
	}
	
	return &model.Track{
		ID:          dbTrack.ID.String(),
		SpotifyID:   dbTrack.SpotifyID,
		Title:       dbTrack.Title,
		Duration:    duration,
		TrackNumber: trackNumber,
		Album:       dbAlbumToGraphQL(dbTrack.Album),
	}
}

func dbReviewToGraphQL(dbReview *models.Review) *model.Review {
	if dbReview == nil {
		return nil
	}
	
	return &model.Review{
		ID:         dbReview.ID.String(),
		User:       dbUserToGraphQL(dbReview.User),
		Album:      dbAlbumToGraphQL(dbReview.Album),
		Rating:     int32(dbReview.Rating),
		ReviewText: dbReview.ReviewText,
		CreatedAt:  dbReview.CreatedAt.Format(time.RFC3339),
	}
}

func dbPlaylistToGraphQL(dbPlaylist *models.Playlist) *model.Playlist {
	if dbPlaylist == nil {
		return nil
	}
	
	return &model.Playlist{
		ID:          dbPlaylist.ID.String(),
		Title:       dbPlaylist.Title,
		Description: dbPlaylist.Description,
		CoverImage:  dbPlaylist.CoverImage,
		Creator:     dbUserToGraphQL(dbPlaylist.Creator),
		CreatedAt:   dbPlaylist.CreatedAt.Format(time.RFC3339),
	}
}

// Connection builders
func buildAlbumConnection(albums []*models.Album, totalCount int, hasNextPage bool) *model.AlbumConnection {
	edges := make([]*model.AlbumEdge, len(albums))
	for i, album := range albums {
		edges[i] = &model.AlbumEdge{
			Cursor: album.ID.String(),
			Node:   dbAlbumToGraphQL(album),
		}
	}
	
	var endCursor *string
	if len(edges) > 0 {
		endCursor = &edges[len(edges)-1].Cursor
	}
	
	return &model.AlbumConnection{
		TotalCount: int32(totalCount),
		Edges:      edges,
		PageInfo: &model.PageInfo{
			EndCursor:   endCursor,
			HasNextPage: hasNextPage,
		},
	}
}

func buildTrackConnection(tracks []*models.Track, totalCount int, hasNextPage bool) *model.TrackConnection {
	edges := make([]*model.TrackEdge, len(tracks))
	for i, track := range tracks {
		edges[i] = &model.TrackEdge{
			Cursor: track.ID.String(),
			Node:   dbTrackToGraphQL(track),
		}
	}
	
	var endCursor *string
	if len(edges) > 0 {
		endCursor = &edges[len(edges)-1].Cursor
	}
	
	return &model.TrackConnection{
		TotalCount: int32(totalCount),
		Edges:      edges,
		PageInfo: &model.PageInfo{
			EndCursor:   endCursor,
			HasNextPage: hasNextPage,
		},
	}
}

func buildReviewConnection(reviews []*models.Review, totalCount int, hasNextPage bool) *model.ReviewConnection {
	edges := make([]*model.ReviewEdge, len(reviews))
	for i, review := range reviews {
		edges[i] = &model.ReviewEdge{
			Cursor: review.ID.String(),
			Node:   dbReviewToGraphQL(review),
		}
	}
	
	var endCursor *string
	if len(edges) > 0 {
		endCursor = &edges[len(edges)-1].Cursor
	}
	
	return &model.ReviewConnection{
		TotalCount: int32(totalCount),
		Edges:      edges,
		PageInfo: &model.PageInfo{
			EndCursor:   endCursor,
			HasNextPage: hasNextPage,
		},
	}
}

func buildPlaylistConnection(playlists []*models.Playlist, totalCount int, hasNextPage bool) *model.PlaylistConnection {
	edges := make([]*model.PlaylistEdge, len(playlists))
	for i, playlist := range playlists {
		edges[i] = &model.PlaylistEdge{
			Cursor: playlist.ID.String(),
			Node:   dbPlaylistToGraphQL(playlist),
		}
	}
	
	var endCursor *string
	if len(edges) > 0 {
		endCursor = &edges[len(edges)-1].Cursor
	}
	
	return &model.PlaylistConnection{
		TotalCount: int32(totalCount),
		Edges:      edges,
		PageInfo: &model.PageInfo{
			EndCursor:   endCursor,
			HasNextPage: hasNextPage,
		},
	}
} 