package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/daedal00/muse/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestPlaylist creates a test playlist for repository tests
func setupTestPlaylist(t *testing.T, creatorID uuid.UUID) *models.Playlist {
	t.Helper()

	return &models.Playlist{
		ID:          uuid.New(),
		Title:       "Test Playlist",
		Description: stringPtr("A test playlist description"),
		CoverImage:  stringPtr("https://example.com/cover.jpg"),
		CreatorID:   creatorID,
		IsPublic:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// cleanupTestPlaylist removes a test playlist from the database
func cleanupTestPlaylist(t *testing.T, ctx context.Context, playlistID uuid.UUID) {
	t.Helper()

	_, err := testDB.Pool.Exec(ctx, "DELETE FROM playlist_tracks WHERE playlist_id = $1", playlistID)
	if err != nil {
		t.Logf("Warning: Failed to cleanup playlist tracks: %v", err)
	}

	_, err = testDB.Pool.Exec(ctx, "DELETE FROM playlists WHERE id = $1", playlistID)
	if err != nil {
		t.Logf("Warning: Failed to cleanup test playlist: %v", err)
	}
}

func TestPlaylistRepository_Create(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewPlaylistRepository(testDB)
	ctx := context.Background()

	// Create a test user first (needed for foreign key)
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	playlist := setupTestPlaylist(t, user.ID)
	defer cleanupTestPlaylist(t, ctx, playlist.ID)

	err = repo.Create(ctx, playlist)
	require.NoError(t, err)

	// Verify playlist was created
	createdPlaylist, err := repo.GetByID(ctx, playlist.ID)
	require.NoError(t, err)

	assert.Equal(t, playlist.Title, createdPlaylist.Title)
	assert.Equal(t, playlist.CreatorID, createdPlaylist.CreatorID)
	assert.Equal(t, playlist.IsPublic, createdPlaylist.IsPublic)
	assert.Equal(t, *playlist.Description, *createdPlaylist.Description)
	assert.Equal(t, *playlist.CoverImage, *createdPlaylist.CoverImage)
}

func TestPlaylistRepository_GetByID(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewPlaylistRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	playlist := setupTestPlaylist(t, user.ID)
	defer cleanupTestPlaylist(t, ctx, playlist.ID)

	// Create playlist first
	err = repo.Create(ctx, playlist)
	require.NoError(t, err)

	// Test getting existing playlist
	foundPlaylist, err := repo.GetByID(ctx, playlist.ID)
	require.NoError(t, err)

	assert.Equal(t, playlist.ID, foundPlaylist.ID)
	assert.Equal(t, playlist.Title, foundPlaylist.Title)
	assert.Equal(t, playlist.CreatorID, foundPlaylist.CreatorID)
}

func TestPlaylistRepository_GetByID_NotFound(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewPlaylistRepository(testDB)
	ctx := context.Background()

	// Test getting non-existent playlist
	nonExistentID := uuid.New()
	_, err := repo.GetByID(ctx, nonExistentID)
	assert.Error(t, err)
}

func TestPlaylistRepository_GetByUserID(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewPlaylistRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	// Create multiple playlists for the user
	playlist1 := setupTestPlaylist(t, user.ID)
	playlist1.Title = "User Playlist 1"
	playlist2 := setupTestPlaylist(t, user.ID)
	playlist2.Title = "User Playlist 2"

	defer func() {
		cleanupTestPlaylist(t, ctx, playlist1.ID)
		cleanupTestPlaylist(t, ctx, playlist2.ID)
	}()

	err = repo.Create(ctx, playlist1)
	require.NoError(t, err)
	err = repo.Create(ctx, playlist2)
	require.NoError(t, err)

	// Test getting playlists by user ID (using GetByCreatorID)
	playlists, err := repo.GetByCreatorID(ctx, user.ID, 10, 0)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(playlists), 2)

	// Check if our playlists are in the results
	playlistIDs := make(map[uuid.UUID]bool)
	for _, playlist := range playlists {
		playlistIDs[playlist.ID] = true
	}

	assert.True(t, playlistIDs[playlist1.ID])
	assert.True(t, playlistIDs[playlist2.ID])
}

func TestPlaylistRepository_Update(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewPlaylistRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	playlist := setupTestPlaylist(t, user.ID)
	defer cleanupTestPlaylist(t, ctx, playlist.ID)

	// Create playlist first
	err = repo.Create(ctx, playlist)
	require.NoError(t, err)

	// Update playlist fields
	playlist.Title = "Updated Playlist Title"
	playlist.Description = stringPtr("Updated description")
	playlist.IsPublic = false

	err = repo.Update(ctx, playlist)
	require.NoError(t, err)

	// Verify updates
	updatedPlaylist, err := repo.GetByID(ctx, playlist.ID)
	require.NoError(t, err)

	assert.Equal(t, "Updated Playlist Title", updatedPlaylist.Title)
	assert.Equal(t, "Updated description", *updatedPlaylist.Description)
	assert.False(t, updatedPlaylist.IsPublic)
}

func TestPlaylistRepository_Delete(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewPlaylistRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	playlist := setupTestPlaylist(t, user.ID)

	// Create playlist first
	err = repo.Create(ctx, playlist)
	require.NoError(t, err)

	// Delete playlist
	err = repo.Delete(ctx, playlist.ID)
	require.NoError(t, err)

	// Verify playlist was deleted
	_, err = repo.GetByID(ctx, playlist.ID)
	assert.Error(t, err)
}

func TestPlaylistRepository_AddTrack(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewPlaylistRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	playlist := setupTestPlaylist(t, user.ID)
	defer cleanupTestPlaylist(t, ctx, playlist.ID)

	// Create playlist first
	err = repo.Create(ctx, playlist)
	require.NoError(t, err)

	// Add track to playlist (using Spotify ID)
	spotifyID := "spotify_track_123"
	position := 1

	err = repo.AddTrack(ctx, playlist.ID, spotifyID, position, user.ID)
	require.NoError(t, err)

	// Verify track was added using GetPlaylistTracks
	tracks, err := repo.GetPlaylistTracks(ctx, playlist.ID, 10, 0)
	require.NoError(t, err)

	assert.Len(t, tracks, 1)
	assert.Equal(t, spotifyID, tracks[0].SpotifyID)
	assert.Equal(t, position, tracks[0].Position)
	assert.Equal(t, user.ID, tracks[0].AddedByUserID)
}

func TestPlaylistRepository_AddMultipleTracks(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewPlaylistRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	playlist := setupTestPlaylist(t, user.ID)
	defer cleanupTestPlaylist(t, ctx, playlist.ID)

	// Create playlist first
	err = repo.Create(ctx, playlist)
	require.NoError(t, err)

	// Add multiple tracks with different positions
	spotifyIDs := []string{"spotify_track_1", "spotify_track_2", "spotify_track_3"}

	for i, spotifyID := range spotifyIDs {
		err = repo.AddTrack(ctx, playlist.ID, spotifyID, i+1, user.ID)
		require.NoError(t, err)
	}

	// Verify tracks were added in correct order
	tracks, err := repo.GetPlaylistTracks(ctx, playlist.ID, 10, 0)
	require.NoError(t, err)

	assert.Len(t, tracks, 3)

	// Tracks should be ordered by position
	assert.Equal(t, 1, tracks[0].Position)
	assert.Equal(t, 2, tracks[1].Position)
	assert.Equal(t, 3, tracks[2].Position)

	// Verify Spotify IDs match
	assert.Equal(t, spotifyIDs[0], tracks[0].SpotifyID)
	assert.Equal(t, spotifyIDs[1], tracks[1].SpotifyID)
	assert.Equal(t, spotifyIDs[2], tracks[2].SpotifyID)
}

func TestPlaylistRepository_RemoveTrack(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewPlaylistRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	playlist := setupTestPlaylist(t, user.ID)
	defer cleanupTestPlaylist(t, ctx, playlist.ID)

	// Create playlist first
	err = repo.Create(ctx, playlist)
	require.NoError(t, err)

	// Add track to playlist
	spotifyID := "spotify_track_123"
	err = repo.AddTrack(ctx, playlist.ID, spotifyID, 1, user.ID)
	require.NoError(t, err)

	// Remove track from playlist
	err = repo.RemoveTrack(ctx, playlist.ID, spotifyID)
	require.NoError(t, err)

	// Verify track was removed
	tracks, err := repo.GetPlaylistTracks(ctx, playlist.ID, 10, 0)
	require.NoError(t, err)

	assert.Len(t, tracks, 0)
}

func TestPlaylistRepository_GetTracks_Pagination(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewPlaylistRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	playlist := setupTestPlaylist(t, user.ID)
	defer cleanupTestPlaylist(t, ctx, playlist.ID)

	// Create playlist first
	err = repo.Create(ctx, playlist)
	require.NoError(t, err)

	// Add 5 tracks
	for i := 1; i <= 5; i++ {
		spotifyID := fmt.Sprintf("spotify_track_%d", i)
		err = repo.AddTrack(ctx, playlist.ID, spotifyID, i, user.ID)
		require.NoError(t, err)
	}

	// Test pagination - first page
	tracks, err := repo.GetPlaylistTracks(ctx, playlist.ID, 3, 0)
	require.NoError(t, err)
	assert.Len(t, tracks, 3)
	assert.Equal(t, 1, tracks[0].Position)
	assert.Equal(t, 2, tracks[1].Position)
	assert.Equal(t, 3, tracks[2].Position)

	// Test pagination - second page
	tracks, err = repo.GetPlaylistTracks(ctx, playlist.ID, 3, 3)
	require.NoError(t, err)
	assert.Len(t, tracks, 2) // Only 2 remaining tracks
	assert.Equal(t, 4, tracks[0].Position)
	assert.Equal(t, 5, tracks[1].Position)
}

func TestPlaylistRepository_UpdateTrackPosition(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewPlaylistRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	playlist := setupTestPlaylist(t, user.ID)
	defer cleanupTestPlaylist(t, ctx, playlist.ID)

	// Create playlist first
	err = repo.Create(ctx, playlist)
	require.NoError(t, err)

	// Add track to playlist
	spotifyID := "spotify_track_123"
	err = repo.AddTrack(ctx, playlist.ID, spotifyID, 1, user.ID)
	require.NoError(t, err)

	// Update track position using ReorderTracks
	trackPositions := map[string]int{
		spotifyID: 5,
	}
	err = repo.ReorderTracks(ctx, playlist.ID, trackPositions)
	require.NoError(t, err)

	// Verify position was updated
	tracks, err := repo.GetPlaylistTracks(ctx, playlist.ID, 10, 0)
	require.NoError(t, err)

	assert.Len(t, tracks, 1)
	assert.Equal(t, 5, tracks[0].Position)
}

func TestPlaylistRepository_GetPublicPlaylists(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewPlaylistRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	// Create public and private playlists
	publicPlaylist := setupTestPlaylist(t, user.ID)
	publicPlaylist.Title = "Public Playlist"
	publicPlaylist.IsPublic = true

	privatePlaylist := setupTestPlaylist(t, user.ID)
	privatePlaylist.Title = "Private Playlist"
	privatePlaylist.IsPublic = false

	defer func() {
		cleanupTestPlaylist(t, ctx, publicPlaylist.ID)
		cleanupTestPlaylist(t, ctx, privatePlaylist.ID)
	}()

	err = repo.Create(ctx, publicPlaylist)
	require.NoError(t, err)
	err = repo.Create(ctx, privatePlaylist)
	require.NoError(t, err)

	// Test getting only public playlists
	playlists, err := repo.GetPublicPlaylists(ctx, 10, 0)
	require.NoError(t, err)

	// Find our public playlist in the results
	var foundPublic bool
	var foundPrivate bool
	for _, playlist := range playlists {
		if playlist.ID == publicPlaylist.ID {
			foundPublic = true
			assert.True(t, playlist.IsPublic)
		}
		if playlist.ID == privatePlaylist.ID {
			foundPrivate = true
		}
	}

	assert.True(t, foundPublic, "Public playlist should be in results")
	assert.False(t, foundPrivate, "Private playlist should not be in results")
}

func TestPlaylistRepository_SearchByTitle(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewPlaylistRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	// Create playlists with different titles
	rockPlaylist := setupTestPlaylist(t, user.ID)
	rockPlaylist.Title = "Best Rock Songs"
	rockPlaylist.IsPublic = true

	popPlaylist := setupTestPlaylist(t, user.ID)
	popPlaylist.Title = "Pop Hits 2023"
	popPlaylist.IsPublic = true

	defer func() {
		cleanupTestPlaylist(t, ctx, rockPlaylist.ID)
		cleanupTestPlaylist(t, ctx, popPlaylist.ID)
	}()

	err = repo.Create(ctx, rockPlaylist)
	require.NoError(t, err)
	err = repo.Create(ctx, popPlaylist)
	require.NoError(t, err)

	// Test by listing all public playlists and filtering (since SearchByTitle method doesn't exist)
	playlists, err := repo.GetPublicPlaylists(ctx, 10, 0)
	require.NoError(t, err)

	// Should find the rock playlist
	var foundRock bool
	for _, playlist := range playlists {
		if playlist.ID == rockPlaylist.ID {
			foundRock = true
		}
	}

	assert.True(t, foundRock, "Rock playlist should be found in public playlists")
}

// Benchmark tests
func BenchmarkPlaylistRepository_Create(b *testing.B) {
	if testDB == nil {
		b.Skip("Database not available")
	}

	repo := NewPlaylistRepository(testDB)
	userRepo := NewUserRepository(testDB)
	ctx := context.Background()

	// Create a test user for benchmarking
	user := setupTestUser(&testing.T{})
	err := userRepo.Create(ctx, user)
	if err != nil {
		b.Fatal(err)
	}
	defer cleanupTestUser(&testing.T{}, ctx, user.ID)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		playlist := setupTestPlaylist(&testing.T{}, user.ID)
		playlist.Title = fmt.Sprintf("Benchmark Playlist %d", i)

		err := repo.Create(ctx, playlist)
		if err != nil {
			b.Fatal(err)
		}

		// Cleanup
		cleanupTestPlaylist(&testing.T{}, ctx, playlist.ID)
	}
}

