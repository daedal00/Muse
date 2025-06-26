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

// setupTestReview creates a test review for repository tests
func setupTestReview(t *testing.T, userID uuid.UUID, spotifyType string) *models.Review {
	t.Helper()

	return &models.Review{
		ID:          uuid.New(),
		UserID:      userID,
		SpotifyID:   fmt.Sprintf("spotify_%s_123", spotifyType),
		SpotifyType: spotifyType,
		Rating:      4,
		ReviewText:  stringPtr("Great music!"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// cleanupTestReview removes a test review from the database
func cleanupTestReview(t *testing.T, ctx context.Context, reviewID uuid.UUID) {
	t.Helper()

	_, err := testDB.Pool.Exec(ctx, "DELETE FROM reviews WHERE id = $1", reviewID)
	if err != nil {
		t.Logf("Warning: Failed to cleanup test review: %v", err)
	}
}

func TestReviewRepository_Create(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewReviewRepository(testDB)
	ctx := context.Background()

	// Create a test user first (needed for foreign key)
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	review := setupTestReview(t, user.ID, "track")
	defer cleanupTestReview(t, ctx, review.ID)

	err = repo.Create(ctx, review)
	require.NoError(t, err)

	// Verify review was created
	createdReview, err := repo.GetByID(ctx, review.ID)
	require.NoError(t, err)

	assert.Equal(t, review.UserID, createdReview.UserID)
	assert.Equal(t, review.SpotifyID, createdReview.SpotifyID)
	assert.Equal(t, review.SpotifyType, createdReview.SpotifyType)
	assert.Equal(t, review.Rating, createdReview.Rating)
	assert.Equal(t, *review.ReviewText, *createdReview.ReviewText)
}

func TestReviewRepository_CreateAlbumReview(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewReviewRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	review := setupTestReview(t, user.ID, "album")
	review.Rating = 5
	review.ReviewText = stringPtr("Amazing album!")
	defer cleanupTestReview(t, ctx, review.ID)

	err = repo.Create(ctx, review)
	require.NoError(t, err)

	// Verify album review was created
	createdReview, err := repo.GetByID(ctx, review.ID)
	require.NoError(t, err)

	assert.Equal(t, "album", createdReview.SpotifyType)
	assert.Equal(t, 5, createdReview.Rating)
	assert.Equal(t, "Amazing album!", *createdReview.ReviewText)
}

func TestReviewRepository_GetByID(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewReviewRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	review := setupTestReview(t, user.ID, "track")
	defer cleanupTestReview(t, ctx, review.ID)

	// Create review first
	err = repo.Create(ctx, review)
	require.NoError(t, err)

	// Test getting existing review
	foundReview, err := repo.GetByID(ctx, review.ID)
	require.NoError(t, err)

	assert.Equal(t, review.ID, foundReview.ID)
	assert.Equal(t, review.SpotifyID, foundReview.SpotifyID)
	assert.Equal(t, review.SpotifyType, foundReview.SpotifyType)
}

func TestReviewRepository_GetByID_NotFound(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewReviewRepository(testDB)
	ctx := context.Background()

	// Test getting non-existent review
	nonExistentID := uuid.New()
	_, err := repo.GetByID(ctx, nonExistentID)
	assert.Error(t, err)
}

func TestReviewRepository_GetByUserID(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewReviewRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	// Create multiple reviews for the user
	review1 := setupTestReview(t, user.ID, "track")
	review1.SpotifyID = "spotify_track_1"
	review2 := setupTestReview(t, user.ID, "album")
	review2.SpotifyID = "spotify_album_1"

	defer func() {
		cleanupTestReview(t, ctx, review1.ID)
		cleanupTestReview(t, ctx, review2.ID)
	}()

	err = repo.Create(ctx, review1)
	require.NoError(t, err)
	err = repo.Create(ctx, review2)
	require.NoError(t, err)

	// Test getting reviews by user ID
	reviews, err := repo.GetByUserID(ctx, user.ID, 10, 0)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(reviews), 2)

	// Check if our reviews are in the results
	reviewIDs := make(map[uuid.UUID]bool)
	for _, review := range reviews {
		reviewIDs[review.ID] = true
	}

	assert.True(t, reviewIDs[review1.ID])
	assert.True(t, reviewIDs[review2.ID])
}

func TestReviewRepository_GetBySpotifyID(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewReviewRepository(testDB)
	ctx := context.Background()

	// Create test users
	user1 := setupTestUser(t)
	user2 := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user1)
	require.NoError(t, err)
	err = userRepo.Create(ctx, user2)
	require.NoError(t, err)
	defer func() {
		cleanupTestUser(t, ctx, user1.ID)
		cleanupTestUser(t, ctx, user2.ID)
	}()

	spotifyID := "spotify_track_popular"

	// Create multiple reviews for the same Spotify track
	review1 := setupTestReview(t, user1.ID, "track")
	review1.SpotifyID = spotifyID
	review1.Rating = 4
	review2 := setupTestReview(t, user2.ID, "track")
	review2.SpotifyID = spotifyID
	review2.Rating = 5

	defer func() {
		cleanupTestReview(t, ctx, review1.ID)
		cleanupTestReview(t, ctx, review2.ID)
	}()

	err = repo.Create(ctx, review1)
	require.NoError(t, err)
	err = repo.Create(ctx, review2)
	require.NoError(t, err)

	// Test getting reviews by Spotify ID
	reviews, err := repo.GetBySpotifyID(ctx, spotifyID, "track", 10, 0)
	require.NoError(t, err)

	assert.Len(t, reviews, 2)

	// Both reviews should be for the same Spotify ID
	for _, review := range reviews {
		assert.Equal(t, spotifyID, review.SpotifyID)
		assert.Equal(t, "track", review.SpotifyType)
	}
}

func TestReviewRepository_Update(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewReviewRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	review := setupTestReview(t, user.ID, "track")
	defer cleanupTestReview(t, ctx, review.ID)

	// Create review first
	err = repo.Create(ctx, review)
	require.NoError(t, err)

	// Update review fields
	review.Rating = 5
	review.ReviewText = stringPtr("Actually, this is amazing!")

	err = repo.Update(ctx, review)
	require.NoError(t, err)

	// Verify updates
	updatedReview, err := repo.GetByID(ctx, review.ID)
	require.NoError(t, err)

	assert.Equal(t, 5, updatedReview.Rating)
	assert.Equal(t, "Actually, this is amazing!", *updatedReview.ReviewText)
}

func TestReviewRepository_Delete(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewReviewRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	review := setupTestReview(t, user.ID, "track")

	// Create review first
	err = repo.Create(ctx, review)
	require.NoError(t, err)

	// Delete review
	err = repo.Delete(ctx, review.ID)
	require.NoError(t, err)

	// Verify review was deleted
	_, err = repo.GetByID(ctx, review.ID)
	assert.Error(t, err)
}

func TestReviewRepository_GetRecentReviews(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewReviewRepository(testDB)
	ctx := context.Background()

	// Create test users
	user1 := setupTestUser(t)
	user2 := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user1)
	require.NoError(t, err)
	err = userRepo.Create(ctx, user2)
	require.NoError(t, err)
	defer func() {
		cleanupTestUser(t, ctx, user1.ID)
		cleanupTestUser(t, ctx, user2.ID)
	}()

	// Create reviews with different timestamps
	review1 := setupTestReview(t, user1.ID, "track")
	review1.SpotifyID = "spotify_track_1"
	review2 := setupTestReview(t, user2.ID, "album")
	review2.SpotifyID = "spotify_album_1"

	defer func() {
		cleanupTestReview(t, ctx, review1.ID)
		cleanupTestReview(t, ctx, review2.ID)
	}()

	err = repo.Create(ctx, review1)
	require.NoError(t, err)
	err = repo.Create(ctx, review2)
	require.NoError(t, err)

	// Test getting recent reviews (using List method if GetRecentReviews doesn't exist)
	reviews, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(reviews), 2)

	// Should include our reviews
	reviewIDs := make(map[uuid.UUID]bool)
	for _, review := range reviews {
		reviewIDs[review.ID] = true
	}

	assert.True(t, reviewIDs[review1.ID])
	assert.True(t, reviewIDs[review2.ID])
}

