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

// SpotifyCacheRepository handles specialized caching for Spotify data
type SpotifyCacheRepository struct {
	client *database.RedisClient
}

// Cache durations optimized for Spotify data
const (
	SpotifyTrackCacheTTL      = 24 * time.Hour   // Tracks don't change often
	SpotifyAlbumCacheTTL      = 24 * time.Hour   // Albums don't change often
	SpotifyArtistCacheTTL     = 12 * time.Hour   // Artist info might update more frequently
	SpotifyPlaylistCacheTTL   = 1 * time.Hour    // Playlists can change frequently
	SpotifySearchCacheTTL     = 30 * time.Minute // Search results cache
	SpotifyUserDataCacheTTL   = 15 * time.Minute // User's Spotify data (recently played, top tracks, etc.)
	SpotifyRecommendationsTTL = 2 * time.Hour    // Recommendations cache
	SpotifyTokenCacheTTL      = 50 * time.Minute // Access token cache (expires in 1 hour, cache for 50 min)
)

func NewSpotifyCacheRepository(client *database.RedisClient) *SpotifyCacheRepository {
	return &SpotifyCacheRepository{client: client}
}

// ============ Track Caching ============

func (r *SpotifyCacheRepository) SetTrack(ctx context.Context, track *models.SpotifyTrack) error {
	key := fmt.Sprintf("spotify:track:%s", track.ID)

	jsonData, err := json.Marshal(track)
	if err != nil {
		return fmt.Errorf("failed to marshal track: %w", err)
	}

	return r.client.Client.Set(ctx, key, jsonData, SpotifyTrackCacheTTL).Err()
}

func (r *SpotifyCacheRepository) GetTrack(ctx context.Context, trackID string) (*models.SpotifyTrack, error) {
	key := fmt.Sprintf("spotify:track:%s", trackID)

	data, err := r.client.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get track: %w", err)
	}

	var track models.SpotifyTrack
	if err := json.Unmarshal([]byte(data), &track); err != nil {
		return nil, fmt.Errorf("failed to unmarshal track: %w", err)
	}

	return &track, nil
}

func (r *SpotifyCacheRepository) SetTracks(ctx context.Context, tracks []models.SpotifyTrack) error {
	if len(tracks) == 0 {
		return nil
	}

	pipe := r.client.Client.Pipeline()

	for _, track := range tracks {
		key := fmt.Sprintf("spotify:track:%s", track.ID)
		jsonData, err := json.Marshal(track)
		if err != nil {
			continue // Skip invalid tracks
		}
		pipe.Set(ctx, key, jsonData, SpotifyTrackCacheTTL)
	}

	_, err := pipe.Exec(ctx)
	return err
}

func (r *SpotifyCacheRepository) GetTracks(ctx context.Context, trackIDs []string) ([]models.SpotifyTrack, []string, error) {
	if len(trackIDs) == 0 {
		return []models.SpotifyTrack{}, []string{}, nil
	}

	keys := make([]string, len(trackIDs))
	for i, id := range trackIDs {
		keys[i] = fmt.Sprintf("spotify:track:%s", id)
	}

	results, err := r.client.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, trackIDs, err // Return all as missing
	}

	var tracks []models.SpotifyTrack
	var missingIDs []string

	for i, result := range results {
		if result == nil {
			missingIDs = append(missingIDs, trackIDs[i])
			continue
		}

		var track models.SpotifyTrack
		if err := json.Unmarshal([]byte(result.(string)), &track); err != nil {
			missingIDs = append(missingIDs, trackIDs[i])
			continue
		}

		tracks = append(tracks, track)
	}

	return tracks, missingIDs, nil
}

// ============ Album Caching ============

func (r *SpotifyCacheRepository) SetAlbum(ctx context.Context, album *models.SpotifyAlbum) error {
	key := fmt.Sprintf("spotify:album:%s", album.ID)

	jsonData, err := json.Marshal(album)
	if err != nil {
		return fmt.Errorf("failed to marshal album: %w", err)
	}

	return r.client.Client.Set(ctx, key, jsonData, SpotifyAlbumCacheTTL).Err()
}

