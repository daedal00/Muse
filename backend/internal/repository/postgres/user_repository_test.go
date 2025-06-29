package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/daedal00/muse/backend/internal/database"
	"github.com/daedal00/muse/backend/internal/models"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDB *database.PostgresDB

func TestMain(m *testing.M) {
	// Load environment variables for testing
	_ = godotenv.Load("../../../.env")

	// Setup test database
	var err error
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		testDBURL = os.Getenv("DATABASE_URL")
	}

	if testDBURL == "" {
		fmt.Println("No DATABASE_URL or TEST_DATABASE_URL set, skipping database tests")
		os.Exit(0)
	}

	fmt.Printf("Connecting to database...\n")
	testDB, err = database.NewPostgresConnection(testDBURL)
	if err != nil {
		fmt.Printf("Warning: Could not connect to test database (%s): %v\n", testDBURL[:20]+"...", err)
		fmt.Println("Skipping database tests. To run tests, ensure DATABASE_URL is set correctly")
		os.Exit(0)
	}

	// Test database connection
	if err := testDB.Health(); err != nil {
		fmt.Printf("Warning: Database health check failed: %v\n", err)
		testDB.Close()
		os.Exit(0)
	}

	// Test if users table exists and has expected columns (including new Spotify columns)
	ctx := context.Background()
	expectedColumns := []string{"bio", "spotify_id", "spotify_access_token", "spotify_refresh_token", "spotify_token_expiry"}

	for _, column := range expectedColumns {
		var columnExists bool
		err = testDB.Pool.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name=$1)", column).Scan(&columnExists)
		if err != nil {
			fmt.Printf("Warning: Could not check table structure: %v\n", err)
			testDB.Close()
			os.Exit(0)
		}

		if !columnExists {
			fmt.Printf("Warning: users table does not have '%s' column. Please run migrations first.\n", column)
			testDB.Close()
			os.Exit(0)
		}
	}

	fmt.Println("Database connection and schema validated successfully")

	// Run tests
	code := m.Run()

	// Cleanup
	testDB.Close()
	os.Exit(code)
}

// setupTestUser creates a test user for repository tests with new schema fields
func setupTestUser(t *testing.T) *models.User {
	t.Helper()

	return &models.User{
		ID:           uuid.New(),
		Name:         "Test User",
		Email:        fmt.Sprintf("test-%s@example.com", uuid.New().String()[:8]),
		PasswordHash: "hashed_password",
		Bio:          stringPtr("Test bio for user"),
		Avatar:       stringPtr("https://example.com/avatar.jpg"),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// setupTestUserWithSpotify creates a test user with Spotify OAuth fields
func setupTestUserWithSpotify(t *testing.T) *models.User {
	t.Helper()

	user := setupTestUser(t)
	user.SpotifyID = stringPtr("spotify_user_123")
	user.SpotifyAccessToken = stringPtr("access_token_123")
	user.SpotifyRefreshToken = stringPtr("refresh_token_123")
	user.SpotifyTokenExpiry = timePtr(time.Now().Add(time.Hour))

	return user
}

// cleanupTestUser removes a test user from the database
func cleanupTestUser(t *testing.T, ctx context.Context, userID uuid.UUID) {
	t.Helper()

	_, err := testDB.Pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		t.Logf("Warning: Failed to cleanup test user: %v", err)
	}
}

// cleanupBenchUser removes a test user from the database (for benchmarks)
func cleanupBenchUser(b *testing.B, ctx context.Context, userID uuid.UUID) {
	b.Helper()

	_, err := testDB.Pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		b.Logf("Warning: Failed to cleanup test user: %v", err)
	}
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestUserRepository_Create(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewUserRepository(testDB)
	ctx := context.Background()

	user := setupTestUser(t)
	defer cleanupTestUser(t, ctx, user.ID)

	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Verify user was created
	createdUser, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)

	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.PasswordHash, createdUser.PasswordHash)
	assert.Equal(t, *user.Bio, *createdUser.Bio)
	assert.Equal(t, *user.Avatar, *createdUser.Avatar)
	assert.Nil(t, createdUser.SpotifyID) // Should be nil for new user
}