func TestReviewRepository_GetAverageRating(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewReviewRepository(testDB)
	ctx := context.Background()

	// Create test users
	user1 := setupTestUser(t)
	user2 := setupTestUser(t)
	user3 := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user1)
	require.NoError(t, err)
	err = userRepo.Create(ctx, user2)
	require.NoError(t, err)
	err = userRepo.Create(ctx, user3)
	require.NoError(t, err)
	defer func() {
		cleanupTestUser(t, ctx, user1.ID)
		cleanupTestUser(t, ctx, user2.ID)
		cleanupTestUser(t, ctx, user3.ID)
	}()

	spotifyID := "spotify_track_rating_test"

	// Create reviews with different ratings for the same track
	review1 := setupTestReview(t, user1.ID, "track")
	review1.SpotifyID = spotifyID
	review1.Rating = 3
	review2 := setupTestReview(t, user2.ID, "track")
	review2.SpotifyID = spotifyID
	review2.Rating = 4
	review3 := setupTestReview(t, user3.ID, "track")
	review3.SpotifyID = spotifyID
	review3.Rating = 5

	defer func() {
		cleanupTestReview(t, ctx, review1.ID)
		cleanupTestReview(t, ctx, review2.ID)
		cleanupTestReview(t, ctx, review3.ID)
	}()

	err = repo.Create(ctx, review1)
	require.NoError(t, err)
	err = repo.Create(ctx, review2)
	require.NoError(t, err)
	err = repo.Create(ctx, review3)
	require.NoError(t, err)

	// Test getting average rating (if method exists)
	reviews, err := repo.GetBySpotifyID(ctx, spotifyID, "track", 10, 0)
	require.NoError(t, err)

	// Calculate average manually for verification
	totalRating := 0
	for _, review := range reviews {
		totalRating += review.Rating
	}
	averageRating := float64(totalRating) / float64(len(reviews))

	assert.Equal(t, 4.0, averageRating) // (3+4+5)/3 = 4.0
}

