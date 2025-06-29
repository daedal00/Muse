package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/daedal00/muse/backend/internal/config"
	"github.com/daedal00/muse/backend/internal/database"
	"github.com/daedal00/muse/backend/internal/spotify"
	spotifyapi "github.com/zmb3/spotify/v2"
)

// Album data structure matching what we'd cache
type CachedAlbum struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Artist      string `json:"artist"`
	ReleaseDate string `json:"release_date"`
	ImageURL    string `json:"image_url"`
	SpotifyID   string `json:"spotify_id"`
}

var (
	testRedisClient   *database.RedisClient
	testSpotifyClient *spotify.Client
	testSpotifyAPI    *spotifyapi.Client
)

func TestMain(m *testing.M) {
	// Setup
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Setup Redis
	if cfg.RedisURL != "" {
		testRedisClient, err = database.NewRedisConnection(cfg.RedisURL)
		if err != nil {
			fmt.Printf("Failed to connect to Redis: %v\n", err)
		}
	}

	// Setup Spotify client
	testSpotifyClient = spotify.NewClient(spotify.Config{
		ClientID:     cfg.SpotifyClientID,
		ClientSecret: cfg.SpotifyClientSecret,
	})

	// Get client credentials token for public API access
	ctx := context.Background()
	if testSpotifyClient != nil {
		testSpotifyAPI, err = testSpotifyClient.GetClientCredentialsClient(ctx)
		if err != nil {
			fmt.Printf("Failed to get Spotify client credentials: %v\n", err)
			testSpotifyAPI = nil
		}
	}

	// Run tests
	code := m.Run()

	// Cleanup
	if testRedisClient != nil {
		testRedisClient.Close()
	}

	os.Exit(code)
}

// Test real Spotify API vs Cache performance
func TestRealSpotifyVsCachePerformance(t *testing.T) {
	if testRedisClient == nil {
		t.Skip("Redis not available")
	}
	if testSpotifyAPI == nil {
		t.Skip("Spotify API not available")
	}

	ctx := context.Background()
	searchQuery := "Abbey Road The Beatles"
	cacheKey := fmt.Sprintf("search:album:%s", searchQuery)

	// First, make a real Spotify API call to get actual data
	t.Run("SpotifyAPI_FirstCall", func(t *testing.T) {
		start := time.Now()

		searchResult, err := testSpotifyAPI.Search(ctx, searchQuery, spotifyapi.SearchTypeAlbum)
		if err != nil {
			t.Fatalf("Spotify search failed: %v", err)
		}

		duration := time.Since(start)

		if searchResult.Albums == nil || len(searchResult.Albums.Albums) == 0 {
			t.Fatalf("No albums found for query: %s", searchQuery)
		}

		// Get the first album (most relevant)
		album := searchResult.Albums.Albums[0]

		// Convert to our cached format
		cachedAlbum := CachedAlbum{
			ID:          fmt.Sprintf("album_%s", album.ID),
			Name:        album.Name,
			Artist:      album.Artists[0].Name,
			ReleaseDate: album.ReleaseDate,
			SpotifyID:   string(album.ID),
		}

		if len(album.Images) > 0 {
			cachedAlbum.ImageURL = album.Images[0].URL
		}

		// Cache the result
		albumJSON, _ := json.Marshal(cachedAlbum)
		err = testRedisClient.Set(ctx, cacheKey, string(albumJSON), 24*time.Hour)
		if err != nil {
			t.Fatalf("Failed to cache result: %v", err)
		}

		t.Logf("Spotify API Call: Found '%s' by %s in %v",
			cachedAlbum.Name, cachedAlbum.Artist, duration)
		t.Logf("Spotify API Response Time: %v", duration)
	})

	// Now test cache retrieval
	t.Run("CacheHit_AfterAPICall", func(t *testing.T) {
		start := time.Now()

		cachedData, err := testRedisClient.Get(ctx, cacheKey)
		if err != nil {
			t.Fatalf("Cache retrieval failed: %v", err)
		}

		duration := time.Since(start)

		var cachedAlbum CachedAlbum
		err = json.Unmarshal([]byte(cachedData), &cachedAlbum)
		if err != nil {
			t.Fatalf("Failed to unmarshal cached data: %v", err)
		}

		t.Logf("Cache Hit: Retrieved '%s' by %s in %v",
			cachedAlbum.Name, cachedAlbum.Artist, duration)
		t.Logf("Cache Response Time: %v", duration)
	})

	// Test multiple cache hits to get average
	t.Run("CacheHit_Multiple", func(t *testing.T) {
		iterations := 100
		start := time.Now()

		for i := 0; i < iterations; i++ {
			_, err := testRedisClient.Get(ctx, cacheKey)
			if err != nil {
				t.Fatalf("Cache hit %d failed: %v", i, err)
			}
		}

		duration := time.Since(start)
		avgTime := duration / time.Duration(iterations)

		t.Logf("Cache Performance: %d hits took %v (avg: %v per hit)",
			iterations, duration, avgTime)
	})
}

