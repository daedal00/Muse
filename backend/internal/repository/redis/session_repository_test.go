package redis

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

var testRedis *database.RedisClient

func TestMain(m *testing.M) {
	// Load environment variables for testing
	_ = godotenv.Load("../../../.env")

	// Setup test Redis
	var err error
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/1" // Use DB 1 for tests
	}

	fmt.Printf("Connecting to Redis at %s...\n", redisURL)
	testRedis, err = database.NewRedisConnection(redisURL)
	if err != nil {
		fmt.Printf("Warning: Could not connect to Redis: %v\n", err)
		fmt.Println("Some tests will be skipped")
	}

	// Run tests
	code := m.Run()

	// Cleanup
	if testRedis != nil {
		ctx := context.Background()
		// Clean up test data
		testRedis.Client.FlushDB(ctx)
		testRedis.Close()
	}

	os.Exit(code)
}

func TestNewSessionRepository(t *testing.T) {
	if testRedis == nil {
		t.Skip("Redis not available")
	}

	repo := NewSessionRepository(testRedis)
	assert.NotNil(t, repo)
}

func TestSessionRepository_Create(t *testing.T) {
	if testRedis == nil {
		t.Skip("Redis not available")
	}

	repo := NewSessionRepository(testRedis)
	ctx := context.Background()

	// Clean up before test
	testRedis.Client.FlushDB(ctx)

	userID := uuid.New()
	session := &models.Session{
		ID:        "test-session-123",
		UserID:    userID,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		CreatedAt: time.Now(),
	}

	err := repo.Create(ctx, session)
	require.NoError(t, err)

	// Verify session was created
	retrievedSession, err := repo.GetByID(ctx, session.ID)
	require.NoError(t, err)
	assert.Equal(t, session.ID, retrievedSession.ID)
	assert.Equal(t, session.UserID, retrievedSession.UserID)
	assert.WithinDuration(t, session.ExpiresAt, retrievedSession.ExpiresAt, time.Second)
}

func TestSessionRepository_Create_ExpiredSession(t *testing.T) {
	if testRedis == nil {
		t.Skip("Redis not available")
	}

	repo := NewSessionRepository(testRedis)
	ctx := context.Background()

	userID := uuid.New()
	session := &models.Session{
		ID:        "expired-session-123",
		UserID:    userID,
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Already expired
		CreatedAt: time.Now(),
	}

	err := repo.Create(ctx, session)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already expired")
}

func TestSessionRepository_GetByID(t *testing.T) {
	if testRedis == nil {
		t.Skip("Redis not available")
	}

	repo := NewSessionRepository(testRedis)
	ctx := context.Background()

	// Clean up before test
	testRedis.Client.FlushDB(ctx)

	userID := uuid.New()
	session := &models.Session{
		ID:        "get-test-session-123",
		UserID:    userID,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		CreatedAt: time.Now(),
	}

	// Create session first
	err := repo.Create(ctx, session)
	require.NoError(t, err)

	// Test getting existing session
	retrievedSession, err := repo.GetByID(ctx, session.ID)
	require.NoError(t, err)
	assert.Equal(t, session.ID, retrievedSession.ID)
	assert.Equal(t, session.UserID, retrievedSession.UserID)

	// Test getting non-existent session
	_, err = repo.GetByID(ctx, "non-existent-session")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestSessionRepository_GetByUserID(t *testing.T) {
	if testRedis == nil {
		t.Skip("Redis not available")
	}

	repo := NewSessionRepository(testRedis)
	ctx := context.Background()

	// Clean up before test
	testRedis.Client.FlushDB(ctx)

	userID := uuid.New()

	// Create multiple sessions for the same user
	sessions := []*models.Session{
		{
			ID:        "user-session-1",
			UserID:    userID,
			ExpiresAt: time.Now().Add(1 * time.Hour),
			CreatedAt: time.Now(),
		},
		{
			ID:        "user-session-2",
			UserID:    userID,
			ExpiresAt: time.Now().Add(2 * time.Hour),
			CreatedAt: time.Now(),
		},
	}

	// Create sessions
	for _, session := range sessions {
		err := repo.Create(ctx, session)
		require.NoError(t, err)
	}

	// Get sessions by user ID
	retrievedSessions, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, retrievedSessions, 2)

	// Verify session IDs are present
	sessionIDs := make(map[string]bool)
	for _, session := range retrievedSessions {
		sessionIDs[session.ID] = true
		assert.Equal(t, userID, session.UserID)
	}
	assert.True(t, sessionIDs["user-session-1"])
	assert.True(t, sessionIDs["user-session-2"])

	// Test getting sessions for non-existent user
	nonExistentUserID := uuid.New()
	emptySessions, err := repo.GetByUserID(ctx, nonExistentUserID)
	require.NoError(t, err)
	assert.Empty(t, emptySessions)
}