func BenchmarkPlaylistRepository_GetByID(b *testing.B) {
	if testDB == nil {
		b.Skip("Database not available")
	}

	repo := NewPlaylistRepository(testDB)
	userRepo := NewUserRepository(testDB)
	ctx := context.Background()

	// Setup test data
	user := setupTestUser(&testing.T{})
	err := userRepo.Create(ctx, user)
	if err != nil {
		b.Fatal(err)
	}
	defer cleanupTestUser(&testing.T{}, ctx, user.ID)

	playlist := setupTestPlaylist(&testing.T{}, user.ID)
	err = repo.Create(ctx, playlist)
	if err != nil {
		b.Fatal(err)
	}
	defer cleanupTestPlaylist(&testing.T{}, ctx, playlist.ID)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := repo.GetByID(ctx, playlist.ID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPlaylistRepository_AddTrack(b *testing.B) {
	if testDB == nil {
		b.Skip("Database not available")
	}

	repo := NewPlaylistRepository(testDB)
	userRepo := NewUserRepository(testDB)
	ctx := context.Background()

	// Setup test data
	user := setupTestUser(&testing.T{})
	err := userRepo.Create(ctx, user)
	if err != nil {
		b.Fatal(err)
	}
	defer cleanupTestUser(&testing.T{}, ctx, user.ID)

	playlist := setupTestPlaylist(&testing.T{}, user.ID)
	err = repo.Create(ctx, playlist)
	if err != nil {
		b.Fatal(err)
	}
	defer cleanupTestPlaylist(&testing.T{}, ctx, playlist.ID)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		spotifyID := fmt.Sprintf("spotify_track_bench_%d", i)
		err := repo.AddTrack(ctx, playlist.ID, spotifyID, i+1, user.ID)
		if err != nil {
			b.Fatal(err)
		}
	}
}
