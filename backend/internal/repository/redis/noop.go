package redis

import (
	"context"
	"fmt"

	"github.com/daedal00/muse/backend/internal/models"
	"github.com/daedal00/muse/backend/internal/repository"
	"github.com/google/uuid"
)

type NoopSessionRepository struct{}

func NewNoopSessionRepository() repository.SessionRepository {
	return &NoopSessionRepository{}
}

func (r *NoopSessionRepository) Create(ctx context.Context, session *models.Session) error {
	return nil
}

func (r *NoopSessionRepository) GetByID(ctx context.Context, id string) (*models.Session, error) {
	return nil, fmt.Errorf("session not found")
}

func (r *NoopSessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Session, error) {
	return []*models.Session{}, nil
}

func (r *NoopSessionRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *NoopSessionRepository) DeleteExpired(ctx context.Context) error {
	return nil
}

func (r *NoopSessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return nil
}

type NoopMusicCacheRepository struct{}

func NewNoopMusicCacheRepository() repository.MusicCacheRepository {
	return &NoopMusicCacheRepository{}
}

func (r *NoopMusicCacheRepository) SetUserMusicData(ctx context.Context, userID uuid.UUID, data interface{}) error {
	return nil
}

func (r *NoopMusicCacheRepository) GetUserMusicData(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	return nil, nil
}

func (r *NoopMusicCacheRepository) AddToRecentlyPlayed(ctx context.Context, userID uuid.UUID, track *models.Track) error {
	return nil
}

func (r *NoopMusicCacheRepository) SetSearchResults(ctx context.Context, query string, resultType string, results interface{}) error {
	return nil
}

func (r *NoopMusicCacheRepository) GetSearchResults(ctx context.Context, query string, resultType string) (interface{}, error) {
	return nil, nil
}

func (r *NoopMusicCacheRepository) SetListeningHistory(ctx context.Context, userID uuid.UUID, history interface{}) error {
	return nil
}

func (r *NoopMusicCacheRepository) GetListeningHistory(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	return nil, nil
}

func (r *NoopMusicCacheRepository) SetPopularAlbums(ctx context.Context, albums []*models.Album) error {
	return nil
}

func (r *NoopMusicCacheRepository) GetPopularAlbums(ctx context.Context) ([]*models.Album, error) {
	return nil, nil
}

func (r *NoopMusicCacheRepository) SetPopularTracks(ctx context.Context, tracks []*models.Track) error {
	return nil
}

func (r *NoopMusicCacheRepository) GetPopularTracks(ctx context.Context) ([]*models.Track, error) {
	return nil, nil
}

func (r *NoopMusicCacheRepository) InvalidateUserCache(ctx context.Context, userID uuid.UUID) error {
	return nil
}

func (r *NoopMusicCacheRepository) InvalidateSearchCache(ctx context.Context, query string) error {
	return nil
}

func (r *NoopMusicCacheRepository) GetCacheStats(ctx context.Context) (map[string]int, error) {
	return map[string]int{}, nil
}
