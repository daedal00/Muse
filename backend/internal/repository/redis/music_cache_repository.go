package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/daedal00/muse/backend/internal/database"
	"github.com/daedal00/muse/backend/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// MusicCacheRepository handles caching of user music data in Redis
type MusicCacheRepository struct {
	client *database.RedisClient
}

// MusicData represents cached music data for a user
type MusicData struct {
	RecentlyPlayed []*models.Track `json:"recently_played"`
	FavoriteAlbums []*models.Album `json:"favorite_albums"`
	LastUpdated    time.Time       `json:"last_updated"`
}

// SearchCacheData represents cached search results
type SearchCacheData struct {
	Query      string        `json:"query"`
	Results    interface{}   `json:"results"` // Can be albums, tracks, or artists
	Timestamp  time.Time     `json:"timestamp"`
	ResultType string        `json:"result_type"` // "albums", "tracks", "artists"
}

// ListeningHistory represents a user's listening activity
type ListeningHistory struct {
	UserID    uuid.UUID              `json:"user_id"`
	Tracks    []*models.Track        `json:"tracks"`
	Albums    []*models.Album        `json:"albums"`
	Artists   []*models.Artist       `json:"artists"`
	Timestamp time.Time              `json:"timestamp"`
}

// Cache durations
const (
	MusicDataCacheTTL    = 1 * time.Hour   // User music data cache for 1 hour
	SearchCacheTTL       = 30 * time.Minute // Search results cache for 30 minutes
	HistoryCacheTTL      = 24 * time.Hour   // Listening history cache for 24 hours
	PopularDataCacheTTL  = 6 * time.Hour    // Popular content cache for 6 hours
)

func NewMusicCacheRepository(client *database.RedisClient) *MusicCacheRepository {
	return &MusicCacheRepository{client: client}
}

// ============ User Music Data Caching ============

// SetUserMusicData caches a user's music data (recently played, favorites, etc.)
func (r *MusicCacheRepository) SetUserMusicData(ctx context.Context, userID uuid.UUID, data interface{}) error {
	key := fmt.Sprintf("user_music:%s", userID.String())
	
	// Cast the data to MusicData
	musicData, ok := data.(*MusicData)
	if !ok {
		return fmt.Errorf("invalid data type for user music data")
	}
	
	// Set last updated timestamp
	musicData.LastUpdated = time.Now()
	
	jsonData, err := json.Marshal(musicData)
	if err != nil {
		return fmt.Errorf("failed to marshal user music data: %w", err)
	}
	
	return r.client.Client.Set(ctx, key, jsonData, MusicDataCacheTTL).Err()
}

// GetUserMusicData retrieves cached user music data
func (r *MusicCacheRepository) GetUserMusicData(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	key := fmt.Sprintf("user_music:%s", userID.String())
	
	data, err := r.client.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get user music data: %w", err)
	}
	
	var musicData MusicData
	if err := json.Unmarshal([]byte(data), &musicData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user music data: %w", err)
	}
	
	return &musicData, nil
}

// AddToRecentlyPlayed adds a track to user's recently played list
func (r *MusicCacheRepository) AddToRecentlyPlayed(ctx context.Context, userID uuid.UUID, track *models.Track) error {
	// Get existing data
	data, err := r.GetUserMusicData(ctx, userID)
	if err != nil {
		return err
	}
	
	var musicData *MusicData
	if data == nil {
		musicData = &MusicData{
			RecentlyPlayed: []*models.Track{},
			FavoriteAlbums: []*models.Album{},
		}
	} else {
		// Cast the interface{} to *MusicData
		var ok bool
		musicData, ok = data.(*MusicData)
		if !ok {
			return fmt.Errorf("invalid data type in cache")
		}
	}
	
	// Add track to beginning of recently played (most recent first)
	musicData.RecentlyPlayed = append([]*models.Track{track}, musicData.RecentlyPlayed...)
	
	// Keep only last 50 tracks
	if len(musicData.RecentlyPlayed) > 50 {
		musicData.RecentlyPlayed = musicData.RecentlyPlayed[:50]
	}
	
	return r.SetUserMusicData(ctx, userID, musicData)
}

// ============ Search Results Caching ============

// SetSearchResults caches search results for faster retrieval
func (r *MusicCacheRepository) SetSearchResults(ctx context.Context, query string, resultType string, results interface{}) error {
	key := fmt.Sprintf("search:%s:%s", resultType, query)
	
	cacheData := SearchCacheData{
		Query:      query,
		Results:    results,
		Timestamp:  time.Now(),
		ResultType: resultType,
	}
	
	jsonData, err := json.Marshal(cacheData)
	if err != nil {
		return fmt.Errorf("failed to marshal search results: %w", err)
	}
	
	return r.client.Client.Set(ctx, key, jsonData, SearchCacheTTL).Err()
}

