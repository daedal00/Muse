package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/daedal00/muse/backend/internal/database"
)

// CacheMetrics tracks cache performance metrics
type CacheMetrics struct {
	client *database.RedisClient
}

// CacheStats represents cache statistics for a specific period
type CacheStats struct {
	Type        string    `json:"type"`          // "album", "track", "search", etc.
	Hits        int64     `json:"hits"`          // Number of cache hits
	Misses      int64     `json:"misses"`        // Number of cache misses
	HitRate     float64   `json:"hit_rate"`      // Hit rate percentage
	TotalOps    int64     `json:"total_ops"`     // Total operations
	Period      string    `json:"period"`        // "hourly", "daily", "weekly"
	Timestamp   time.Time `json:"timestamp"`     // When this stat was recorded
	AvgHitTime  int64     `json:"avg_hit_time"`  // Average hit time in microseconds
	AvgMissTime int64     `json:"avg_miss_time"` // Average miss time in microseconds
}

// DetailedCacheStats provides comprehensive cache analysis
type DetailedCacheStats struct {
	Overall   CacheStats            `json:"overall"`
	ByType    map[string]CacheStats `json:"by_type"`
	ByHour    []CacheStats          `json:"by_hour"`
	TopMisses []string              `json:"top_misses"` // Most frequently missed keys
	TopHits   []string              `json:"top_hits"`   // Most frequently hit keys
}

func NewCacheMetrics(client *database.RedisClient) *CacheMetrics {
	return &CacheMetrics{client: client}
}

// RecordHit records a cache hit with timing information
func (m *CacheMetrics) RecordHit(ctx context.Context, cacheType, key string, duration time.Duration) error {
	now := time.Now()
	hour := now.Format("2006-01-02-15") // YYYY-MM-DD-HH format

	pipe := m.client.Client.Pipeline()

	// Increment hit counters
	pipe.Incr(ctx, fmt.Sprintf("cache:hits:%s:%s", cacheType, hour))
	pipe.Incr(ctx, fmt.Sprintf("cache:hits:%s:total", cacheType))
	pipe.Incr(ctx, "cache:hits:total")

	// Record timing (convert to microseconds)
	timingKey := fmt.Sprintf("cache:timing:hits:%s:%s", cacheType, hour)
	pipe.LPush(ctx, timingKey, duration.Microseconds())
	pipe.LTrim(ctx, timingKey, 0, 999)        // Keep last 1000 measurements
	pipe.Expire(ctx, timingKey, 25*time.Hour) // Expire after 25 hours

	// Record popular hits (for analysis)
	popularHitsKey := fmt.Sprintf("cache:popular:hits:%s", cacheType)
	pipe.ZIncrBy(ctx, popularHitsKey, 1, key)
	pipe.Expire(ctx, popularHitsKey, 7*24*time.Hour) // Keep for 7 days

	// Set expiration for hourly counters
	pipe.Expire(ctx, fmt.Sprintf("cache:hits:%s:%s", cacheType, hour), 25*time.Hour)

	_, err := pipe.Exec(ctx)
	return err
}

// RecordMiss records a cache miss with timing information
func (m *CacheMetrics) RecordMiss(ctx context.Context, cacheType, key string, duration time.Duration) error {
	now := time.Now()
	hour := now.Format("2006-01-02-15")

	pipe := m.client.Client.Pipeline()

	// Increment miss counters
	pipe.Incr(ctx, fmt.Sprintf("cache:misses:%s:%s", cacheType, hour))
	pipe.Incr(ctx, fmt.Sprintf("cache:misses:%s:total", cacheType))
	pipe.Incr(ctx, "cache:misses:total")

	// Record timing
	timingKey := fmt.Sprintf("cache:timing:misses:%s:%s", cacheType, hour)
	pipe.LPush(ctx, timingKey, duration.Microseconds())
	pipe.LTrim(ctx, timingKey, 0, 999)
	pipe.Expire(ctx, timingKey, 25*time.Hour)

	// Record frequent misses (for cache warming opportunities)
	frequentMissesKey := fmt.Sprintf("cache:frequent:misses:%s", cacheType)
	pipe.ZIncrBy(ctx, frequentMissesKey, 1, key)
	pipe.Expire(ctx, frequentMissesKey, 7*24*time.Hour)

	// Set expiration for hourly counters
	pipe.Expire(ctx, fmt.Sprintf("cache:misses:%s:%s", cacheType, hour), 25*time.Hour)

	_, err := pipe.Exec(ctx)
	return err
}

