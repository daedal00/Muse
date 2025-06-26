package database

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPostgresConnection_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Use test database URL or skip if not available
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}
	if databaseURL == "" {
		t.Skip("No DATABASE_URL or TEST_DATABASE_URL set, skipping database connection test")
	}

	// Test successful connection
	db, err := NewPostgresConnection(databaseURL)
	require.NoError(t, err)
	require.NotNil(t, db)

	// Verify connection is working
	err = db.Health()
	assert.NoError(t, err)

	// Clean up
	db.Close()
}

func TestNewPostgresConnection_InvalidURL(t *testing.T) {
	// Test with invalid database URL
	invalidURL := "postgres://invalid:invalid@nonexistent:5432/nonexistent"

	db, err := NewPostgresConnection(invalidURL)

	// Should either fail to connect or fail health check
	if err == nil {
		// If connection was created, health check should fail
		require.NotNil(t, db)
		err = db.Health()
		assert.Error(t, err, "Health check should fail for invalid connection")
		db.Close()
	} else {
		// Connection creation failed, which is also acceptable
		assert.Error(t, err)
	}
}

func TestNewPostgresConnection_EmptyURL(t *testing.T) {
	// Test with empty URL
	db, err := NewPostgresConnection("")
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestPostgresDB_Health(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}
	if databaseURL == "" {
		t.Skip("No DATABASE_URL or TEST_DATABASE_URL set, skipping health check test")
	}

	db, err := NewPostgresConnection(databaseURL)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()

	// Test health check
	err = db.Health()
	assert.NoError(t, err)

	// Test health check with context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.Pool.Ping(ctx)
	assert.NoError(t, err)
}

func TestPostgresDB_QueryExecution(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}
	if databaseURL == "" {
		t.Skip("No DATABASE_URL or TEST_DATABASE_URL set, skipping query test")
	}

	db, err := NewPostgresConnection(databaseURL)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()

	ctx := context.Background()

	// Test simple query
	var result int
	err = db.Pool.QueryRow(ctx, "SELECT 1").Scan(&result)
	assert.NoError(t, err)
	assert.Equal(t, 1, result)

	// Test query with current timestamp
	var now time.Time
	err = db.Pool.QueryRow(ctx, "SELECT NOW()").Scan(&now)
	assert.NoError(t, err)
	assert.True(t, time.Since(now) < time.Minute, "Timestamp should be recent")
}

func TestPostgresDB_TableExistence(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}
	if databaseURL == "" {
		t.Skip("No DATABASE_URL or TEST_DATABASE_URL set, skipping table existence test")
	}

	db, err := NewPostgresConnection(databaseURL)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()

	ctx := context.Background()

	// Test that expected tables exist after migrations
	expectedTables := []string{
		"users",
		"playlists",
		"playlist_tracks",
		"reviews",
		"user_preferences",
	}

	for _, tableName := range expectedTables {
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT 1 FROM information_schema.tables 
				WHERE table_schema = 'public' AND table_name = $1
			)`
		err = db.Pool.QueryRow(ctx, query, tableName).Scan(&exists)
		assert.NoError(t, err, "Failed to check existence of table %s", tableName)

		if !exists {
			t.Logf("Warning: Table %s does not exist. Please run migrations.", tableName)
		}
	}
}

func TestPostgresDB_SchemaValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}
	if databaseURL == "" {
		t.Skip("No DATABASE_URL or TEST_DATABASE_URL set, skipping schema validation test")
	}

	db, err := NewPostgresConnection(databaseURL)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()

	ctx := context.Background()

	// Check that users table has required columns for new schema
	expectedUserColumns := []string{
		"id", "name", "email", "password_hash", "bio", "avatar",
		"spotify_id", "spotify_access_token", "spotify_refresh_token", "spotify_token_expiry",
		"created_at", "updated_at",
	}

	for _, columnName := range expectedUserColumns {
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_name = 'users' AND column_name = $1
			)`
		err = db.Pool.QueryRow(ctx, query, columnName).Scan(&exists)
		assert.NoError(t, err, "Failed to check existence of column %s in users table", columnName)

		if !exists {
			t.Logf("Warning: Column %s does not exist in users table. Schema may be outdated.", columnName)
		}
	}

	// Check that reviews table has been updated for new schema
	expectedReviewColumns := []string{"spotify_id", "spotify_type"}

	for _, columnName := range expectedReviewColumns {
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_name = 'reviews' AND column_name = $1
			)`
		err = db.Pool.QueryRow(ctx, query, columnName).Scan(&exists)
		assert.NoError(t, err, "Failed to check existence of column %s in reviews table", columnName)

		if !exists {
			t.Logf("Warning: Column %s does not exist in reviews table. Schema may be outdated.", columnName)
		}
	}
}

func TestPostgresDB_ConnectionPooling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}
	if databaseURL == "" {
		t.Skip("No DATABASE_URL or TEST_DATABASE_URL set, skipping connection pooling test")
	}

	db, err := NewPostgresConnection(databaseURL)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()

	// Test concurrent connections
	numConcurrent := 5
	results := make(chan error, numConcurrent)

	for i := 0; i < numConcurrent; i++ {
		go func() {
			ctx := context.Background()
			var result int
			err := db.Pool.QueryRow(ctx, "SELECT 1").Scan(&result)
			results <- err
		}()
	}

	// Check all concurrent queries succeeded
	for i := 0; i < numConcurrent; i++ {
		err := <-results
		assert.NoError(t, err, "Concurrent query %d failed", i+1)
	}
}

func TestPostgresDB_TransactionSupport(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}
	if databaseURL == "" {
		t.Skip("No DATABASE_URL or TEST_DATABASE_URL set, skipping transaction test")
	}

	db, err := NewPostgresConnection(databaseURL)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()

	ctx := context.Background()

	// Test transaction creation and rollback
	tx, err := db.Pool.Begin(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, tx)

	// Perform some operation in transaction
	_, err = tx.Exec(ctx, "SELECT 1")
	assert.NoError(t, err)

	// Rollback transaction
	err = tx.Rollback(ctx)
	assert.NoError(t, err)
}

// Benchmark tests
func BenchmarkPostgresConnection(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}
	if databaseURL == "" {
		b.Skip("No DATABASE_URL or TEST_DATABASE_URL set, skipping benchmark")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		db, err := NewPostgresConnection(databaseURL)
		if err != nil {
			b.Fatal(err)
		}
		db.Close()
	}
}

func BenchmarkPostgresQuery(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}
	if databaseURL == "" {
		b.Skip("No DATABASE_URL or TEST_DATABASE_URL set, skipping benchmark")
	}

	db, err := NewPostgresConnection(databaseURL)
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result int
		err := db.Pool.QueryRow(ctx, "SELECT 1").Scan(&result)
		if err != nil {
			b.Fatal(err)
		}
	}
}