func (r *SpotifyCacheRepository) GetAlbum(ctx context.Context, albumID string) (*models.SpotifyAlbum, error) {
	key := fmt.Sprintf("spotify:album:%s", albumID)

	data, err := r.client.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get album: %w", err)
	}

	var album models.SpotifyAlbum
	if err := json.Unmarshal([]byte(data), &album); err != nil {
		return nil, fmt.Errorf("failed to unmarshal album: %w", err)
	}

	return &album, nil
}

func (r *SpotifyCacheRepository) SetAlbums(ctx context.Context, albums []models.SpotifyAlbum) error {
	if len(albums) == 0 {
		return nil
	}

	pipe := r.client.Client.Pipeline()

	for _, album := range albums {
		key := fmt.Sprintf("spotify:album:%s", album.ID)
		jsonData, err := json.Marshal(album)
		if err != nil {
			continue
		}
		pipe.Set(ctx, key, jsonData, SpotifyAlbumCacheTTL)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// ============ Artist Caching ============

func (r *SpotifyCacheRepository) SetArtist(ctx context.Context, artist *models.SpotifyArtist) error {
	key := fmt.Sprintf("spotify:artist:%s", artist.ID)

	jsonData, err := json.Marshal(artist)
	if err != nil {
		return fmt.Errorf("failed to marshal artist: %w", err)
	}

	return r.client.Client.Set(ctx, key, jsonData, SpotifyArtistCacheTTL).Err()
}

func (r *SpotifyCacheRepository) GetArtist(ctx context.Context, artistID string) (*models.SpotifyArtist, error) {
	key := fmt.Sprintf("spotify:artist:%s", artistID)

	data, err := r.client.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get artist: %w", err)
	}

	var artist models.SpotifyArtist
	if err := json.Unmarshal([]byte(data), &artist); err != nil {
		return nil, fmt.Errorf("failed to unmarshal artist: %w", err)
	}

	return &artist, nil
}

func (r *SpotifyCacheRepository) SetArtists(ctx context.Context, artists []models.SpotifyArtist) error {
	if len(artists) == 0 {
		return nil
	}

	pipe := r.client.Client.Pipeline()

	for _, artist := range artists {
		key := fmt.Sprintf("spotify:artist:%s", artist.ID)
		jsonData, err := json.Marshal(artist)
		if err != nil {
			continue
		}
		pipe.Set(ctx, key, jsonData, SpotifyArtistCacheTTL)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// ============ User Data Caching ============

func (r *SpotifyCacheRepository) SetUserData(ctx context.Context, userID uuid.UUID, data *models.CachedUserData) error {
	key := fmt.Sprintf("spotify:user:%s", userID.String())

	data.LastUpdated = time.Now()
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal user data: %w", err)
	}

	return r.client.Client.Set(ctx, key, jsonData, SpotifyUserDataCacheTTL).Err()
}

func (r *SpotifyCacheRepository) GetUserData(ctx context.Context, userID uuid.UUID) (*models.CachedUserData, error) {
	key := fmt.Sprintf("spotify:user:%s", userID.String())

	data, err := r.client.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user data: %w", err)
	}

	var userData models.CachedUserData
	if err := json.Unmarshal([]byte(data), &userData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	return &userData, nil
}

// ============ Recommendations Caching ============

func (r *SpotifyCacheRepository) SetRecommendations(ctx context.Context, userID uuid.UUID, recommendations *models.CachedRecommendations) error {
	key := fmt.Sprintf("spotify:recommendations:%s", userID.String())

	recommendations.GeneratedAt = time.Now()
	jsonData, err := json.Marshal(recommendations)
	if err != nil {
		return fmt.Errorf("failed to marshal recommendations: %w", err)
	}

	return r.client.Client.Set(ctx, key, jsonData, SpotifyRecommendationsTTL).Err()
}

func (r *SpotifyCacheRepository) GetRecommendations(ctx context.Context, userID uuid.UUID) (*models.CachedRecommendations, error) {
	key := fmt.Sprintf("spotify:recommendations:%s", userID.String())

	data, err := r.client.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get recommendations: %w", err)
	}

	var recommendations models.CachedRecommendations
	if err := json.Unmarshal([]byte(data), &recommendations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal recommendations: %w", err)
	}

	return &recommendations, nil
}

// ============ Search Results Caching ============

func (r *SpotifyCacheRepository) SetSearchResults(ctx context.Context, query, resultType string, results interface{}) error {
	key := fmt.Sprintf("spotify:search:%s:%s", resultType, query)

	cacheData := map[string]interface{}{
		"query":       query,
		"results":     results,
		"timestamp":   time.Now().Unix(),
		"result_type": resultType,
	}

	jsonData, err := json.Marshal(cacheData)
	if err != nil {
		return fmt.Errorf("failed to marshal search results: %w", err)
	}

	return r.client.Client.Set(ctx, key, jsonData, SpotifySearchCacheTTL).Err()
}

func (r *SpotifyCacheRepository) GetSearchResults(ctx context.Context, query, resultType string) (interface{}, error) {
	key := fmt.Sprintf("spotify:search:%s:%s", resultType, query)

	data, err := r.client.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get search results: %w", err)
	}

	var cacheData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &cacheData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}

	return cacheData["results"], nil
}

