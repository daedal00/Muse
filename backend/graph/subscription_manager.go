package graph

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/daedal00/muse/backend/graph/model"
	"github.com/daedal00/muse/backend/internal/database"
)

// SubscriptionManager handles real-time subscriptions using Redis pub/sub
type SubscriptionManager struct {
	redis       *database.RedisClient
	subscribers map[string]map[chan *model.Review]bool // albumID -> subscribers
	mutex       sync.RWMutex
}

// NewSubscriptionManager creates a new subscription manager
func NewSubscriptionManager(redisClient *database.RedisClient) *SubscriptionManager {
	sm := &SubscriptionManager{
		redis:       redisClient,
		subscribers: make(map[string]map[chan *model.Review]bool),
	}
	
	// Start listening to Redis pub/sub
	go sm.listenToRedis()
	
	return sm
}

// Subscribe adds a new subscriber for review updates on a specific album
func (sm *SubscriptionManager) Subscribe(ctx context.Context, albumID string) (<-chan *model.Review, func()) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	// Create subscriber channel
	subscriber := make(chan *model.Review, 10) // Buffered channel
	
	// Initialize album subscribers map if needed
	if sm.subscribers[albumID] == nil {
		sm.subscribers[albumID] = make(map[chan *model.Review]bool)
	}
	
	// Add subscriber
	sm.subscribers[albumID][subscriber] = true
	
	// Create cleanup function
	cleanup := func() {
		sm.mutex.Lock()
		defer sm.mutex.Unlock()
		
		if sm.subscribers[albumID] != nil {
			delete(sm.subscribers[albumID], subscriber)
			
			// Clean up empty album subscriber map
			if len(sm.subscribers[albumID]) == 0 {
				delete(sm.subscribers, albumID)
			}
		}
		
		close(subscriber)
	}
	
	// Handle context cancellation
	go func() {
		<-ctx.Done()
		cleanup()
	}()
	
	return subscriber, cleanup
}

// PublishReview publishes a new review to Redis for distribution
func (sm *SubscriptionManager) PublishReview(ctx context.Context, review *model.Review) error {
	if review.Album == nil {
		return nil // Skip if no album
	}
	
	albumID := review.Album.ID
	
	// Serialize review
	reviewData, err := json.Marshal(review)
	if err != nil {
		return err
	}
	
	// Publish to Redis channel
	channelName := "reviews:" + albumID
	return sm.redis.Client.Publish(ctx, channelName, reviewData).Err()
}

// listenToRedis listens to Redis pub/sub for review updates
func (sm *SubscriptionManager) listenToRedis() {
	ctx := context.Background()
	
	// Subscribe to all review channels using pattern
	pubsub := sm.redis.Client.PSubscribe(ctx, "reviews:*")
	defer pubsub.Close()
	
	// Listen for messages
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("Error receiving pub/sub message: %v", err)
			time.Sleep(time.Second) // Wait before retrying
			continue
		}
		
		// Extract album ID from channel name
		channelName := msg.Channel
		if len(channelName) < 8 { // "reviews:" is 8 characters
			continue
		}
		albumID := channelName[8:] // Remove "reviews:" prefix
		
		// Deserialize review
		var review model.Review
		if err := json.Unmarshal([]byte(msg.Payload), &review); err != nil {
			log.Printf("Error deserializing review: %v", err)
			continue
		}
		
		// Distribute to local subscribers
		sm.distributeReview(albumID, &review)
	}
}

// distributeReview distributes a review to all local subscribers
func (sm *SubscriptionManager) distributeReview(albumID string, review *model.Review) {
	sm.mutex.RLock()
	subscribers := sm.subscribers[albumID]
	sm.mutex.RUnlock()
	
	if subscribers == nil {
		return // No subscribers for this album
	}
	
	// Send to all subscribers (non-blocking)
	for subscriber := range subscribers {
		select {
		case subscriber <- review:
			// Successfully sent
		default:
			// Channel is full or closed, skip
			log.Printf("Warning: Failed to send review to subscriber for album %s", albumID)
		}
	}
}

// GetActiveSubscriptions returns the number of active subscriptions per album
func (sm *SubscriptionManager) GetActiveSubscriptions() map[string]int {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	
	result := make(map[string]int)
	for albumID, subscribers := range sm.subscribers {
		result[albumID] = len(subscribers)
	}
	
	return result
} 