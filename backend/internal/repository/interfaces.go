package repository

import (
	"context"

	"github.com/daedal00/muse/backend/internal/models"
	"github.com/google/uuid"
)

// Repository interfaces for all entities

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
	GetAll(ctx context.Context, limit int, after *string) ([]*models.User, error)
	Count(ctx context.Context) (int, error)
}

type ArtistRepository interface {
	Create(ctx context.Context, artist *models.Artist) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Artist, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.Artist, error)
	GetBySpotifyID(ctx context.Context, spotifyID string) (*models.Artist, error)
	Update(ctx context.Context, artist *models.Artist) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.Artist, error)
}

type AlbumRepository interface {
	Create(ctx context.Context, album *models.Album) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Album, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.Album, error)
	GetBySpotifyID(ctx context.Context, spotifyID string) (*models.Album, error)
	GetByArtistID(ctx context.Context, artistID uuid.UUID, limit, offset int) ([]*models.Album, error)
	CountByArtistID(ctx context.Context, artistID uuid.UUID) (int, error)
	Update(ctx context.Context, album *models.Album) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.Album, error)
	Count(ctx context.Context) (int, error)
}

type TrackRepository interface {
	Create(ctx context.Context, track *models.Track) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Track, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.Track, error)
	GetBySpotifyID(ctx context.Context, spotifyID string) (*models.Track, error)
	GetByAlbumID(ctx context.Context, albumID uuid.UUID, limit, offset int) ([]*models.Track, error)
	CountByAlbumID(ctx context.Context, albumID uuid.UUID) (int, error)
	Update(ctx context.Context, track *models.Track) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.Track, error)
	Count(ctx context.Context) (int, error)
}

type ReviewRepository interface {
	Create(ctx context.Context, review *models.Review) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Review, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Review, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int, error)
	GetByAlbumID(ctx context.Context, albumID uuid.UUID, limit, offset int) ([]*models.Review, error)
	CountByAlbumID(ctx context.Context, albumID uuid.UUID) (int, error)
	AverageRatingByAlbumID(ctx context.Context, albumID uuid.UUID) (*float64, error)
	GetByUserAndAlbum(ctx context.Context, userID, albumID uuid.UUID) (*models.Review, error)
	Update(ctx context.Context, review *models.Review) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.Review, error)
	Count(ctx context.Context) (int, error)
}

type TrackReviewRepository interface {
	Create(ctx context.Context, review *models.TrackReview) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.TrackReview, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.TrackReview, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int, error)
	GetByTrackID(ctx context.Context, trackID uuid.UUID, limit, offset int) ([]*models.TrackReview, error)
	CountByTrackID(ctx context.Context, trackID uuid.UUID) (int, error)
	AverageRatingByTrackID(ctx context.Context, trackID uuid.UUID) (*float64, error)
	GetByUserAndTrack(ctx context.Context, userID, trackID uuid.UUID) (*models.TrackReview, error)
	ListTopTracksByAverageRating(ctx context.Context, limit int) ([]*models.TrackReviewSummary, error)
	ListTopTracksByReviewCount(ctx context.Context, limit int) ([]*models.TrackReviewSummary, error)
	ListFavoritesByUserID(ctx context.Context, userID uuid.UUID, minRating, limit int) ([]*models.TrackReview, error)
	List(ctx context.Context, limit, offset int) ([]*models.TrackReview, error)
	Count(ctx context.Context) (int, error)
	Update(ctx context.Context, review *models.TrackReview) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PlaylistRepository interface {
	Create(ctx context.Context, playlist *models.Playlist) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Playlist, error)
	GetBySpotifyID(ctx context.Context, spotifyID string) (*models.Playlist, error)
	GetByCreatorID(ctx context.Context, creatorID uuid.UUID, limit, offset int) ([]*models.Playlist, error)
	CountByCreatorID(ctx context.Context, creatorID uuid.UUID) (int, error)
	Update(ctx context.Context, playlist *models.Playlist) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.Playlist, error)
	Count(ctx context.Context) (int, error)

	// Playlist track operations
	AddTrack(ctx context.Context, playlistID, trackID uuid.UUID, position int) error
	RemoveTrack(ctx context.Context, playlistID, trackID uuid.UUID) error
	GetTracks(ctx context.Context, playlistID uuid.UUID, limit, offset int) ([]*models.Track, error)
	CountTracks(ctx context.Context, playlistID uuid.UUID) (int, error)
	ReorderTracks(ctx context.Context, playlistID uuid.UUID, trackPositions map[uuid.UUID]int) error
}

type SpotifyTokenRepository interface {
	Upsert(ctx context.Context, token *models.SpotifyToken) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.SpotifyToken, error)
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	GetByID(ctx context.Context, id string) (*models.Session, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Session, error)
	Delete(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

// MusicCacheRepository handles caching of user music data and search results
type MusicCacheRepository interface {
	// User music data caching
	SetUserMusicData(ctx context.Context, userID uuid.UUID, data interface{}) error
	GetUserMusicData(ctx context.Context, userID uuid.UUID) (interface{}, error)
	AddToRecentlyPlayed(ctx context.Context, userID uuid.UUID, track *models.Track) error

	// Search results caching
	SetSearchResults(ctx context.Context, query string, resultType string, results interface{}) error
	GetSearchResults(ctx context.Context, query string, resultType string) (interface{}, error)

	// Listening history
	SetListeningHistory(ctx context.Context, userID uuid.UUID, history interface{}) error
	GetListeningHistory(ctx context.Context, userID uuid.UUID) (interface{}, error)

	// Popular content caching
	SetPopularAlbums(ctx context.Context, albums []*models.Album) error
	GetPopularAlbums(ctx context.Context) ([]*models.Album, error)
	SetPopularTracks(ctx context.Context, tracks []*models.Track) error
	GetPopularTracks(ctx context.Context) ([]*models.Track, error)

	// Cache management
	InvalidateUserCache(ctx context.Context, userID uuid.UUID) error
	InvalidateSearchCache(ctx context.Context, query string) error
	GetCacheStats(ctx context.Context) (map[string]int, error)
}

// Repository container
type Repositories struct {
	User        UserRepository
	Artist      ArtistRepository
	Album       AlbumRepository
	Track       TrackRepository
	Review      ReviewRepository
	TrackReview TrackReviewRepository
	Playlist    PlaylistRepository
	Spotify     SpotifyTokenRepository
	Session     SessionRepository
	Comment     CommentRepository
	MusicCache  MusicCacheRepository // New: Redis music cache
}

type CommentRepository interface {
	Create(ctx context.Context, comment *models.Comment) error
	GetByAlbumID(ctx context.Context, albumID uuid.UUID, limit, offset int) ([]*models.Comment, error)
	CountByAlbumID(ctx context.Context, albumID uuid.UUID) (int, error)
	GetByTrackID(ctx context.Context, trackID uuid.UUID, limit, offset int) ([]*models.Comment, error)
	CountByTrackID(ctx context.Context, trackID uuid.UUID) (int, error)
}