func TestUserRepository_CreateWithSpotify(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewUserRepository(testDB)
	ctx := context.Background()

	user := setupTestUserWithSpotify(t)
	defer cleanupTestUser(t, ctx, user.ID)

	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Verify user was created with Spotify fields
	createdUser, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)

	assert.NotNil(t, createdUser.SpotifyID)
	assert.Equal(t, *user.SpotifyID, *createdUser.SpotifyID)
	assert.NotNil(t, createdUser.SpotifyAccessToken)
	assert.NotNil(t, createdUser.SpotifyRefreshToken)
	assert.NotNil(t, createdUser.SpotifyTokenExpiry)
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewUserRepository(testDB)
	ctx := context.Background()

	// Test getting non-existent user
	nonExistentID := uuid.New()
	_, err := repo.GetByID(ctx, nonExistentID)
	assert.Error(t, err)
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewUserRepository(testDB)
	ctx := context.Background()

	// Test getting non-existent user
	_, err := repo.GetByEmail(ctx, "nonexistent@example.com")
	assert.Error(t, err)
}

func TestUserRepository_GetByID(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewUserRepository(testDB)
	ctx := context.Background()

	user := setupTestUser(t)
	defer cleanupTestUser(t, ctx, user.ID)

	// Create user first
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Test getting existing user
	foundUser, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)

	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Name, foundUser.Name)
	assert.Equal(t, user.Email, foundUser.Email)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewUserRepository(testDB)
	ctx := context.Background()

	user := setupTestUser(t)
	defer cleanupTestUser(t, ctx, user.ID)

	// Create user first
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Test getting existing user by email
	foundUser, err := repo.GetByEmail(ctx, user.Email)
	require.NoError(t, err)

	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Email, foundUser.Email)
}

func TestUserRepository_GetBySpotifyID(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewUserRepository(testDB)
	ctx := context.Background()

	user := setupTestUserWithSpotify(t)
	defer cleanupTestUser(t, ctx, user.ID)

	// Create user first
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Test getting existing user by Spotify ID
	foundUser, err := repo.GetBySpotifyID(ctx, *user.SpotifyID)
	require.NoError(t, err)

	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, *user.SpotifyID, *foundUser.SpotifyID)
}

func TestUserRepository_Update(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewUserRepository(testDB)
	ctx := context.Background()

	user := setupTestUser(t)
	defer cleanupTestUser(t, ctx, user.ID)

	// Create user first
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Update user fields
	user.Name = "Updated User Name"
	user.Bio = stringPtr("Updated bio")
	user.Avatar = stringPtr("https://example.com/new-avatar.jpg")

	err = repo.Update(ctx, user)
	require.NoError(t, err)

	// Verify updates
	updatedUser, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)

	assert.Equal(t, "Updated User Name", updatedUser.Name)
	assert.Equal(t, "Updated bio", *updatedUser.Bio)
	assert.Equal(t, "https://example.com/new-avatar.jpg", *updatedUser.Avatar)
}

func TestUserRepository_UpdateSpotifyCredentials(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewUserRepository(testDB)
	ctx := context.Background()

	user := setupTestUser(t)
	defer cleanupTestUser(t, ctx, user.ID)

	// Create user first
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Update with Spotify credentials
	user.SpotifyID = stringPtr("new_spotify_id_123")
	user.SpotifyAccessToken = stringPtr("new_access_token")
	user.SpotifyRefreshToken = stringPtr("new_refresh_token")
	user.SpotifyTokenExpiry = timePtr(time.Now().Add(2 * time.Hour))

	err = repo.Update(ctx, user)
	require.NoError(t, err)

	// Verify Spotify fields were updated
	updatedUser, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)

	assert.NotNil(t, updatedUser.SpotifyID)
	assert.Equal(t, "new_spotify_id_123", *updatedUser.SpotifyID)
	assert.NotNil(t, updatedUser.SpotifyAccessToken)
	assert.NotNil(t, updatedUser.SpotifyRefreshToken)
	assert.NotNil(t, updatedUser.SpotifyTokenExpiry)
}

