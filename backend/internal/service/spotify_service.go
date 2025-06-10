package service

import (
	"context"

	"github.com/daedal00/muse/backend/graph/model"
)

type SpotifyService interface {
	SearchAlbums(ctx context.Context, query string, limit, offset int) ([]*model.AlbumSearchResult, error)
	GetAlbumById(ctx context.Context, id string) ([]*model.Album ,error)
}

