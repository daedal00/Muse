package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID  `db:"id"`
	UserID    uuid.UUID  `db:"user_id"`
	AlbumID   *uuid.UUID `db:"album_id"`
	TrackID   *uuid.UUID `db:"track_id"`
	Content   string     `db:"content"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`

	// Associations - ignored by DB mapper
	User  *User  `db:"-"`
	Album *Album `db:"-"`
	Track *Track `db:"-"`
}