func TestUserRepository_Delete(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewUserRepository(testDB)
	ctx := context.Background()

	user := setupTestUser(t)

	// Create user first
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Delete user
	err = repo.Delete(ctx, user.ID)
	require.NoError(t, err)

	// Verify user was deleted
	_, err = repo.GetByID(ctx, user.ID)
	assert.Error(t, err)
}

func TestUserRepository_List(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewUserRepository(testDB)
	ctx := context.Background()

	// Create multiple test users
	user1 := setupTestUser(t)
	user2 := setupTestUser(t)
	user3 := setupTestUser(t)

	defer func() {
		cleanupTestUser(t, ctx, user1.ID)
		cleanupTestUser(t, ctx, user2.ID)
		cleanupTestUser(t, ctx, user3.ID)
	}()

	err := repo.Create(ctx, user1)
	require.NoError(t, err)
	err = repo.Create(ctx, user2)
	require.NoError(t, err)
	err = repo.Create(ctx, user3)
	require.NoError(t, err)

	// Test listing users
	users, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)

	// Should contain at least our 3 test users
	assert.GreaterOrEqual(t, len(users), 3)

	// Check if our users are in the list
	userIDs := make(map[uuid.UUID]bool)
	for _, user := range users {
		userIDs[user.ID] = true
	}

	assert.True(t, userIDs[user1.ID])
	assert.True(t, userIDs[user2.ID])
	assert.True(t, userIDs[user3.ID])
}

func TestUserRepository_ListPagination(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewUserRepository(testDB)
	ctx := context.Background()

	// Test with small limit and offset
	users, err := repo.List(ctx, 2, 0)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(users), 2)

	// Test with offset
	usersWithOffset, err := repo.List(ctx, 2, 1)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(usersWithOffset), 2)

	// First user from offset query should be different (if we have enough users)
	if len(users) > 0 && len(usersWithOffset) > 0 {
		assert.NotEqual(t, users[0].ID, usersWithOffset[0].ID)
	}
}

func TestUserRepository_EmailUniqueConstraint(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	repo := NewUserRepository(testDB)
	ctx := context.Background()

	user1 := setupTestUser(t)
	user2 := setupTestUser(t)
	user2.Email = user1.Email // Same email

	defer func() {
		cleanupTestUser(t, ctx, user1.ID)
		cleanupTestUser(t, ctx, user2.ID)
	}()

	// Create first user
	err := repo.Create(ctx, user1)
	require.NoError(t, err)

	// Try to create second user with same email - should fail
	err = repo.Create(ctx, user2)
	assert.Error(t, err, "Should fail due to unique email constraint")
}

// Benchmark tests
func BenchmarkUserRepository_Create(b *testing.B) {
	if testDB == nil {
		b.Skip("Database not available")
	}

	repo := NewUserRepository(testDB)
	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		user := setupTestUser(&testing.T{}) // Use empty testing.T for benchmark
		user.Email = fmt.Sprintf("bench-%d@example.com", i)

		err := repo.Create(ctx, user)
		if err != nil {
			b.Fatal(err)
		}

		// Cleanup
		cleanupBenchUser(b, ctx, user.ID)
	}
}

func BenchmarkUserRepository_GetByID(b *testing.B) {
	if testDB == nil {
		b.Skip("Database not available")
	}

	repo := NewUserRepository(testDB)
	ctx := context.Background()

	// Setup test user
	user := setupTestUser(&testing.T{})
	err := repo.Create(ctx, user)
	if err != nil {
		b.Fatal(err)
	}
	defer cleanupBenchUser(b, ctx, user.ID)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := repo.GetByID(ctx, user.ID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUserRepository_GetByEmail(b *testing.B) {
	if testDB == nil {
		b.Skip("Database not available")
	}

	repo := NewUserRepository(testDB)
	ctx := context.Background()

	// Setup test user
	user := setupTestUser(&testing.T{})
	err := repo.Create(ctx, user)
	if err != nil {
		b.Fatal(err)
	}
	defer cleanupBenchUser(b, ctx, user.ID)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := repo.GetByEmail(ctx, user.Email)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUserRepository_List(b *testing.B) {
	if testDB == nil {
		b.Skip("Database not available")
	}

	repo := NewUserRepository(testDB)
	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := repo.List(ctx, 10, 0)
		if err != nil {
			b.Fatal(err)
		}
	}
}
