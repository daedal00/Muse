package repository

import (
	"context"

	"github.com/daedal00/muse/backend/internal/models"
	"github.com/google/uuid"
)

// UserRepository handles user data operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetBySpotifyID(ctx context.Context, spotifyID string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
}

// ReviewRepository handles review operations (now using Spotify IDs)
type ReviewRepository interface {
	Create(ctx context.Context, review *models.Review) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Review, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Review, error)
	GetBySpotifyID(ctx context.Context, spotifyID string, spotifyType string, limit, offset int) ([]*models.Review, error)
	Update(ctx context.Context, review *models.Review) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.Review, error)
	GetAverageRating(ctx context.Context, spotifyID string, spotifyType string) (float64, error)
	GetRatingCount(ctx context.Context, spotifyID string, spotifyType string) (int, error)
}

// PlaylistRepository handles playlist operations
type PlaylistRepository interface {
	Create(ctx context.Context, playlist *models.Playlist) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Playlist, error)
	GetByCreatorID(ctx context.Context, creatorID uuid.UUID, limit, offset int) ([]*models.Playlist, error)
	GetPublicPlaylists(ctx context.Context, limit, offset int) ([]*models.Playlist, error)
	Update(ctx context.Context, playlist *models.Playlist) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.Playlist, error)

	// Playlist track operations (now using Spotify IDs)
	AddTrack(ctx context.Context, playlistID uuid.UUID, spotifyID string, position int, addedByUserID uuid.UUID) error
	RemoveTrack(ctx context.Context, playlistID uuid.UUID, spotifyID string) error
	GetTrackSpotifyIDs(ctx context.Context, playlistID uuid.UUID, limit, offset int) ([]string, error)
	GetTrackCount(ctx context.Context, playlistID uuid.UUID) (int, error)
	ReorderTracks(ctx context.Context, playlistID uuid.UUID, trackPositions map[string]int) error
	GetPlaylistTracks(ctx context.Context, playlistID uuid.UUID, limit, offset int) ([]*models.PlaylistTrack, error)

	// Track management methods
	GetTracksByPlaylistID(ctx context.Context, playlistID uuid.UUID, limit, offset int) ([]*models.PlaylistTrack, error)
}

// UserPreferencesRepository handles user preferences and settings
type UserPreferencesRepository interface {
	Create(ctx context.Context, preferences *models.UserPreferences) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserPreferences, error)
	Update(ctx context.Context, preferences *models.UserPreferences) error
	Delete(ctx context.Context, userID uuid.UUID) error
	AddFavoriteArtist(ctx context.Context, userID uuid.UUID, spotifyArtistID string) error
	RemoveFavoriteArtist(ctx context.Context, userID uuid.UUID, spotifyArtistID string) error
	UpdatePreferredGenres(ctx context.Context, userID uuid.UUID, genres []string) error
}

// SessionRepository handles user sessions (Redis-based)
type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	GetByID(ctx context.Context, id string) (*models.Session, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Session, error)
	Delete(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

// SpotifyCacheRepository handles caching of Spotify data
type SpotifyCacheRepository interface {
	// Track caching
	SetTrack(ctx context.Context, track *models.SpotifyTrack) error
	GetTrack(ctx context.Context, trackID string) (*models.SpotifyTrack, error)
	SetTracks(ctx context.Context, tracks []models.SpotifyTrack) error
	GetTracks(ctx context.Context, trackIDs []string) ([]models.SpotifyTrack, []string, error)

	// Album caching
	SetAlbum(ctx context.Context, album *models.SpotifyAlbum) error
	GetAlbum(ctx context.Context, albumID string) (*models.SpotifyAlbum, error)
	SetAlbums(ctx context.Context, albums []models.SpotifyAlbum) error

	// Artist caching
	SetArtist(ctx context.Context, artist *models.SpotifyArtist) error
	GetArtist(ctx context.Context, artistID string) (*models.SpotifyArtist, error)
	SetArtists(ctx context.Context, artists []models.SpotifyArtist) error

	// User data caching
	SetUserData(ctx context.Context, userID uuid.UUID, data *models.CachedUserData) error
	GetUserData(ctx context.Context, userID uuid.UUID) (*models.CachedUserData, error)

	// Recommendations caching
	SetRecommendations(ctx context.Context, userID uuid.UUID, recommendations *models.CachedRecommendations) error
	GetRecommendations(ctx context.Context, userID uuid.UUID) (*models.CachedRecommendations, error)

	// Search results caching
	SetSearchResults(ctx context.Context, query, resultType string, results interface{}) error
	GetSearchResults(ctx context.Context, query, resultType string) (interface{}, error)

	// Token caching
	SetAccessToken(ctx context.Context, userID uuid.UUID, token string) error
	GetAccessToken(ctx context.Context, userID uuid.UUID) (string, error)
	DeleteAccessToken(ctx context.Context, userID uuid.UUID) error

	// Cache management
	InvalidateUserCache(ctx context.Context, userID uuid.UUID) error
	InvalidateSearchCache(ctx context.Context, query string) error
	GetCacheStats(ctx context.Context) (map[string]int, error)

	// Batch operations
	WarmUpCache(ctx context.Context, userID uuid.UUID, tracks []models.SpotifyTrack, albums []models.SpotifyAlbum, artists []models.SpotifyArtist) error
}

// MusicCacheRepository handles general music data caching (kept for backward compatibility)
type MusicCacheRepository interface {
	// User music data caching
	SetUserMusicData(ctx context.Context, userID uuid.UUID, data interface{}) error
	GetUserMusicData(ctx context.Context, userID uuid.UUID) (interface{}, error)
	AddToRecentlyPlayed(ctx context.Context, userID uuid.UUID, track *models.SpotifyTrack) error

	// Search results caching
	SetSearchResults(ctx context.Context, query string, resultType string, results interface{}) error
	GetSearchResults(ctx context.Context, query string, resultType string) (interface{}, error)

	// Listening history
	SetListeningHistory(ctx context.Context, userID uuid.UUID, history interface{}) error
	GetListeningHistory(ctx context.Context, userID uuid.UUID) (interface{}, error)

	// Cache management
	InvalidateUserCache(ctx context.Context, userID uuid.UUID) error
	InvalidateSearchCache(ctx context.Context, query string) error
	GetCacheStats(ctx context.Context) (map[string]int, error)
}

// Repository container with optimized structure
type Repositories struct {
	User            UserRepository
	Review          ReviewRepository
	Playlist        PlaylistRepository
	UserPreferences UserPreferencesRepository
	Session         SessionRepository
	SpotifyCache    SpotifyCacheRepository // New: Optimized Spotify data cache
	MusicCache      MusicCacheRepository   // Legacy: General music cache (will be phased out)
}
