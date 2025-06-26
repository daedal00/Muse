package graph

import (
	"math"
	"time"

	"github.com/daedal00/muse/backend/graph/model"
	"github.com/daedal00/muse/backend/internal/models"
	"github.com/zmb3/spotify/v2"
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

// Convert Spotify models to GraphQL models
func spotifyArtistToGraphQL(spotifyArtist *models.SpotifyArtist) *model.Artist {
	if spotifyArtist == nil {
		return nil
	}

	return &model.Artist{
		ID:        spotifyArtist.ID, // Use Spotify ID directly
		SpotifyID: &spotifyArtist.ID,
		Name:      spotifyArtist.Name,
	}
}

func spotifyAlbumToGraphQL(spotifyAlbum *models.SpotifyAlbum) *model.Album {
	if spotifyAlbum == nil {
		return nil
	}

	var releaseDate *time.Time
	if spotifyAlbum.ReleaseDate != "" {
		// Parse Spotify date format (YYYY-MM-DD or YYYY)
		if parsed, err := time.Parse("2006-01-02", spotifyAlbum.ReleaseDate); err == nil {
			releaseDate = &parsed
		} else if parsed, err := time.Parse("2006", spotifyAlbum.ReleaseDate); err == nil {
			releaseDate = &parsed
		}
	}

	var formattedReleaseDate *string
	if releaseDate != nil {
		formatted := releaseDate.Format(time.RFC3339)
		formattedReleaseDate = &formatted
	}

	// Convert artists
	var artists []*model.Artist
	for _, artist := range spotifyAlbum.Artists {
		artists = append(artists, spotifyArtistToGraphQL(&artist))
	}

	var primaryArtist *model.Artist
	if len(artists) > 0 {
		primaryArtist = artists[0]
	}

	// Get first image if available
	var coverImage *string
	if len(spotifyAlbum.Images) > 0 {
		coverImage = &spotifyAlbum.Images[0]
	}

	return &model.Album{
		ID:          spotifyAlbum.ID, // Use Spotify ID directly
		SpotifyID:   &spotifyAlbum.ID,
		Title:       spotifyAlbum.Name,
		ReleaseDate: formattedReleaseDate,
		CoverImage:  coverImage,
		Artist:      primaryArtist,
	}
}

func spotifyTrackToGraphQL(spotifyTrack *models.SpotifyTrack) *model.Track {
	if spotifyTrack == nil {
		return nil
	}

	var duration *int32
	if spotifyTrack.DurationMS > 0 {
		// Convert from milliseconds to seconds
		durationSec := spotifyTrack.DurationMS / 1000
		val := safeIntToInt32(durationSec)
		duration = &val
	}

	var trackNumber *int32
	if spotifyTrack.TrackNumber > 0 {
		val := safeIntToInt32(spotifyTrack.TrackNumber)
		trackNumber = &val
	}

	return &model.Track{
		ID:          spotifyTrack.ID, // Use Spotify ID directly
		SpotifyID:   &spotifyTrack.ID,
		Title:       spotifyTrack.Name,
		Duration:    duration,
		TrackNumber: trackNumber,
		// Album will be populated by field resolver if needed
	}
}

func dbReviewToGraphQL(dbReview *models.Review) *model.Review {
	if dbReview == nil {
		return nil
	}

	// Note: In the new architecture, reviews reference Spotify IDs directly
	// The Album field will need to be resolved by fetching from Spotify API/cache
	// For now, we'll set it to nil and handle it in field resolvers
	return &model.Review{
		ID:         dbReview.ID.String(),
		User:       dbUserToGraphQL(dbReview.User),
		Album:      nil, // Will be resolved by field resolver using SpotifyID
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
		CreatedAt:   dbPlaylist.CreatedAt.Format(time.RFC3339),
	}
}

// Helper functions to convert Spotify API data to internal models
func spotifyAPITrackToModel(track *spotify.FullTrack) *models.SpotifyTrack {
	if track == nil {
		return nil
	}

	// Convert artists
	var artists []models.SpotifyArtist
	for _, artist := range track.Artists {
		artists = append(artists, models.SpotifyArtist{
			ID:   string(artist.ID),
			Name: artist.Name,
		})
	}

	// Convert album with basic info
	albumModel := models.SpotifyAlbum{
		ID:          string(track.Album.ID),
		Name:        track.Album.Name,
		ReleaseDate: track.Album.ReleaseDate,
	}

	// Convert album artists
	for _, artist := range track.Album.Artists {
		albumModel.Artists = append(albumModel.Artists, models.SpotifyArtist{
			ID:   string(artist.ID),
			Name: artist.Name,
		})
	}

	// Convert album images
	for _, image := range track.Album.Images {
		albumModel.Images = append(albumModel.Images, image.URL)
	}

	return &models.SpotifyTrack{
		ID:          string(track.ID),
		Name:        track.Name,
		Artists:     artists,
		Album:       albumModel,
		DurationMS:  int(track.Duration),
		TrackNumber: int(track.TrackNumber),
		PreviewURL:  &track.PreviewURL,
	}
}

func spotifyAPIAlbumToModel(album *spotify.FullAlbum) *models.SpotifyAlbum {
	if album == nil {
		return nil
	}

	// Convert artists
	var artists []models.SpotifyArtist
	for _, artist := range album.Artists {
		artists = append(artists, models.SpotifyArtist{
			ID:   string(artist.ID),
			Name: artist.Name,
		})
	}

	// Convert images
	var images []string
	for _, image := range album.Images {
		images = append(images, image.URL)
	}

	return &models.SpotifyAlbum{
		ID:          string(album.ID),
		Name:        album.Name,
		Artists:     artists,
		ReleaseDate: album.ReleaseDate,
		Images:      images,
	}
}

func spotifyAPIArtistToModel(artist *spotify.FullArtist) *models.SpotifyArtist {
	if artist == nil {
		return nil
	}

	// Convert images
	var images []string
	for _, image := range artist.Images {
		images = append(images, image.URL)
	}

	// Convert genres - simplified loop
	genres := append([]string(nil), artist.Genres...)

	return &models.SpotifyArtist{
		ID:        string(artist.ID),
		Name:      artist.Name,
		Images:    images,
		Genres:    genres,
		Followers: int(artist.Followers.Count),
	}
}

// Helper function to convert Spotify API AlbumDetails to GraphQL AlbumDetails
func spotifyAPIAlbumToAlbumDetails(album *spotify.FullAlbum, tracks []spotify.SimpleTrack) *model.AlbumDetails {
	if album == nil {
		return nil
	}

	// Convert main album info
	albumModel := spotifyAPIAlbumToModel(album)
	albumGraphQL := spotifyAlbumToGraphQL(albumModel)

	// Convert tracks to TrackDetails
	var trackDetails []*model.TrackDetails
	for _, track := range tracks {
		// Convert track artists
		var featuredArtists []*model.Artist
		for _, artist := range track.Artists {
			featuredArtists = append(featuredArtists, &model.Artist{
				ID:        string(artist.ID),
				SpotifyID: func() *string { s := string(artist.ID); return &s }(),
				Name:      artist.Name,
			})
		}

		var duration *int32
		if track.Duration > 0 {
			durationSec := int(track.Duration) / 1000
			val := safeIntToInt32(durationSec)
			duration = &val
		}

		var trackNumber *int32
		if track.TrackNumber > 0 {
			val := safeIntToInt32(int(track.TrackNumber))
			trackNumber = &val
		}

		trackDetails = append(trackDetails, &model.TrackDetails{
			ID:              string(track.ID),
			SpotifyID:       func() *string { s := string(track.ID); return &s }(),
			Title:           track.Name,
			Duration:        duration,
			TrackNumber:     trackNumber,
			FeaturedArtists: featuredArtists,
			AverageRating:   nil, // TODO: Calculate from reviews
			TotalReviews:    0,   // TODO: Count from reviews
		})
	}

	return &model.AlbumDetails{
		ID:            albumGraphQL.ID,
		SpotifyID:     albumGraphQL.SpotifyID,
		Title:         albumGraphQL.Title,
		Artist:        albumGraphQL.Artist,
		ReleaseDate:   albumGraphQL.ReleaseDate,
		CoverImage:    albumGraphQL.CoverImage,
		Tracks:        trackDetails,
		AverageRating: nil, // TODO: Calculate from reviews
		TotalReviews:  0,   // TODO: Count from reviews
	}
}

// Helper function to convert Spotify API ArtistDetails to GraphQL ArtistDetails
func spotifyAPIArtistToArtistDetails(artist *spotify.FullArtist, albums []spotify.SimpleAlbum, topTracks []spotify.FullTrack) *model.ArtistDetails {
	if artist == nil {
		return nil
	}

	// Convert basic artist info
	artistModel := spotifyAPIArtistToModel(artist)
	artistGraphQL := spotifyArtistToGraphQL(artistModel)

	// Convert albums
	var albumList []*model.Album
	for _, album := range albums {
		// Convert album artists
		var albumArtists []models.SpotifyArtist
		for _, artist := range album.Artists {
			albumArtists = append(albumArtists, models.SpotifyArtist{
				ID:   string(artist.ID),
				Name: artist.Name,
			})
		}

		// Convert album images
		var images []string
		for _, image := range album.Images {
			images = append(images, image.URL)
		}

		albumModel := &models.SpotifyAlbum{
			ID:          string(album.ID),
			Name:        album.Name,
			Artists:     albumArtists,
			ReleaseDate: album.ReleaseDate,
			Images:      images,
		}

		albumList = append(albumList, spotifyAlbumToGraphQL(albumModel))
	}

	// Convert top tracks
	var trackList []*model.Track
	for _, track := range topTracks {
		trackModel := spotifyAPITrackToModel(&track)
		trackList = append(trackList, spotifyTrackToGraphQL(trackModel))
	}

	// Get first image if available
	var image *string
	if len(artistModel.Images) > 0 {
		image = &artistModel.Images[0]
	}

	return &model.ArtistDetails{
		ID:        artistGraphQL.ID,
		SpotifyID: artistGraphQL.SpotifyID,
		Name:      artistGraphQL.Name,
		Albums:    albumList,
		TopTracks: trackList,
		Image:     image,
		Followers: func() *int32 { f := safeIntToInt32(artistModel.Followers); return &f }(),
		Genres:    artistModel.Genres,
	}
}