func TestSessionRepository_Delete(t *testing.T) {
	if testRedis == nil {
		t.Skip("Redis not available")
	}

	repo := NewSessionRepository(testRedis)
	ctx := context.Background()

	// Clean up before test
	testRedis.Client.FlushDB(ctx)

	userID := uuid.New()
	session := &models.Session{
		ID:        "delete-test-session",
		UserID:    userID,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		CreatedAt: time.Now(),
	}

	// Create session first
	err := repo.Create(ctx, session)
	require.NoError(t, err)

	// Verify session exists
	_, err = repo.GetByID(ctx, session.ID)
	require.NoError(t, err)

	// Delete session
	err = repo.Delete(ctx, session.ID)
	require.NoError(t, err)

	// Verify session is deleted
	_, err = repo.GetByID(ctx, session.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// Test deleting non-existent session
	err = repo.Delete(ctx, "non-existent-session")
	assert.Error(t, err)
}

func TestSessionRepository_DeleteByUserID(t *testing.T) {
	if testRedis == nil {
		t.Skip("Redis not available")
	}

	repo := NewSessionRepository(testRedis)
	ctx := context.Background()

	// Clean up before test
	testRedis.Client.FlushDB(ctx)

	userID := uuid.New()

	// Create multiple sessions for the user
	sessions := []*models.Session{
		{
			ID:        "delete-user-session-1",
			UserID:    userID,
			ExpiresAt: time.Now().Add(1 * time.Hour),
			CreatedAt: time.Now(),
		},
		{
			ID:        "delete-user-session-2",
			UserID:    userID,
			ExpiresAt: time.Now().Add(2 * time.Hour),
			CreatedAt: time.Now(),
		},
	}

	// Create sessions
	for _, session := range sessions {
		err := repo.Create(ctx, session)
		require.NoError(t, err)
	}

	// Verify sessions exist
	retrievedSessions, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, retrievedSessions, 2)

	// Delete all sessions for user
	err = repo.DeleteByUserID(ctx, userID)
	require.NoError(t, err)

	// Verify all sessions are deleted
	retrievedSessions, err = repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Empty(t, retrievedSessions)
}

func TestSessionRepository_DeleteExpired(t *testing.T) {
	if testRedis == nil {
		t.Skip("Redis not available")
	}

	repo := NewSessionRepository(testRedis)
	ctx := context.Background()

	// Clean up before test
	testRedis.Client.FlushDB(ctx)

	userID := uuid.New()

	// Create a session that will expire quickly
	shortSession := &models.Session{
		ID:        "short-lived-session",
		UserID:    userID,
		ExpiresAt: time.Now().Add(100 * time.Millisecond), // Very short expiration
		CreatedAt: time.Now(),
	}

	// Create a session with longer expiration
	longSession := &models.Session{
		ID:        "long-lived-session",
		UserID:    userID,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		CreatedAt: time.Now(),
	}

	// Create both sessions
	err := repo.Create(ctx, shortSession)
	require.NoError(t, err)
	err = repo.Create(ctx, longSession)
	require.NoError(t, err)

	// Wait for short session to expire
	time.Sleep(200 * time.Millisecond)

	// Run cleanup
	err = repo.DeleteExpired(ctx)
	require.NoError(t, err)

	// Verify short session is gone (expired naturally by Redis)
	_, err = repo.GetByID(ctx, shortSession.ID)
	assert.Error(t, err)

	// Verify long session still exists
	_, err = repo.GetByID(ctx, longSession.ID)
	require.NoError(t, err)
}

func TestSessionRepository_ConcurrentOperations(t *testing.T) {
	if testRedis == nil {
		t.Skip("Redis not available")
	}

	repo := NewSessionRepository(testRedis)
	ctx := context.Background()

	// Clean up before test
	testRedis.Client.FlushDB(ctx)

	userID := uuid.New()

	// Test concurrent session creation
	numGoroutines := 10
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			session := &models.Session{
				ID:        fmt.Sprintf("concurrent-session-%d", index),
				UserID:    userID,
				ExpiresAt: time.Now().Add(1 * time.Hour),
				CreatedAt: time.Now(),
			}

			err := repo.Create(ctx, session)
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify all sessions were created
	sessions, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, sessions, numGoroutines)
}

