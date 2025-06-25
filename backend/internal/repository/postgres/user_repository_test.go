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
	
	// Test if users table exists and has expected columns
	ctx := context.Background()
	var columnExists bool
	err = testDB.Pool.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='bio')").Scan(&columnExists)
	if err != nil {
		fmt.Printf("Warning: Could not check table structure: %v\n", err)
		testDB.Close()
		os.Exit(0)
	}
	
	if !columnExists {
		fmt.Println("Warning: users table does not have 'bio' column. Please run migrations first.")
		testDB.Close()
		os.Exit(0)
	}
	
	fmt.Println("Database connection and schema validated successfully")
	
	// Run tests
	code := m.Run()
	
	// Cleanup
	testDB.Close()
	os.Exit(code)
}

// setupTestUser creates a test user for repository tests
func setupTestUser(t *testing.T) *models.User {
	t.Helper()
	
	return &models.User{
		ID:           uuid.New(),
		Name:         "Test User",
		Email:        fmt.Sprintf("test-%s@example.com", uuid.New().String()[:8]),
		PasswordHash: "hashed_password",
		Bio:          stringPtr("Test bio"),
		Avatar:       stringPtr("https://example.com/avatar.jpg"),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// cleanupTestUser removes a test user from the database
func cleanupTestUser(t *testing.T, ctx context.Context, userID uuid.UUID) {
	t.Helper()
	
	_, err := testDB.Pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		t.Logf("Warning: Failed to cleanup test user: %v", err)
	}
}

func stringPtr(s string) *string {
	return &s
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
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Verify user was created
	createdUser, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get created user: %v", err)
	}
	
	if createdUser.Name != user.Name {
		t.Errorf("Expected name %s, got %s", user.Name, createdUser.Name)
	}
	if createdUser.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, createdUser.Email)
	}
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
	if err == nil {
		t.Error("Expected error when getting non-existent user")
	}
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}
	
	repo := NewUserRepository(testDB)
	ctx := context.Background()
	
	// Test getting non-existent user
	_, err := repo.GetByEmail(ctx, "nonexistent@example.com")
	if err == nil {
		t.Error("Expected error when getting non-existent user")
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	repo := NewUserRepository(testDB)
	ctx := context.Background()
	
	user := setupTestUser(t)
	defer cleanupTestUser(t, ctx, user.ID)
	
	// Create user first
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Test getting existing user
	foundUser, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}
	
	if foundUser.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, foundUser.ID)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	repo := NewUserRepository(testDB)
	ctx := context.Background()
	
	user := setupTestUser(t)
	user.Email = "unique@example.com" // Use unique email
	defer cleanupTestUser(t, ctx, user.ID)
	
	// Create user first
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Test getting existing user by email
	foundUser, err := repo.GetByEmail(ctx, user.Email)
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}
	
	if foundUser.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, foundUser.Email)
	}
}

func TestUserRepository_Update(t *testing.T) {
	repo := NewUserRepository(testDB)
	ctx := context.Background()
	
	user := setupTestUser(t)
	defer cleanupTestUser(t, ctx, user.ID)
	
	// Create user first
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Update user
	user.Name = "Updated Name"
	user.Bio = stringPtr("Updated bio")
	
	err = repo.Update(ctx, user)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}
	
	// Verify update
	updatedUser, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}
	
	if updatedUser.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got %s", updatedUser.Name)
	}
	if *updatedUser.Bio != "Updated bio" {
		t.Errorf("Expected bio 'Updated bio', got %s", *updatedUser.Bio)
	}
	
	// Test updating non-existent user
	nonExistentUser := setupTestUser(t)
	err = repo.Update(ctx, nonExistentUser)
	if err == nil {
		t.Error("Expected error when updating non-existent user")
	}
}

func TestUserRepository_Delete(t *testing.T) {
	repo := NewUserRepository(testDB)
	ctx := context.Background()
	
	user := setupTestUser(t)
	
	// Create user first
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Delete user
	err = repo.Delete(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}
	
	// Verify deletion
	_, err = repo.GetByID(ctx, user.ID)
	if err == nil {
		t.Error("Expected error when getting deleted user")
	}
	
	// Test deleting non-existent user
	err = repo.Delete(ctx, uuid.New())
	if err == nil {
		t.Error("Expected error when deleting non-existent user")
	}
}

func TestUserRepository_List(t *testing.T) {
	repo := NewUserRepository(testDB)
	ctx := context.Background()
	
	// Create multiple test users
	var userIDs []uuid.UUID
	for i := 0; i < 3; i++ {
		user := setupTestUser(t)
		user.Email = fmt.Sprintf("test%d@example.com", i)
		userIDs = append(userIDs, user.ID)
		
		err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user %d: %v", i, err)
		}
	}
	
	// Cleanup
	defer func() {
		for _, id := range userIDs {
			cleanupTestUser(t, ctx, id)
		}
	}()
	
	// Test listing users
	users, err := repo.List(ctx, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}
	
	if len(users) < 3 {
		t.Errorf("Expected at least 3 users, got %d", len(users))
	}
	
	// Test with limit
	users, err = repo.List(ctx, 2, 0)
	if err != nil {
		t.Fatalf("Failed to list users with limit: %v", err)
	}
	
	if len(users) > 2 {
		t.Errorf("Expected at most 2 users, got %d", len(users))
	}
} 