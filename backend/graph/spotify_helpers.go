package graph

import (
	"strings"
	"time"

	"github.com/daedal00/muse/backend/graph/model"
	spotifyapi "github.com/zmb3/spotify/v2"
)

func spotifyArtistsToSearchResults(artists []spotifyapi.SimpleArtist) []*model.ArtistSearchResult {
	results := make([]*model.ArtistSearchResult, 0, len(artists))
	for _, artist := range artists {
		results = append(results, &model.ArtistSearchResult{
			ID:             string(artist.ID),
			Name:           artist.Name,
			ExternalSource: model.ExternalSourceSpotify,
		})
	}
	return results
}

func spotifyAlbumToSearchResult(album spotifyapi.SimpleAlbum) *model.AlbumSearchResult {
	var releaseDate *string
	if album.ReleaseDate != "" {
		releaseDate = &album.ReleaseDate
	}

	var coverImage *string
	if len(album.Images) > 0 {
		coverImage = &album.Images[0].URL
	}

	return &model.AlbumSearchResult{
		ID:             string(album.ID),
		Title:          album.Name,
		Artist:         spotifyArtistsToSearchResults(album.Artists),
		ReleaseDate:    releaseDate,
		CoverImage:     coverImage,
		ExternalSource: model.ExternalSourceSpotify,
	}
}

func spotifyFullAlbumToSearchResult(album *spotifyapi.FullAlbum) *model.AlbumSearchResult {
	if album == nil {
		return nil
	}

	var releaseDate *string
	if album.ReleaseDate != "" {
		releaseDate = &album.ReleaseDate
	}

	var coverImage *string
	if len(album.Images) > 0 {
		coverImage = &album.Images[0].URL
	}

	return &model.AlbumSearchResult{
		ID:             string(album.ID),
		Title:          album.Name,
		Artist:         spotifyArtistsToSearchResults(album.Artists),
		ReleaseDate:    releaseDate,
		CoverImage:     coverImage,
		ExternalSource: model.ExternalSourceSpotify,
	}
}

func spotifyFullTrackToSearchResult(track spotifyapi.FullTrack) *model.TrackSearchResult {
	var durationSeconds *int32
	if track.Duration > 0 {
		val := safeIntToInt32(int(track.Duration) / 1000)
		durationSeconds = &val
	}

	var trackNumber *int32
	if track.TrackNumber > 0 {
		val := safeIntToInt32(int(track.TrackNumber))
		trackNumber = &val
	}

	var albumResult *model.AlbumSearchResult
	if track.Album.ID != "" {
		albumResult = spotifyAlbumToSearchResult(track.Album)
	}

	return &model.TrackSearchResult{
		ID:             string(track.ID),
		Title:          track.Name,
		Duration:       durationSeconds,
		TrackNumber:    trackNumber,
		Album:          albumResult,
		Artists:        spotifyArtistsToSearchResults(track.Artists),
		ExternalSource: model.ExternalSourceSpotify,
	}
}

func spotifySimpleTrackToSearchResult(track spotifyapi.SimpleTrack, album *spotifyapi.SimpleAlbum) *model.TrackSearchResult {
	var durationSeconds *int32
	if track.Duration > 0 {
		val := safeIntToInt32(int(track.Duration) / 1000)
		durationSeconds = &val
	}

	var trackNumber *int32
	if track.TrackNumber > 0 {
		val := safeIntToInt32(int(track.TrackNumber))
		trackNumber = &val
	}

	var albumResult *model.AlbumSearchResult
	if album != nil && album.ID != "" {
		albumResult = spotifyAlbumToSearchResult(*album)
	}

	return &model.TrackSearchResult{
		ID:             string(track.ID),
		Title:          track.Name,
		Duration:       durationSeconds,
		TrackNumber:    trackNumber,
		Album:          albumResult,
		Artists:        spotifyArtistsToSearchResults(track.Artists),
		ExternalSource: model.ExternalSourceSpotify,
	}
}

func spotifyTimeRangeToRequest(rangeInput string) spotifyapi.Range {
	switch strings.ToUpper(rangeInput) {
	case "SHORT_TERM":
		return spotifyapi.ShortTermRange
	case "LONG_TERM":
		return spotifyapi.LongTermRange
	default:
		return spotifyapi.MediumTermRange
	}
}

func optionalDateTime(value time.Time) *string {
	if value.IsZero() {
		return nil
	}
	formatted := value.Format(time.RFC3339)
	return &formatted
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
