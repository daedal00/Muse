package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/daedal00/muse/backend/internal/models"
	"github.com/google/uuid"
)

func setupTestAlbum(t *testing.T, artistID uuid.UUID) *models.Album {
	t.Helper()

	releaseDate, _ := time.Parse("2006-01-02", "2023-01-15")

	return &models.Album{
		ID:          uuid.New(),
		SpotifyID:   stringPtr(fmt.Sprintf("spotify_album_%s", uuid.New().String()[:8])),
		Title:       fmt.Sprintf("Test Album %s", uuid.New().String()[:8]),
		ArtistID:    artistID,
		ReleaseDate: &releaseDate,
		CoverImage:  stringPtr("https://example.com/cover.jpg"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func cleanupTestAlbum(t *testing.T, ctx context.Context, albumID uuid.UUID) {
	t.Helper()

	_, err := testDB.Pool.Exec(ctx, "DELETE FROM albums WHERE id = $1", albumID)
	if err != nil {
		t.Logf("Warning: Failed to cleanup test album: %v", err)
	}
}

func timePtr(s string) *time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return &t
}

func TestAlbumRepository_Create(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	albumRepo := NewAlbumRepository(testDB)
	artistRepo := NewArtistRepository(testDB)
	ctx := context.Background()

	// Create test artist first
	artist := setupTestArtist(t)
	defer cleanupTestArtist(t, ctx, artist.ID)

	err := artistRepo.Create(ctx, artist)
	if err != nil {
		t.Fatalf("Failed to create test artist: %v", err)
	}

	// Create test album
	album := setupTestAlbum(t, artist.ID)
	defer cleanupTestAlbum(t, ctx, album.ID)

	err = albumRepo.Create(ctx, album)
	if err != nil {
		t.Fatalf("Failed to create album: %v", err)
	}

	// Verify album was created
	createdAlbum, err := albumRepo.GetByID(ctx, album.ID)
	if err != nil {
		t.Fatalf("Failed to get created album: %v", err)
	}

	if createdAlbum.Title != album.Title {
		t.Errorf("Expected title %s, got %s", album.Title, createdAlbum.Title)
	}
	if createdAlbum.ArtistID != album.ArtistID {
		t.Errorf("Expected artist ID %s, got %s", album.ArtistID, createdAlbum.ArtistID)
	}
}

func TestAlbumRepository_GetByID(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	albumRepo := NewAlbumRepository(testDB)
	artistRepo := NewArtistRepository(testDB)
	ctx := context.Background()

	// Create test artist first
	artist := setupTestArtist(t)
	defer cleanupTestArtist(t, ctx, artist.ID)

	err := artistRepo.Create(ctx, artist)
	if err != nil {
		t.Fatalf("Failed to create test artist: %v", err)
	}

	// Create test album
	album := setupTestAlbum(t, artist.ID)
	defer cleanupTestAlbum(t, ctx, album.ID)

	err = albumRepo.Create(ctx, album)
	if err != nil {
		t.Fatalf("Failed to create album: %v", err)
	}

	// Test getting existing album
	foundAlbum, err := albumRepo.GetByID(ctx, album.ID)
	if err != nil {
		t.Fatalf("Failed to get album by ID: %v", err)
	}

	if foundAlbum.ID != album.ID {
		t.Errorf("Expected ID %s, got %s", album.ID, foundAlbum.ID)
	}
}

func TestAlbumRepository_GetByID_NotFound(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewAlbumRepository(testDB)
	ctx := context.Background()

	// Test getting non-existent album
	nonExistentID := uuid.New()
	_, err := repo.GetByID(ctx, nonExistentID)
	if err == nil {
		t.Error("Expected error when getting non-existent album")
	}
}

func TestAlbumRepository_GetBySpotifyID(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	albumRepo := NewAlbumRepository(testDB)
	artistRepo := NewArtistRepository(testDB)
	ctx := context.Background()

	// Create test artist first
	artist := setupTestArtist(t)
	defer cleanupTestArtist(t, ctx, artist.ID)

	err := artistRepo.Create(ctx, artist)
	if err != nil {
		t.Fatalf("Failed to create test artist: %v", err)
	}

	// Create test album
	album := setupTestAlbum(t, artist.ID)
	defer cleanupTestAlbum(t, ctx, album.ID)

	err = albumRepo.Create(ctx, album)
	if err != nil {
		t.Fatalf("Failed to create album: %v", err)
	}

	// Test getting album by Spotify ID
	if album.SpotifyID != nil {
		foundAlbum, err := albumRepo.GetBySpotifyID(ctx, *album.SpotifyID)
		if err != nil {
			t.Fatalf("Failed to get album by Spotify ID: %v", err)
		}

		if foundAlbum.ID != album.ID {
			t.Errorf("Expected ID %s, got %s", album.ID, foundAlbum.ID)
		}
	}
}

func TestAlbumRepository_GetByArtistID(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	albumRepo := NewAlbumRepository(testDB)
	artistRepo := NewArtistRepository(testDB)
	ctx := context.Background()

	// Create test artist first
	artist := setupTestArtist(t)
	defer cleanupTestArtist(t, ctx, artist.ID)

	err := artistRepo.Create(ctx, artist)
	if err != nil {
		t.Fatalf("Failed to create test artist: %v", err)
	}

	// Create multiple test albums for the same artist
	var albumIDs []uuid.UUID
	for i := 0; i < 3; i++ {
		album := setupTestAlbum(t, artist.ID)
		albumIDs = append(albumIDs, album.ID)

		err := albumRepo.Create(ctx, album)
		if err != nil {
			t.Fatalf("Failed to create album %d: %v", i, err)
		}
	}

	// Cleanup
	defer func() {
		for _, id := range albumIDs {
			cleanupTestAlbum(t, ctx, id)
		}
	}()

	// Test getting albums by artist ID
	albums, err := albumRepo.GetByArtistID(ctx, artist.ID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to get albums by artist ID: %v", err)
	}

	if len(albums) < 3 {
		t.Errorf("Expected at least 3 albums, got %d", len(albums))
	}

	// Verify all albums belong to the artist
	for _, album := range albums {
		if album.ArtistID != artist.ID {
			t.Errorf("Expected artist ID %s, got %s", artist.ID, album.ArtistID)
		}
	}
}

func TestAlbumRepository_Update(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	albumRepo := NewAlbumRepository(testDB)
	artistRepo := NewArtistRepository(testDB)
	ctx := context.Background()

	// Create test artist first
	artist := setupTestArtist(t)
	defer cleanupTestArtist(t, ctx, artist.ID)

	err := artistRepo.Create(ctx, artist)
	if err != nil {
		t.Fatalf("Failed to create test artist: %v", err)
	}

	// Create test album
	album := setupTestAlbum(t, artist.ID)
	defer cleanupTestAlbum(t, ctx, album.ID)

	err = albumRepo.Create(ctx, album)
	if err != nil {
		t.Fatalf("Failed to create album: %v", err)
	}

	// Update album
	album.Title = "Updated Album Title"
	album.ReleaseDate = timePtr("2024-01-15")

	err = albumRepo.Update(ctx, album)
	if err != nil {
		t.Fatalf("Failed to update album: %v", err)
	}

	// Verify update
	updatedAlbum, err := albumRepo.GetByID(ctx, album.ID)
	if err != nil {
		t.Fatalf("Failed to get updated album: %v", err)
	}

	if updatedAlbum.Title != "Updated Album Title" {
		t.Errorf("Expected title 'Updated Album Title', got %s", updatedAlbum.Title)
	}
	if updatedAlbum.ReleaseDate != nil {
		expectedDate, _ := time.Parse("2006-01-02", "2024-01-15")
		if !updatedAlbum.ReleaseDate.Equal(expectedDate) {
			t.Errorf("Expected release date '2024-01-15', got %s", updatedAlbum.ReleaseDate.Format("2006-01-02"))
		}
	}
}

func TestAlbumRepository_Delete(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	albumRepo := NewAlbumRepository(testDB)
	artistRepo := NewArtistRepository(testDB)
	ctx := context.Background()

	// Create test artist first
	artist := setupTestArtist(t)
	defer cleanupTestArtist(t, ctx, artist.ID)

	err := artistRepo.Create(ctx, artist)
	if err != nil {
		t.Fatalf("Failed to create test artist: %v", err)
	}

	// Create test album
	album := setupTestAlbum(t, artist.ID)

	err = albumRepo.Create(ctx, album)
	if err != nil {
		t.Fatalf("Failed to create album: %v", err)
	}

	// Delete album
	err = albumRepo.Delete(ctx, album.ID)
	if err != nil {
		t.Fatalf("Failed to delete album: %v", err)
	}

	// Verify deletion
	_, err = albumRepo.GetByID(ctx, album.ID)
	if err == nil {
		t.Error("Expected error when getting deleted album")
	}
}

func TestAlbumRepository_List(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	albumRepo := NewAlbumRepository(testDB)
	artistRepo := NewArtistRepository(testDB)
	ctx := context.Background()

	// Create test artist first
	artist := setupTestArtist(t)
	defer cleanupTestArtist(t, ctx, artist.ID)

	err := artistRepo.Create(ctx, artist)
	if err != nil {
		t.Fatalf("Failed to create test artist: %v", err)
	}

	// Create multiple test albums
	var albumIDs []uuid.UUID
	for i := 0; i < 3; i++ {
		album := setupTestAlbum(t, artist.ID)
		albumIDs = append(albumIDs, album.ID)

		err := albumRepo.Create(ctx, album)
		if err != nil {
			t.Fatalf("Failed to create album %d: %v", i, err)
		}
	}

	// Cleanup
	defer func() {
		for _, id := range albumIDs {
			cleanupTestAlbum(t, ctx, id)
		}
	}()

	// Test listing albums
	albums, err := albumRepo.List(ctx, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list albums: %v", err)
	}

	if len(albums) < 3 {
		t.Errorf("Expected at least 3 albums, got %d", len(albums))
	}

	// Test with limit
	albums, err = albumRepo.List(ctx, 2, 0)
	if err != nil {
		t.Fatalf("Failed to list albums with limit: %v", err)
	}

	if len(albums) > 2 {
		t.Errorf("Expected at most 2 albums, got %d", len(albums))
	}
}