// GetSearchResults retrieves cached search results
func (r *MusicCacheRepository) GetSearchResults(ctx context.Context, query string, resultType string) (interface{}, error) {
	key := fmt.Sprintf("search:%s:%s", resultType, query)
	
	data, err := r.client.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get search results: %w", err)
	}
	
	var searchData SearchCacheData
	if err := json.Unmarshal([]byte(data), &searchData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}
	
	return &searchData, nil
}

// ============ Listening History ============

// SetListeningHistory caches user's listening history
func (r *MusicCacheRepository) SetListeningHistory(ctx context.Context, userID uuid.UUID, history interface{}) error {
	key := fmt.Sprintf("history:%s", userID.String())
	
	// Cast the data to ListeningHistory
	listeningHistory, ok := history.(*ListeningHistory)
	if !ok {
		return fmt.Errorf("invalid data type for listening history")
	}
	
	listeningHistory.Timestamp = time.Now()
	
	jsonData, err := json.Marshal(listeningHistory)
	if err != nil {
		return fmt.Errorf("failed to marshal listening history: %w", err)
	}
	
	return r.client.Client.Set(ctx, key, jsonData, HistoryCacheTTL).Err()
}

// GetListeningHistory retrieves cached listening history
func (r *MusicCacheRepository) GetListeningHistory(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	key := fmt.Sprintf("history:%s", userID.String())
	
	data, err := r.client.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get listening history: %w", err)
	}
	
	var history ListeningHistory
	if err := json.Unmarshal([]byte(data), &history); err != nil {
		return nil, fmt.Errorf("failed to unmarshal listening history: %w", err)
	}
	
	return &history, nil
}

// ============ Popular Content Caching ============

// SetPopularAlbums caches popular albums for faster recommendations
func (r *MusicCacheRepository) SetPopularAlbums(ctx context.Context, albums []*models.Album) error {
	key := "popular:albums"
	
	jsonData, err := json.Marshal(albums)
	if err != nil {
		return fmt.Errorf("failed to marshal popular albums: %w", err)
	}
	
	return r.client.Client.Set(ctx, key, jsonData, PopularDataCacheTTL).Err()
}

// GetPopularAlbums retrieves cached popular albums
func (r *MusicCacheRepository) GetPopularAlbums(ctx context.Context) ([]*models.Album, error) {
	key := "popular:albums"
	
	data, err := r.client.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get popular albums: %w", err)
	}
	
	var albums []*models.Album
	if err := json.Unmarshal([]byte(data), &albums); err != nil {
		return nil, fmt.Errorf("failed to unmarshal popular albums: %w", err)
	}
	
	return albums, nil
}

// SetPopularTracks caches popular tracks
func (r *MusicCacheRepository) SetPopularTracks(ctx context.Context, tracks []*models.Track) error {
	key := "popular:tracks"
	
	jsonData, err := json.Marshal(tracks)
	if err != nil {
		return fmt.Errorf("failed to marshal popular tracks: %w", err)
	}
	
	return r.client.Client.Set(ctx, key, jsonData, PopularDataCacheTTL).Err()
}

// GetPopularTracks retrieves cached popular tracks
func (r *MusicCacheRepository) GetPopularTracks(ctx context.Context) ([]*models.Track, error) {
	key := "popular:tracks"
	
	data, err := r.client.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get popular tracks: %w", err)
	}
	
	var tracks []*models.Track
	if err := json.Unmarshal([]byte(data), &tracks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal popular tracks: %w", err)
	}
	
	return tracks, nil
}

// ============ Cache Management ============

// InvalidateUserCache removes all cached data for a user
func (r *MusicCacheRepository) InvalidateUserCache(ctx context.Context, userID uuid.UUID) error {
	keys := []string{
		fmt.Sprintf("user_music:%s", userID.String()),
		fmt.Sprintf("history:%s", userID.String()),
	}
	
	return r.client.Client.Del(ctx, keys...).Err()
}

// InvalidateSearchCache removes cached search results for a query
func (r *MusicCacheRepository) InvalidateSearchCache(ctx context.Context, query string) error {
	pattern := fmt.Sprintf("search:*:%s", query)
	
	// Get all matching keys
	keys, err := r.client.Client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	
	if len(keys) > 0 {
		return r.client.Client.Del(ctx, keys...).Err()
	}
	
	return nil
}

// GetCacheStats returns cache statistics
func (r *MusicCacheRepository) GetCacheStats(ctx context.Context) (map[string]int, error) {
	stats := make(map[string]int)
	
	// Count different types of cached data
	patterns := map[string]string{
		"user_music": "user_music:*",
		"searches":   "search:*",
		"history":    "history:*",
		"popular":    "popular:*",
		"sessions":   "session:*",
	}
	
	for name, pattern := range patterns {
		keys, err := r.client.Client.Keys(ctx, pattern).Result()
		if err != nil {
			continue // Skip on error
		}
		stats[name] = len(keys)
	}
	
	return stats, nil
} 