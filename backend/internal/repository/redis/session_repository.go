package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/daedal00/muse/backend/internal/database"
	"github.com/daedal00/muse/backend/internal/models"
	"github.com/daedal00/muse/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type sessionRepository struct {
	client *database.RedisClient
}

func NewSessionRepository(client *database.RedisClient) repository.SessionRepository {
	return &sessionRepository{client: client}
}

func (r *sessionRepository) Create(ctx context.Context, session *models.Session) error {
	sessionKey := fmt.Sprintf("session:%s", session.ID)
	userSessionsKey := fmt.Sprintf("user_sessions:%s", session.UserID.String())

	// Serialize session data
	sessionData := map[string]interface{}{
		"id":         session.ID,
		"user_id":    session.UserID.String(),
		"expires_at": session.ExpiresAt.Unix(),
		"created_at": session.CreatedAt.Unix(),
	}

	serializedData, err := json.Marshal(sessionData)
	if err != nil {
		return fmt.Errorf("failed to serialize session: %w", err)
	}

	// Calculate TTL
	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("session already expired")
	}

	// Use pipeline for atomic operations
	pipe := r.client.Client.Pipeline()

	// Store session data with expiration
	pipe.Set(ctx, sessionKey, serializedData, ttl)

	// Add session ID to user's session set with same expiration
	pipe.SAdd(ctx, userSessionsKey, session.ID)
	pipe.Expire(ctx, userSessionsKey, ttl)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (r *sessionRepository) GetByID(ctx context.Context, id string) (*models.Session, error) {
	sessionKey := fmt.Sprintf("session:%s", id)

	data, err := r.client.Client.Get(ctx, sessionKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("session not found or expired")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var sessionData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &sessionData); err != nil {
		return nil, fmt.Errorf("failed to deserialize session: %w", err)
	}

	userID, err := uuid.Parse(sessionData["user_id"].(string))
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in session: %w", err)
	}

	expiresAt := time.Unix(int64(sessionData["expires_at"].(float64)), 0)
	createdAt := time.Unix(int64(sessionData["created_at"].(float64)), 0)

	return &models.Session{
		ID:        sessionData["id"].(string),
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
	}, nil
}

func (r *sessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Session, error) {
	userSessionsKey := fmt.Sprintf("user_sessions:%s", userID.String())

	// Get all session IDs for the user
	sessionIDs, err := r.client.Client.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		if err == redis.Nil {
			return []*models.Session{}, nil
		}
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}

	if len(sessionIDs) == 0 {
		return []*models.Session{}, nil
	}

	// Build keys for batch retrieval
	sessionKeys := make([]string, len(sessionIDs))
	for i, sessionID := range sessionIDs {
		sessionKeys[i] = fmt.Sprintf("session:%s", sessionID)
	}

	// Get all sessions in batch
	results, err := r.client.Client.MGet(ctx, sessionKeys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}

	var sessions []*models.Session
	for i, result := range results {
		if result == nil {
			// Session expired, remove from user's session set
			r.client.Client.SRem(ctx, userSessionsKey, sessionIDs[i])
			continue
		}

		var sessionData map[string]interface{}
		if err := json.Unmarshal([]byte(result.(string)), &sessionData); err != nil {
			continue // Skip invalid sessions
		}

		userIDParsed, err := uuid.Parse(sessionData["user_id"].(string))
		if err != nil {
			continue // Skip invalid sessions
		}

		expiresAt := time.Unix(int64(sessionData["expires_at"].(float64)), 0)
		createdAt := time.Unix(int64(sessionData["created_at"].(float64)), 0)

		sessions = append(sessions, &models.Session{
			ID:        sessionData["id"].(string),
			UserID:    userIDParsed,
			ExpiresAt: expiresAt,
			CreatedAt: createdAt,
		})
	}

	return sessions, nil
}

func (r *sessionRepository) Delete(ctx context.Context, id string) error {
	// First get the session to find the user ID
	session, err := r.GetByID(ctx, id)
	if err != nil {
		return err // Will return "session not found" if it doesn't exist
	}

	sessionKey := fmt.Sprintf("session:%s", id)
	userSessionsKey := fmt.Sprintf("user_sessions:%s", session.UserID.String())

	// Use pipeline for atomic operations
	pipe := r.client.Client.Pipeline()

	// Delete session data
	pipe.Del(ctx, sessionKey)

	// Remove session ID from user's session set
	pipe.SRem(ctx, userSessionsKey, id)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (r *sessionRepository) DeleteExpired(ctx context.Context) error {
	// Redis automatically handles expiration, but we can clean up user session sets
	// This is a maintenance operation to clean orphaned user session sets

	pattern := "user_sessions:*"
	iter := r.client.Client.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		userSessionsKey := iter.Val()

		// Get all session IDs for this user
		sessionIDs, err := r.client.Client.SMembers(ctx, userSessionsKey).Result()
		if err != nil {
			continue
		}

		// Check which sessions still exist
		if len(sessionIDs) > 0 {
			sessionKeys := make([]string, len(sessionIDs))
			for i, sessionID := range sessionIDs {
				sessionKeys[i] = fmt.Sprintf("session:%s", sessionID)
			}

			// Check existence of sessions
			pipe := r.client.Client.Pipeline()
			for _, key := range sessionKeys {
				pipe.Exists(ctx, key)
			}

			results, err := pipe.Exec(ctx)
			if err != nil {
				continue
			}

			// Remove expired session IDs from user set
			expiredSessionIDs := []interface{}{}
			for i, result := range results {
				if result.(*redis.IntCmd).Val() == 0 {
					expiredSessionIDs = append(expiredSessionIDs, sessionIDs[i])
				}
			}

			if len(expiredSessionIDs) > 0 {
				r.client.Client.SRem(ctx, userSessionsKey, expiredSessionIDs...)
			}

			// If user has no sessions left, delete the set
			count, _ := r.client.Client.SCard(ctx, userSessionsKey).Result()
			if count == 0 {
				r.client.Client.Del(ctx, userSessionsKey)
			}
		}
	}

	return iter.Err()
}

func (r *sessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	userSessionsKey := fmt.Sprintf("user_sessions:%s", userID.String())

	// Get all session IDs for the user
	sessionIDs, err := r.client.Client.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil // No sessions to delete
		}
		return fmt.Errorf("failed to get user sessions: %w", err)
	}

	if len(sessionIDs) == 0 {
		return nil
	}

	// Build keys for batch deletion
	sessionKeys := make([]string, len(sessionIDs)+1)
	for i, sessionID := range sessionIDs {
		sessionKeys[i] = fmt.Sprintf("session:%s", sessionID)
	}
	sessionKeys[len(sessionIDs)] = userSessionsKey // Also delete the user sessions set

	// Delete all sessions and the user session set
	err = r.client.Client.Del(ctx, sessionKeys...).Err()
	if err != nil {
		return fmt.Errorf("failed to delete sessions by user: %w", err)
	}

	return nil
}