// ============ Token Caching ============

func (r *SpotifyCacheRepository) SetAccessToken(ctx context.Context, userID uuid.UUID, token string) error {
	key := fmt.Sprintf("spotify:token:%s", userID.String())
	return r.client.Client.Set(ctx, key, token, SpotifyTokenCacheTTL).Err()
}

func (r *SpotifyCacheRepository) GetAccessToken(ctx context.Context, userID uuid.UUID) (string, error) {
	key := fmt.Sprintf("spotify:token:%s", userID.String())

	token, err := r.client.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("failed to get access token: %w", err)
	}

	return token, nil
}

func (r *SpotifyCacheRepository) DeleteAccessToken(ctx context.Context, userID uuid.UUID) error {
	key := fmt.Sprintf("spotify:token:%s", userID.String())
	return r.client.Client.Del(ctx, key).Err()
}

// ============ Cache Management ============

func (r *SpotifyCacheRepository) InvalidateUserCache(ctx context.Context, userID uuid.UUID) error {
	pattern := fmt.Sprintf("spotify:*:%s", userID.String())

	keys, err := r.client.Client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return r.client.Client.Del(ctx, keys...).Err()
	}

	return nil
}

func (r *SpotifyCacheRepository) InvalidateSearchCache(ctx context.Context, query string) error {
	pattern := fmt.Sprintf("spotify:search:*:%s", query)

	keys, err := r.client.Client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return r.client.Client.Del(ctx, keys...).Err()
	}

	return nil
}

func (r *SpotifyCacheRepository) GetCacheStats(ctx context.Context) (map[string]int, error) {
	stats := make(map[string]int)

	patterns := map[string]string{
		"tracks":          "spotify:track:*",
		"albums":          "spotify:album:*",
		"artists":         "spotify:artist:*",
		"user_data":       "spotify:user:*",
		"recommendations": "spotify:recommendations:*",
		"searches":        "spotify:search:*",
		"tokens":          "spotify:token:*",
	}

	for name, pattern := range patterns {
		keys, err := r.client.Client.Keys(ctx, pattern).Result()
		if err != nil {
			continue
		}
		stats[name] = len(keys)
	}

	return stats, nil
}

// ============ Batch Operations ============

func (r *SpotifyCacheRepository) WarmUpCache(ctx context.Context, userID uuid.UUID, tracks []models.SpotifyTrack, albums []models.SpotifyAlbum, artists []models.SpotifyArtist) error {
	pipe := r.client.Client.Pipeline()

	// Cache tracks
	for _, track := range tracks {
		key := fmt.Sprintf("spotify:track:%s", track.ID)
		jsonData, err := json.Marshal(track)
		if err != nil {
			continue
		}
		pipe.Set(ctx, key, jsonData, SpotifyTrackCacheTTL)
	}

	// Cache albums
	for _, album := range albums {
		key := fmt.Sprintf("spotify:album:%s", album.ID)
		jsonData, err := json.Marshal(album)
		if err != nil {
			continue
		}
		pipe.Set(ctx, key, jsonData, SpotifyAlbumCacheTTL)
	}

	// Cache artists
	for _, artist := range artists {
		key := fmt.Sprintf("spotify:artist:%s", artist.ID)
		jsonData, err := json.Marshal(artist)
		if err != nil {
			continue
		}
		pipe.Set(ctx, key, jsonData, SpotifyArtistCacheTTL)
	}

	_, err := pipe.Exec(ctx)
	return err
}