// Benchmark tests for performance monitoring
func BenchmarkSessionRepository_Create(b *testing.B) {
	if testRedis == nil {
		b.Skip("Redis not available")
	}

	repo := NewSessionRepository(testRedis)
	ctx := context.Background()

	// Clean up before benchmark
	testRedis.Client.FlushDB(ctx)

	userID := uuid.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		session := &models.Session{
			ID:        fmt.Sprintf("benchmark-session-%d", i),
			UserID:    userID,
			ExpiresAt: time.Now().Add(1 * time.Hour),
			CreatedAt: time.Now(),
		}

		err := repo.Create(ctx, session)
		if err != nil {
			b.Fatalf("Failed to create session: %v", err)
		}
	}
}

func BenchmarkSessionRepository_GetByID(b *testing.B) {
	if testRedis == nil {
		b.Skip("Redis not available")
	}

	repo := NewSessionRepository(testRedis)
	ctx := context.Background()

	// Clean up and setup
	testRedis.Client.FlushDB(ctx)

	userID := uuid.New()
	session := &models.Session{
		ID:        "benchmark-get-session",
		UserID:    userID,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		CreatedAt: time.Now(),
	}

	err := repo.Create(ctx, session)
	if err != nil {
		b.Fatalf("Failed to create session: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.GetByID(ctx, session.ID)
		if err != nil {
			b.Fatalf("Failed to get session: %v", err)
		}
	}
}

func BenchmarkSessionRepository_GetByUserID(b *testing.B) {
	if testRedis == nil {
		b.Skip("Redis not available")
	}

	repo := NewSessionRepository(testRedis)
	ctx := context.Background()

	// Clean up and setup
	testRedis.Client.FlushDB(ctx)

	userID := uuid.New()

	// Create multiple sessions for benchmarking
	for i := 0; i < 10; i++ {
		session := &models.Session{
			ID:        fmt.Sprintf("benchmark-user-session-%d", i),
			UserID:    userID,
			ExpiresAt: time.Now().Add(1 * time.Hour),
			CreatedAt: time.Now(),
		}

		err := repo.Create(ctx, session)
		if err != nil {
			b.Fatalf("Failed to create session: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.GetByUserID(ctx, userID)
		if err != nil {
			b.Fatalf("Failed to get user sessions: %v", err)
		}
	}
}

func BenchmarkSessionRepository_Delete(b *testing.B) {
	if testRedis == nil {
		b.Skip("Redis not available")
	}

	repo := NewSessionRepository(testRedis)
	ctx := context.Background()

	userID := uuid.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Setup: Create session to delete
		session := &models.Session{
			ID:        fmt.Sprintf("benchmark-delete-session-%d", i),
			UserID:    userID,
			ExpiresAt: time.Now().Add(1 * time.Hour),
			CreatedAt: time.Now(),
		}

		err := repo.Create(ctx, session)
		if err != nil {
			b.Fatalf("Failed to create session: %v", err)
		}
		b.StartTimer()

		// Benchmark: Delete session
		err = repo.Delete(ctx, session.ID)
		if err != nil {
			b.Fatalf("Failed to delete session: %v", err)
		}
	}
}

func BenchmarkSessionRepository_DeleteByUserID(b *testing.B) {
	if testRedis == nil {
		b.Skip("Redis not available")
	}

	repo := NewSessionRepository(testRedis)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Setup: Create multiple sessions for user
		userID := uuid.New()
		for j := 0; j < 5; j++ {
			session := &models.Session{
				ID:        fmt.Sprintf("benchmark-user-delete-session-%d-%d", i, j),
				UserID:    userID,
				ExpiresAt: time.Now().Add(1 * time.Hour),
				CreatedAt: time.Now(),
			}

			err := repo.Create(ctx, session)
			if err != nil {
				b.Fatalf("Failed to create session: %v", err)
			}
		}
		b.StartTimer()

		// Benchmark: Delete all user sessions
		err := repo.DeleteByUserID(ctx, userID)
		if err != nil {
			b.Fatalf("Failed to delete user sessions: %v", err)
		}
	}
}