// Benchmark comparing real Spotify API vs Cache
func BenchmarkRealSpotifyVsCache(b *testing.B) {
	if testRedisClient == nil || testSpotifyAPI == nil {
		b.Skip("Redis or Spotify API not available")
	}

	ctx := context.Background()
	searchQuery := "Thriller Michael Jackson"
	cacheKey := fmt.Sprintf("search:album:%s", searchQuery)

	// Pre-populate cache with real Spotify data
	searchResult, err := testSpotifyAPI.Search(ctx, searchQuery, spotifyapi.SearchTypeAlbum)
	if err != nil {
		b.Fatalf("Failed to get Spotify data for benchmark: %v", err)
	}

	if searchResult.Albums == nil || len(searchResult.Albums.Albums) == 0 {
		b.Fatalf("No albums found for benchmark query: %s", searchQuery)
	}

	album := searchResult.Albums.Albums[0]
	cachedAlbum := CachedAlbum{
		ID:          fmt.Sprintf("album_%s", album.ID),
		Name:        album.Name,
		Artist:      album.Artists[0].Name,
		ReleaseDate: album.ReleaseDate,
		SpotifyID:   string(album.ID),
	}

	if len(album.Images) > 0 {
		cachedAlbum.ImageURL = album.Images[0].URL
	}

	albumJSON, _ := json.Marshal(cachedAlbum)
	err = testRedisClient.Set(ctx, cacheKey, string(albumJSON), 24*time.Hour)
	if err != nil {
		b.Fatalf("Failed to cache album for benchmark: %v", err)
	}

	b.Run("RedisCache", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := testRedisClient.Get(ctx, cacheKey)
			if err != nil {
				b.Fatalf("Cache query failed: %v", err)
			}
		}
	})

	b.Run("SpotifyAPI", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Make real Spotify API calls (be careful with rate limits)
			_, err := testSpotifyAPI.Search(ctx, searchQuery, spotifyapi.SearchTypeAlbum)
			if err != nil {
				b.Fatalf("Spotify API call failed: %v", err)
			}
		}
	})
}

// Test realistic cache hit patterns with different albums
func TestAlbumCacheHitPatterns(t *testing.T) {
	if testRedisClient == nil || testSpotifyAPI == nil {
		t.Skip("Redis or Spotify API not available")
	}

	ctx := context.Background()

	// Popular albums that would be frequently searched
	popularSearches := []string{
		"Abbey Road The Beatles",
		"Thriller Michael Jackson",
		"Dark Side of the Moon Pink Floyd",
		"Nevermind Nirvana",
		"Back in Black AC/DC",
	}

	// Pre-populate cache with popular albums
	t.Logf("Pre-populating cache with popular albums...")
	for _, query := range popularSearches {
		cacheKey := fmt.Sprintf("search:album:%s", query)

		// Check if already cached
		_, err := testRedisClient.Get(ctx, cacheKey)
		if err != nil {
			// Not cached, fetch from Spotify
			searchResult, err := testSpotifyAPI.Search(ctx, query, spotifyapi.SearchTypeAlbum)
			if err != nil {
				t.Logf("Failed to search for %s: %v", query, err)
				continue
			}

			if searchResult.Albums != nil && len(searchResult.Albums.Albums) > 0 {
				album := searchResult.Albums.Albums[0]
				cachedAlbum := CachedAlbum{
					ID:          fmt.Sprintf("album_%s", album.ID),
					Name:        album.Name,
					Artist:      album.Artists[0].Name,
					ReleaseDate: album.ReleaseDate,
					SpotifyID:   string(album.ID),
				}

				if len(album.Images) > 0 {
					cachedAlbum.ImageURL = album.Images[0].URL
				}

				albumJSON, _ := json.Marshal(cachedAlbum)
				err = testRedisClient.Set(ctx, cacheKey, string(albumJSON), 24*time.Hour)
				if err != nil {
					t.Logf("Failed to cache album %s: %v", query, err)
					continue
				}
				t.Logf("Cached: %s by %s", cachedAlbum.Name, cachedAlbum.Artist)
			}
		}
	}

	// Simulate user search patterns (80% popular, 20% random)
	totalSearches := 100
	cacheHits := 0
	cacheMisses := 0

	for i := 0; i < totalSearches; i++ {
		var query string

		if i < 80 { // 80% search for popular albums
			query = popularSearches[i%len(popularSearches)]
		} else { // 20% search for random albums
			query = fmt.Sprintf("random album %d", i)
		}

		cacheKey := fmt.Sprintf("search:album:%s", query)
		_, err := testRedisClient.Get(ctx, cacheKey)

		if err != nil {
			cacheMisses++
		} else {
			cacheHits++
		}
	}

	hitRate := float64(cacheHits) / float64(totalSearches) * 100

	t.Logf("Album Search Cache Analysis:")
	t.Logf("Total Searches: %d", totalSearches)
	t.Logf("Cache Hits: %d", cacheHits)
	t.Logf("Cache Misses: %d", cacheMisses)
	t.Logf("Hit Rate: %.1f%%", hitRate)
}
