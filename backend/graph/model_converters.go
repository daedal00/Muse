package graph

import (
	"math"
	"time"

	"github.com/daedal00/muse/backend/graph/model"
	"github.com/daedal00/muse/backend/internal/models"
)

// Helper functions to convert between database models and GraphQL models

// safeIntToInt32 safely converts int to int32, clamping to valid range
func safeIntToInt32(val int) int32 {
	if val > math.MaxInt32 {
		return math.MaxInt32
	}
	if val < math.MinInt32 {
		return math.MinInt32
	}
	return int32(val)
}

// safeLenToInt32 safely converts slice length to int32
func safeLenToInt32(length int) int32 {
	return safeIntToInt32(length)
}

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
		val := safeIntToInt32(*dbTrack.DurationMs / 1000)
		duration = &val
	}

	var trackNumber *int32
	if dbTrack.TrackNumber != nil {
		val := safeIntToInt32(*dbTrack.TrackNumber)
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
		Rating:     safeIntToInt32(dbReview.Rating),
		ReviewText: dbReview.ReviewText,
		CreatedAt:  dbReview.CreatedAt.Format(time.RFC3339),
	}
}

func dbTrackReviewToGraphQL(dbReview *models.TrackReview) *model.TrackReview {
	if dbReview == nil {
		return nil
	}

	return &model.TrackReview{
		ID:         dbReview.ID.String(),
		User:       dbUserToGraphQL(dbReview.User),
		Track:      dbTrackToGraphQL(dbReview.Track),
		Rating:     safeIntToInt32(dbReview.Rating),
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

func dbProfileSettingsToGraphQL(dbSettings *models.ProfileSettings) *model.ProfileSettings {
	if dbSettings == nil {
		return nil
	}

	// Convert layout
	layout := model.ProfileLayoutGrid
	switch dbSettings.Layout {
	case models.ProfileLayoutList:
		layout = model.ProfileLayoutList
	case models.ProfileLayoutBento:
		layout = model.ProfileLayoutBento
	}

	// Convert background style
	var bgStyle *model.BackgroundStyle
	if dbSettings.BackgroundStyle != nil {
		switch *dbSettings.BackgroundStyle {
		case models.BackgroundStyleSolid:
			bs := model.BackgroundStyleSolid
			bgStyle = &bs
		case models.BackgroundStyleGradient:
			bs := model.BackgroundStyleGradient
			bgStyle = &bs
		case models.BackgroundStyleImage:
			bs := model.BackgroundStyleImage
			bgStyle = &bs
		}
	}

	// Convert bio style
	var bioStyle *model.BioStyle
	if dbSettings.BioStyle != nil {
		switch *dbSettings.BioStyle {
		case models.BioStyleMinimal:
			bs := model.BioStyleMinimal
			bioStyle = &bs
		case models.BioStyleDetailed:
			bs := model.BioStyleDetailed
			bioStyle = &bs
		case models.BioStyleQuote:
			bs := model.BioStyleQuote
			bioStyle = &bs
		}
	}

	// Convert UUID arrays to string arrays
	pinnedAlbumIds := make([]string, len(dbSettings.PinnedAlbumIDs))
	for i, id := range dbSettings.PinnedAlbumIDs {
		pinnedAlbumIds[i] = id.String()
	}

	pinnedTrackIds := make([]string, len(dbSettings.PinnedTrackIDs))
	for i, id := range dbSettings.PinnedTrackIDs {
		pinnedTrackIds[i] = id.String()
	}

	featuredArtistIds := make([]string, len(dbSettings.FeaturedArtistIDs))
	for i, id := range dbSettings.FeaturedArtistIDs {
		featuredArtistIds[i] = id.String()
	}

	return &model.ProfileSettings{
		Layout:               layout,
		PrimaryColor:         dbSettings.PrimaryColor,
		AccentColor:          dbSettings.AccentColor,
		BackgroundStyle:      bgStyle,
		BackgroundValue:      dbSettings.BackgroundValue,
		PinnedAlbumIds:       pinnedAlbumIds,
		PinnedTrackIds:       pinnedTrackIds,
		FeaturedArtistIds:    featuredArtistIds,
		SectionsOrder:        dbSettings.SectionsOrder,
		ShowSpotifyStats:     dbSettings.ShowSpotifyStats,
		ShowListeningHistory: dbSettings.ShowListeningHistory,
		BioStyle:             bioStyle,
		CustomTags:           dbSettings.CustomTags,
	}
}
