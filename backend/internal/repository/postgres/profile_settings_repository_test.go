package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/daedal00/muse/backend/internal/models"
	"github.com/google/uuid"
)

func TestProfileSettingsRepository_GetByUserID_NotFound(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewProfileSettingsRepository(testDB)
	ctx := context.Background()

	// Test getting settings for non-existent user returns defaults
	nonExistentID := uuid.New()
	settings, err := repo.GetByUserID(ctx, nonExistentID)
	if err != nil {
		t.Fatalf("Expected no error for non-existent user, got: %v", err)
	}

	if settings == nil {
		t.Fatal("Expected default settings, got nil")
	}

	if settings.UserID != nonExistentID {
		t.Errorf("Expected user ID %s, got %s", nonExistentID, settings.UserID)
	}

	if settings.Layout != models.ProfileLayoutGrid {
		t.Errorf("Expected default layout GRID, got %s", settings.Layout)
	}

	if settings.ShowSpotifyStats != true {
		t.Error("Expected ShowSpotifyStats to be true by default")
	}
}

func TestProfileSettingsRepository_Upsert_Create(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	userRepo := NewUserRepository(testDB)
	repo := NewProfileSettingsRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	defer cleanupTestUser(t, ctx, user.ID)

	err := userRepo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create profile settings
	primaryColor := "#ff5500"
	accentColor := "#000000"
	bgStyle := models.BackgroundStyleGradient
	bgValue := "linear-gradient(135deg, #667eea 0%, #764ba2 100%)"
	bioStyle := models.BioStyleQuote

	settings := &models.ProfileSettings{
		UserID:               user.ID,
		Layout:               models.ProfileLayoutBento,
		PrimaryColor:         &primaryColor,
		AccentColor:          &accentColor,
		BackgroundStyle:      &bgStyle,
		BackgroundValue:      &bgValue,
		PinnedAlbumIDs:       []uuid.UUID{},
		PinnedTrackIDs:       []uuid.UUID{},
		FeaturedArtistIDs:    []uuid.UUID{},
		SectionsOrder:        []string{"topTracks", "featuredArtists"},
		ShowSpotifyStats:     false,
		ShowListeningHistory: true,
		BioStyle:             &bioStyle,
		CustomTags:           []string{"indie", "electronic", "chill"},
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	err = repo.Upsert(ctx, settings)
	if err != nil {
		t.Fatalf("Failed to upsert profile settings: %v", err)
	}

	// Verify settings were created
	retrieved, err := repo.GetByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get profile settings: %v", err)
	}

	if retrieved.Layout != models.ProfileLayoutBento {
		t.Errorf("Expected layout BENTO, got %s", retrieved.Layout)
	}

	if *retrieved.PrimaryColor != "#ff5500" {
		t.Errorf("Expected primary color #ff5500, got %s", *retrieved.PrimaryColor)
	}

	if retrieved.ShowSpotifyStats != false {
		t.Error("Expected ShowSpotifyStats to be false")
	}

	if len(retrieved.CustomTags) != 3 {
		t.Errorf("Expected 3 custom tags, got %d", len(retrieved.CustomTags))
	}
}

func TestProfileSettingsRepository_Upsert_Update(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	userRepo := NewUserRepository(testDB)
	repo := NewProfileSettingsRepository(testDB)
	ctx := context.Background()

	// Create a test user
	user := setupTestUser(t)
	defer cleanupTestUser(t, ctx, user.ID)

	err := userRepo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create initial settings
	settings := models.DefaultProfileSettings(user.ID)
	err = repo.Upsert(ctx, settings)
	if err != nil {
		t.Fatalf("Failed to create initial settings: %v", err)
	}

	// Update settings
	newPrimaryColor := "#1DB954"
	settings.PrimaryColor = &newPrimaryColor
	settings.Layout = models.ProfileLayoutList
	settings.CustomTags = []string{"jazz", "blues"}

	err = repo.Upsert(ctx, settings)
	if err != nil {
		t.Fatalf("Failed to update profile settings: %v", err)
	}

	// Verify update
	retrieved, err := repo.GetByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get updated settings: %v", err)
	}

	if *retrieved.PrimaryColor != "#1DB954" {
		t.Errorf("Expected primary color #1DB954, got %s", *retrieved.PrimaryColor)
	}

	if retrieved.Layout != models.ProfileLayoutList {
		t.Errorf("Expected layout LIST, got %s", retrieved.Layout)
	}

	if len(retrieved.CustomTags) != 2 || retrieved.CustomTags[0] != "jazz" {
		t.Errorf("Expected custom tags [jazz, blues], got %v", retrieved.CustomTags)
	}
}

func TestProfileSettingsRepository_Delete(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	userRepo := NewUserRepository(testDB)
	repo := NewProfileSettingsRepository(testDB)
	ctx := context.Background()

	// Create a test user
	user := setupTestUser(t)
	defer cleanupTestUser(t, ctx, user.ID)

	err := userRepo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create profile settings
	settings := models.DefaultProfileSettings(user.ID)
	err = repo.Upsert(ctx, settings)
	if err != nil {
		t.Fatalf("Failed to create profile settings: %v", err)
	}

	// Delete settings
	err = repo.Delete(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to delete profile settings: %v", err)
	}

	// Verify deletion (should return defaults now)
	retrieved, err := repo.GetByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get settings after deletion: %v", err)
	}

	// Should be default values since we deleted
	if retrieved.Layout != models.ProfileLayoutGrid {
		t.Errorf("Expected default layout GRID after deletion, got %s", retrieved.Layout)
	}
}

func TestProfileSettingsRepository_Delete_NotFound(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewProfileSettingsRepository(testDB)
	ctx := context.Background()

	// Try to delete non-existent settings
	err := repo.Delete(ctx, uuid.New())
	if err == nil {
		t.Error("Expected error when deleting non-existent profile settings")
	}
}