// GetHitRate gets the current hit rate for a specific cache type
func (m *CacheMetrics) GetHitRate(ctx context.Context, cacheType string) (float64, error) {
	pipe := m.client.Client.Pipeline()
	hitsCmd := pipe.Get(ctx, fmt.Sprintf("cache:hits:%s:total", cacheType))
	missesCmd := pipe.Get(ctx, fmt.Sprintf("cache:misses:%s:total", cacheType))

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	hits, _ := hitsCmd.Int64()
	misses, _ := missesCmd.Int64()

	total := hits + misses
	if total == 0 {
		return 0, nil
	}

	return float64(hits) / float64(total) * 100, nil
}

// GetHourlyStats gets cache statistics for the last N hours
func (m *CacheMetrics) GetHourlyStats(ctx context.Context, cacheType string, hours int) ([]CacheStats, error) {
	stats := make([]CacheStats, 0, hours)
	now := time.Now()

	for i := 0; i < hours; i++ {
		hour := now.Add(-time.Duration(i) * time.Hour).Format("2006-01-02-15")

		pipe := m.client.Client.Pipeline()
		hitsCmd := pipe.Get(ctx, fmt.Sprintf("cache:hits:%s:%s", cacheType, hour))
		missesCmd := pipe.Get(ctx, fmt.Sprintf("cache:misses:%s:%s", cacheType, hour))
		hitTimingsCmd := pipe.LRange(ctx, fmt.Sprintf("cache:timing:hits:%s:%s", cacheType, hour), 0, -1)
		missTimingsCmd := pipe.LRange(ctx, fmt.Sprintf("cache:timing:misses:%s:%s", cacheType, hour), 0, -1)

		_, err := pipe.Exec(ctx)
		if err != nil {
			continue
		}

		hits, _ := hitsCmd.Int64()
		misses, _ := missesCmd.Int64()
		total := hits + misses

		hitRate := float64(0)
		if total > 0 {
			hitRate = float64(hits) / float64(total) * 100
		}

		// Calculate average timings
		avgHitTime := int64(0)
		avgMissTime := int64(0)

		hitTimings := hitTimingsCmd.Val()
		if len(hitTimings) > 0 {
			var sum int64
			count := 0
			for _, timing := range hitTimings {
				if val, err := strconv.ParseInt(timing, 10, 64); err == nil {
					sum += val
					count++
				}
			}
			if count > 0 {
				avgHitTime = sum / int64(count)
			}
		}

		missTimings := missTimingsCmd.Val()
		if len(missTimings) > 0 {
			var sum int64
			count := 0
			for _, timing := range missTimings {
				if val, err := strconv.ParseInt(timing, 10, 64); err == nil {
					sum += val
					count++
				}
			}
			if count > 0 {
				avgMissTime = sum / int64(count)
			}
		}

		stat := CacheStats{
			Type:        cacheType,
			Hits:        hits,
			Misses:      misses,
			HitRate:     hitRate,
			TotalOps:    total,
			Period:      "hourly",
			Timestamp:   now.Add(-time.Duration(i) * time.Hour),
			AvgHitTime:  avgHitTime,
			AvgMissTime: avgMissTime,
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

// GetDetailedStats provides comprehensive cache analysis
func (m *CacheMetrics) GetDetailedStats(ctx context.Context) (*DetailedCacheStats, error) {
	cacheTypes := []string{"album", "track", "artist", "search", "playlist", "user"}

	detailed := &DetailedCacheStats{
		ByType: make(map[string]CacheStats),
	}

	// Get overall stats
	var totalHits, totalMisses int64

	for _, cacheType := range cacheTypes {
		hits, _ := m.client.Client.Get(ctx, fmt.Sprintf("cache:hits:%s:total", cacheType)).Int64()
		misses, _ := m.client.Client.Get(ctx, fmt.Sprintf("cache:misses:%s:total", cacheType)).Int64()

		total := hits + misses
		hitRate := float64(0)
		if total > 0 {
			hitRate = float64(hits) / float64(total) * 100
		}

		detailed.ByType[cacheType] = CacheStats{
			Type:     cacheType,
			Hits:     hits,
			Misses:   misses,
			HitRate:  hitRate,
			TotalOps: total,
			Period:   "total",
		}

		totalHits += hits
		totalMisses += misses
	}

	// Overall stats
	totalOps := totalHits + totalMisses
	overallHitRate := float64(0)
	if totalOps > 0 {
		overallHitRate = float64(totalHits) / float64(totalOps) * 100
	}

	detailed.Overall = CacheStats{
		Type:     "overall",
		Hits:     totalHits,
		Misses:   totalMisses,
		HitRate:  overallHitRate,
		TotalOps: totalOps,
		Period:   "total",
	}

	// Get hourly breakdown for the most active cache type
	mostActiveType := "album" // Default
	maxOps := int64(0)
	for cacheType, stats := range detailed.ByType {
		if stats.TotalOps > maxOps {
			maxOps = stats.TotalOps
			mostActiveType = cacheType
		}
	}

	if hourlyStats, err := m.GetHourlyStats(ctx, mostActiveType, 24); err == nil {
		detailed.ByHour = hourlyStats
	}

	// Get top misses and hits
	if topMisses, err := m.getTopItems(ctx, "misses", "album", 10); err == nil {
		detailed.TopMisses = topMisses
	}

	if topHits, err := m.getTopItems(ctx, "hits", "album", 10); err == nil {
		detailed.TopHits = topHits
	}

	return detailed, nil
}

// getTopItems gets the most frequently hit or missed items
func (m *CacheMetrics) getTopItems(ctx context.Context, itemType, cacheType string, limit int) ([]string, error) {
	var keyPrefix string
	if itemType == "hits" {
		keyPrefix = "popular"
	} else {
		keyPrefix = "frequent"
	}

	key := fmt.Sprintf("cache:%s:%s:%s", keyPrefix, itemType, cacheType)

	results, err := m.client.Client.ZRevRange(ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetRealtimeHitRate gets the hit rate for the current hour
func (m *CacheMetrics) GetRealtimeHitRate(ctx context.Context, cacheType string) (float64, int64, int64, error) {
	now := time.Now()
	hour := now.Format("2006-01-02-15")

	pipe := m.client.Client.Pipeline()
	hitsCmd := pipe.Get(ctx, fmt.Sprintf("cache:hits:%s:%s", cacheType, hour))
	missesCmd := pipe.Get(ctx, fmt.Sprintf("cache:misses:%s:%s", cacheType, hour))

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, 0, 0, err
	}

	hits, _ := hitsCmd.Int64()
	misses, _ := missesCmd.Int64()
	total := hits + misses

	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(hits) / float64(total) * 100
	}

	return hitRate, hits, misses, nil
}

// ExportMetrics exports cache metrics as JSON for external monitoring
func (m *CacheMetrics) ExportMetrics(ctx context.Context) (string, error) {
	stats, err := m.GetDetailedStats(ctx)
	if err != nil {
		return "", err
	}

	jsonData, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// ResetMetrics clears all cache metrics (use with caution)
func (m *CacheMetrics) ResetMetrics(ctx context.Context) error {
	patterns := []string{
		"cache:hits:*",
		"cache:misses:*",
		"cache:timing:*",
		"cache:popular:*",
		"cache:frequent:*",
	}

	for _, pattern := range patterns {
		keys, err := m.client.Client.Keys(ctx, pattern).Result()
		if err != nil {
			continue
		}

		if len(keys) > 0 {
			m.client.Client.Del(ctx, keys...)
		}
	}

	return nil
}