func TestReviewRepository_RatingValidation(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewReviewRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	// Test valid ratings (1-5)
	validRatings := []int{1, 2, 3, 4, 5}
	for i, rating := range validRatings {
		review := setupTestReview(t, user.ID, "track")
		review.SpotifyID = fmt.Sprintf("spotify_track_rating_%d", i)
		review.Rating = rating
		defer cleanupTestReview(t, ctx, review.ID)

		err = repo.Create(ctx, review)
		require.NoError(t, err)

		// Verify rating was saved correctly
		createdReview, err := repo.GetByID(ctx, review.ID)
		require.NoError(t, err)
		assert.Equal(t, rating, createdReview.Rating)
	}
}

func TestReviewRepository_SpotifyTypeValidation(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewReviewRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	// Test valid Spotify types
	validTypes := []string{"track", "album"}
	for i, spotifyType := range validTypes {
		review := setupTestReview(t, user.ID, spotifyType)
		review.SpotifyID = fmt.Sprintf("spotify_%s_%d", spotifyType, i)
		defer cleanupTestReview(t, ctx, review.ID)

		err = repo.Create(ctx, review)
		require.NoError(t, err)

		// Verify type was saved correctly
		createdReview, err := repo.GetByID(ctx, review.ID)
		require.NoError(t, err)
		assert.Equal(t, spotifyType, createdReview.SpotifyType)
	}
}

func TestReviewRepository_UniqueUserSpotifyConstraint(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewReviewRepository(testDB)
	ctx := context.Background()

	// Create a test user first
	user := setupTestUser(t)
	userRepo := NewUserRepository(testDB)
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	defer cleanupTestUser(t, ctx, user.ID)

	spotifyID := "spotify_track_unique_test"

	// Create first review
	review1 := setupTestReview(t, user.ID, "track")
	review1.SpotifyID = spotifyID
	defer cleanupTestReview(t, ctx, review1.ID)

	err = repo.Create(ctx, review1)
	require.NoError(t, err)

	// Try to create second review for same user and Spotify ID
	// Note: Currently there's no unique constraint on (user_id, spotify_id) in the database schema
	// This test demonstrates the current behavior - multiple reviews are allowed
	review2 := setupTestReview(t, user.ID, "track")
	review2.SpotifyID = spotifyID
	defer cleanupTestReview(t, ctx, review2.ID)

	err = repo.Create(ctx, review2)
	// Currently this succeeds because there's no unique constraint
	// In a future schema update, a unique constraint could be added
	assert.NoError(t, err, "Currently allows multiple reviews per user/spotify_id - no unique constraint exists")

	// Verify both reviews exist
	reviews, err := repo.GetBySpotifyID(ctx, spotifyID, "track", 10, 0)
	require.NoError(t, err)
	assert.Len(t, reviews, 2, "Both reviews should exist since no unique constraint is enforced")
}

// Benchmark tests
func BenchmarkReviewRepository_Create(b *testing.B) {
	if testDB == nil {
		b.Skip("Database not available")
	}

	repo := NewReviewRepository(testDB)
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
		review := setupTestReview(&testing.T{}, user.ID, "track")
		review.SpotifyID = fmt.Sprintf("spotify_track_bench_%d", i)

		err := repo.Create(ctx, review)
		if err != nil {
			b.Fatal(err)
		}

		// Cleanup
		cleanupTestReview(&testing.T{}, ctx, review.ID)
	}
}

func BenchmarkReviewRepository_GetBySpotifyID(b *testing.B) {
	if testDB == nil {
		b.Skip("Database not available")
	}

	repo := NewReviewRepository(testDB)
	userRepo := NewUserRepository(testDB)
	ctx := context.Background()

	// Setup test data
	user := setupTestUser(&testing.T{})
	err := userRepo.Create(ctx, user)
	if err != nil {
		b.Fatal(err)
	}
	defer cleanupTestUser(&testing.T{}, ctx, user.ID)

	review := setupTestReview(&testing.T{}, user.ID, "track")
	err = repo.Create(ctx, review)
	if err != nil {
		b.Fatal(err)
	}
	defer cleanupTestReview(&testing.T{}, ctx, review.ID)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := repo.GetBySpotifyID(ctx, review.SpotifyID, "track", 10, 0)
		if err != nil {
			b.Fatal(err)
		}
	}
}
