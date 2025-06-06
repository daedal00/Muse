package graph

import "github.com/daedal00/muse/backend/graph/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	users []*model.User
	passwordHashes map[string]string

	albums []*model.Album
	reviews []*model.Review
	playlists []*model.Playlist
}

func NewResolver() *Resolver {
	return &Resolver{
		users: make([]*model.User, 0),
		passwordHashes: make(map[string]string),
	}
}