package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_Creation(t *testing.T) {
	user := &User{
		ID:           uuid.New(),
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		Bio:          stringPtr("Test bio"),
		Avatar:       stringPtr("https://example.com/avatar.jpg"),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	assert.NotEqual(t, uuid.Nil, user.ID)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NotNil(t, user.Bio)
	assert.Equal(t, "Test bio", *user.Bio)
}

func TestUser_JSONMarshaling(t *testing.T) {
	user := &User{
		ID:        uuid.New(),
		Name:      "Test User",
		Email:     "test@example.com",
		Bio:       stringPtr("Test bio"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test marshaling
	data, err := json.Marshal(user)
	require.NoError(t, err)
	assert.Contains(t, string(data), "Test User")
	assert.NotContains(t, string(data), "password_hash") // Should be omitted

	// Test unmarshaling
	var unmarshaled User
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, user.Name, unmarshaled.Name)
	assert.Equal(t, user.Email, unmarshaled.Email)
}

func TestUser_SpotifyFields(t *testing.T) {
	user := &User{
		ID:                  uuid.New(),
		Name:                "Test User",
		Email:               "test@example.com",
		SpotifyID:           stringPtr("spotify_user_123"),
		SpotifyAccessToken:  stringPtr("access_token"),
		SpotifyRefreshToken: stringPtr("refresh_token"),
		SpotifyTokenExpiry:  timePtr(time.Now().Add(time.Hour)),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	assert.NotNil(t, user.SpotifyID)
	assert.Equal(t, "spotify_user_123", *user.SpotifyID)
	assert.NotNil(t, user.SpotifyAccessToken)
	assert.NotNil(t, user.SpotifyRefreshToken)
	assert.NotNil(t, user.SpotifyTokenExpiry)
}

func TestReview_Creation(t *testing.T) {
	userID := uuid.New()
	review := &Review{
		ID:          uuid.New(),
		UserID:      userID,
		SpotifyID:   "spotify_track_123",
		SpotifyType: "track",
		Rating:      4,
		ReviewText:  stringPtr("Great track!"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	assert.Equal(t, userID, review.UserID)
	assert.Equal(t, "spotify_track_123", review.SpotifyID)
	assert.Equal(t, "track", review.SpotifyType)
	assert.Equal(t, 4, review.Rating)
	assert.Equal(t, "Great track!", *review.ReviewText)
}

func TestReview_ValidSpotifyTypes(t *testing.T) {
	validTypes := []string{"album", "track"}

	for _, spotifyType := range validTypes {
		review := &Review{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			SpotifyID:   "spotify_" + spotifyType + "_123",
			SpotifyType: spotifyType,
			Rating:      5,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		assert.Equal(t, spotifyType, review.SpotifyType)
	}
}

func TestPlaylist_Creation(t *testing.T) {
	creatorID := uuid.New()
	playlist := &Playlist{
		ID:          uuid.New(),
		Title:       "My Awesome Playlist",
		Description: stringPtr("A great collection of tracks"),
		CoverImage:  stringPtr("https://example.com/cover.jpg"),
		CreatorID:   creatorID,
		IsPublic:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	assert.Equal(t, "My Awesome Playlist", playlist.Title)
	assert.Equal(t, creatorID, playlist.CreatorID)
	assert.True(t, playlist.IsPublic)
	assert.NotNil(t, playlist.Description)
}

func TestPlaylistTrack_Creation(t *testing.T) {
	playlistID := uuid.New()
	userID := uuid.New()

	track := &PlaylistTrack{
		ID:            uuid.New(),
		PlaylistID:    playlistID,
		SpotifyID:     "spotify_track_456",
		Position:      1,
		AddedAt:       time.Now(),
		AddedByUserID: userID,
	}

	assert.Equal(t, playlistID, track.PlaylistID)
	assert.Equal(t, "spotify_track_456", track.SpotifyID)
	assert.Equal(t, 1, track.Position)
	assert.Equal(t, userID, track.AddedByUserID)
}

func TestUserPreferences_Creation(t *testing.T) {
	userID := uuid.New()
	prefs := &UserPreferences{
		ID:                uuid.New(),
		UserID:            userID,
		PreferredGenres:   []string{"rock", "pop", "electronic"},
		FavoriteArtistIDs: []string{"spotify_artist_1", "spotify_artist_2"},
		NotificationSettings: map[string]bool{
			"email_notifications": true,
			"push_notifications":  false,
		},
		PrivacySettings: map[string]bool{
			"public_profile":   true,
			"public_playlists": false,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assert.Equal(t, userID, prefs.UserID)
	assert.Len(t, prefs.PreferredGenres, 3)
	assert.Contains(t, prefs.PreferredGenres, "rock")
	assert.Len(t, prefs.FavoriteArtistIDs, 2)
	assert.True(t, prefs.NotificationSettings["email_notifications"])
	assert.False(t, prefs.NotificationSettings["push_notifications"])
}

func TestSpotifyArtist_Creation(t *testing.T) {
	artist := &SpotifyArtist{
		ID:        "spotify_artist_123",
		Name:      "Test Artist",
		Images:    []string{"https://example.com/image1.jpg", "https://example.com/image2.jpg"},
		Genres:    []string{"rock", "alternative"},
		Followers: 50000,
	}

	assert.Equal(t, "spotify_artist_123", artist.ID)
	assert.Equal(t, "Test Artist", artist.Name)
	assert.Len(t, artist.Images, 2)
	assert.Len(t, artist.Genres, 2)
	assert.Equal(t, 50000, artist.Followers)
}

func TestSpotifyAlbum_Creation(t *testing.T) {
	album := &SpotifyAlbum{
		ID:          "spotify_album_123",
		Name:        "Test Album",
		Artists:     []SpotifyArtist{{ID: "artist_1", Name: "Artist 1"}},
		Images:      []string{"https://example.com/album.jpg"},
		ReleaseDate: "2023-01-01",
		TotalTracks: 12,
		TrackIds:    []string{"track_1", "track_2", "track_3"},
	}

	assert.Equal(t, "spotify_album_123", album.ID)
	assert.Equal(t, "Test Album", album.Name)
	assert.Len(t, album.Artists, 1)
	assert.Equal(t, 12, album.TotalTracks)
	assert.Len(t, album.TrackIds, 3)
}

func TestSpotifyTrack_Creation(t *testing.T) {
	track := &SpotifyTrack{
		ID:   "spotify_track_123",
		Name: "Test Track",
		Artists: []SpotifyArtist{
			{ID: "artist_1", Name: "Artist 1"},
		},
		Album: SpotifyAlbum{
			ID:   "album_1",
			Name: "Album 1",
		},
		DurationMS:  240000, // 4 minutes
		TrackNumber: 5,
		PreviewURL:  stringPtr("https://example.com/preview.mp3"),
		ExternalURLs: map[string]string{
			"spotify": "https://open.spotify.com/track/123",
		},
	}

	assert.Equal(t, "spotify_track_123", track.ID)
	assert.Equal(t, "Test Track", track.Name)
	assert.Len(t, track.Artists, 1)
	assert.Equal(t, 240000, track.DurationMS)
	assert.Equal(t, 5, track.TrackNumber)
	assert.NotNil(t, track.PreviewURL)
}

func TestConnection_Creation(t *testing.T) {
	users := []User{
		{ID: uuid.New(), Name: "User 1"},
		{ID: uuid.New(), Name: "User 2"},
	}

	edges := []Edge[User]{
		{Cursor: "cursor1", Node: users[0]},
		{Cursor: "cursor2", Node: users[1]},
	}

	connection := Connection[User]{
		TotalCount: 2,
		Edges:      edges,
		PageInfo: PageInfo{
			EndCursor:   stringPtr("cursor2"),
			HasNextPage: false,
		},
	}

	assert.Equal(t, 2, connection.TotalCount)
	assert.Len(t, connection.Edges, 2)
	assert.False(t, connection.PageInfo.HasNextPage)
	assert.Equal(t, "cursor2", *connection.PageInfo.EndCursor)
}

func TestCachedUserData_JSONMarshaling(t *testing.T) {
	data := &CachedUserData{
		UserID: uuid.New(),
		RecentlyPlayed: []SpotifyTrack{
			{ID: "track1", Name: "Track 1"},
		},
		TopTracks: []SpotifyTrack{
			{ID: "track2", Name: "Track 2"},
		},
		TopArtists: []SpotifyArtist{
			{ID: "artist1", Name: "Artist 1"},
		},
		SavedTracks: []string{"track3", "track4"},
		SavedAlbums: []string{"album1", "album2"},
		PlaylistIDs: []string{"playlist1", "playlist2"},
		LastUpdated: time.Now(),
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(data)
	require.NoError(t, err)
	assert.Contains(t, string(jsonData), "track1")
	assert.Contains(t, string(jsonData), "recently_played")

	// Test JSON unmarshaling
	var unmarshaled CachedUserData
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, data.UserID, unmarshaled.UserID)
	assert.Len(t, unmarshaled.RecentlyPlayed, 1)
	assert.Len(t, unmarshaled.SavedTracks, 2)
}

func TestCachedRecommendations_Creation(t *testing.T) {
	userID := uuid.New()
	recommendations := &CachedRecommendations{
		UserID: userID,
		RecommendedTracks: []SpotifyTrack{
			{ID: "rec_track1", Name: "Recommended Track 1"},
		},
		RecommendedAlbums: []SpotifyAlbum{
			{ID: "rec_album1", Name: "Recommended Album 1"},
		},
		RecommendedArtists: []SpotifyArtist{
			{ID: "rec_artist1", Name: "Recommended Artist 1"},
		},
		BasedOnGenres:  []string{"rock", "pop"},
		BasedOnArtists: []string{"artist1", "artist2"},
		GeneratedAt:    time.Now(),
	}

	assert.Equal(t, userID, recommendations.UserID)
	assert.Len(t, recommendations.RecommendedTracks, 1)
	assert.Len(t, recommendations.RecommendedAlbums, 1)
	assert.Len(t, recommendations.RecommendedArtists, 1)
	assert.Len(t, recommendations.BasedOnGenres, 2)
	assert.Len(t, recommendations.BasedOnArtists, 2)
}

// Test edge cases and validation
func TestReview_RatingValidation(t *testing.T) {
	// Test valid ratings (1-5)
	validRatings := []int{1, 2, 3, 4, 5}
	for _, rating := range validRatings {
		review := &Review{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			SpotifyID:   "spotify_track_123",
			SpotifyType: "track",
			Rating:      rating,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		assert.True(t, rating >= 1 && rating <= 5, "Rating should be between 1-5")
		assert.Equal(t, rating, review.Rating)
	}
}

func TestPlaylistTrack_PositionOrdering(t *testing.T) {
	playlistID := uuid.New()
	tracks := []PlaylistTrack{
		{ID: uuid.New(), PlaylistID: playlistID, Position: 3, SpotifyID: "track3"},
		{ID: uuid.New(), PlaylistID: playlistID, Position: 1, SpotifyID: "track1"},
		{ID: uuid.New(), PlaylistID: playlistID, Position: 2, SpotifyID: "track2"},
	}

	// Verify positions are different
	positions := make(map[int]bool)
	for _, track := range tracks {
		assert.False(t, positions[track.Position], "Positions should be unique")
		positions[track.Position] = true
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
