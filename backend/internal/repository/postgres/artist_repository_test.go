package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/daedal00/muse/backend/internal/models"
	"github.com/google/uuid"
)

func setupTestArtist(t *testing.T) *models.Artist {
	t.Helper()
	
	return &models.Artist{
		ID:        uuid.New(),
		SpotifyID: stringPtr(fmt.Sprintf("spotify_%s", uuid.New().String()[:8])),
		Name:      fmt.Sprintf("Test Artist %s", uuid.New().String()[:8]),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func cleanupTestArtist(t *testing.T, ctx context.Context, artistID uuid.UUID) {
	t.Helper()
	
	_, err := testDB.Pool.Exec(ctx, "DELETE FROM artists WHERE id = $1", artistID)
	if err != nil {
		t.Logf("Warning: Failed to cleanup test artist: %v", err)
	}
}

func TestArtistRepository_Create(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}
	
	repo := NewArtistRepository(testDB)
	ctx := context.Background()
	
	artist := setupTestArtist(t)
	defer cleanupTestArtist(t, ctx, artist.ID)
	
	err := repo.Create(ctx, artist)
	if err != nil {
		t.Fatalf("Failed to create artist: %v", err)
	}
	
	// Verify artist was created
	createdArtist, err := repo.GetByID(ctx, artist.ID)
	if err != nil {
		t.Fatalf("Failed to get created artist: %v", err)
	}
	
	if createdArtist.Name != artist.Name {
		t.Errorf("Expected name %s, got %s", artist.Name, createdArtist.Name)
	}
	if createdArtist.SpotifyID != nil && artist.SpotifyID != nil && *createdArtist.SpotifyID != *artist.SpotifyID {
		t.Errorf("Expected Spotify ID %s, got %s", *artist.SpotifyID, *createdArtist.SpotifyID)
	}
}

func TestArtistRepository_GetByID(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}
	
	repo := NewArtistRepository(testDB)
	ctx := context.Background()
	
	artist := setupTestArtist(t)
	defer cleanupTestArtist(t, ctx, artist.ID)
	
	// Create artist first
	err := repo.Create(ctx, artist)
	if err != nil {
		t.Fatalf("Failed to create artist: %v", err)
	}
	
	// Test getting existing artist
	foundArtist, err := repo.GetByID(ctx, artist.ID)
	if err != nil {
		t.Fatalf("Failed to get artist by ID: %v", err)
	}
	
	if foundArtist.ID != artist.ID {
		t.Errorf("Expected ID %s, got %s", artist.ID, foundArtist.ID)
	}
}

func TestArtistRepository_GetByID_NotFound(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}
	
	repo := NewArtistRepository(testDB)
	ctx := context.Background()
	
	// Test getting non-existent artist
	nonExistentID := uuid.New()
	_, err := repo.GetByID(ctx, nonExistentID)
	if err == nil {
		t.Error("Expected error when getting non-existent artist")
	}
}

func TestArtistRepository_GetBySpotifyID(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}
	
	repo := NewArtistRepository(testDB)
	ctx := context.Background()
	
	artist := setupTestArtist(t)
	defer cleanupTestArtist(t, ctx, artist.ID)
	
	// Create artist first
	err := repo.Create(ctx, artist)
	if err != nil {
		t.Fatalf("Failed to create artist: %v", err)
	}
	
	// Test getting artist by Spotify ID
	if artist.SpotifyID != nil {
		foundArtist, err := repo.GetBySpotifyID(ctx, *artist.SpotifyID)
		if err != nil {
			t.Fatalf("Failed to get artist by Spotify ID: %v", err)
		}
		
		if foundArtist.ID != artist.ID {
			t.Errorf("Expected ID %s, got %s", artist.ID, foundArtist.ID)
		}
	}
}

func TestArtistRepository_Update(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}
	
	repo := NewArtistRepository(testDB)
	ctx := context.Background()
	
	artist := setupTestArtist(t)
	defer cleanupTestArtist(t, ctx, artist.ID)
	
	// Create artist first
	err := repo.Create(ctx, artist)
	if err != nil {
		t.Fatalf("Failed to create artist: %v", err)
	}
	
	// Update artist
	artist.Name = "Updated Artist Name"
	
	err = repo.Update(ctx, artist)
	if err != nil {
		t.Fatalf("Failed to update artist: %v", err)
	}
	
	// Verify update
	updatedArtist, err := repo.GetByID(ctx, artist.ID)
	if err != nil {
		t.Fatalf("Failed to get updated artist: %v", err)
	}
	
	if updatedArtist.Name != "Updated Artist Name" {
		t.Errorf("Expected name 'Updated Artist Name', got %s", updatedArtist.Name)
	}
}

func TestArtistRepository_Delete(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}
	
	repo := NewArtistRepository(testDB)
	ctx := context.Background()
	
	artist := setupTestArtist(t)
	
	// Create artist first
	err := repo.Create(ctx, artist)
	if err != nil {
		t.Fatalf("Failed to create artist: %v", err)
	}
	
	// Delete artist
	err = repo.Delete(ctx, artist.ID)
	if err != nil {
		t.Fatalf("Failed to delete artist: %v", err)
	}
	
	// Verify deletion
	_, err = repo.GetByID(ctx, artist.ID)
	if err == nil {
		t.Error("Expected error when getting deleted artist")
	}
}

func TestArtistRepository_List(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}
	
	repo := NewArtistRepository(testDB)
	ctx := context.Background()
	
	// Create multiple test artists
	var artistIDs []uuid.UUID
	for i := 0; i < 3; i++ {
		artist := setupTestArtist(t)
		artistIDs = append(artistIDs, artist.ID)
		
		err := repo.Create(ctx, artist)
		if err != nil {
			t.Fatalf("Failed to create artist %d: %v", i, err)
		}
	}
	
	// Cleanup
	defer func() {
		for _, id := range artistIDs {
			cleanupTestArtist(t, ctx, id)
		}
	}()
	
	// Test listing artists
	artists, err := repo.List(ctx, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list artists: %v", err)
	}
	
	if len(artists) < 3 {
		t.Errorf("Expected at least 3 artists, got %d", len(artists))
	}
	
	// Test with limit
	artists, err = repo.List(ctx, 2, 0)
	if err != nil {
		t.Fatalf("Failed to list artists with limit: %v", err)
	}
	
	if len(artists) > 2 {
		t.Errorf("Expected at most 2 artists, got %d", len(artists))
	}
} 