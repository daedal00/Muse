package postgres

import (
	"context"

	"github.com/daedal00/muse/backend/internal/database"
	"github.com/daedal00/muse/backend/internal/models"
	"github.com/google/uuid"
)

type commentRepository struct {
	db *database.PostgresDB
}

func NewCommentRepository(db *database.PostgresDB) *commentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *models.Comment) error {
	query := `
		INSERT INTO comments (id, user_id, album_id, track_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Pool.Exec(ctx, query,
		comment.ID, comment.UserID, comment.AlbumID, comment.TrackID,
		comment.Content, comment.CreatedAt, comment.UpdatedAt,
	)
	return err
}

func (r *commentRepository) GetByAlbumID(ctx context.Context, albumID uuid.UUID, limit, offset int) ([]*models.Comment, error) {
	query := `
		SELECT id, user_id, album_id, track_id, content, created_at, updated_at
		FROM comments
		WHERE album_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Pool.Query(ctx, query, albumID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		c := &models.Comment{}
		err := rows.Scan(
			&c.ID, &c.UserID, &c.AlbumID, &c.TrackID, &c.Content, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (r *commentRepository) CountByAlbumID(ctx context.Context, albumID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE album_id = $1`
	var count int
	err := r.db.Pool.QueryRow(ctx, query, albumID).Scan(&count)
	return count, err
}

func (r *commentRepository) GetByTrackID(ctx context.Context, trackID uuid.UUID, limit, offset int) ([]*models.Comment, error) {
	query := `
		SELECT id, user_id, album_id, track_id, content, created_at, updated_at
		FROM comments
		WHERE track_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Pool.Query(ctx, query, trackID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		c := &models.Comment{}
		err := rows.Scan(
			&c.ID, &c.UserID, &c.AlbumID, &c.TrackID, &c.Content, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (r *commentRepository) CountByTrackID(ctx context.Context, trackID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE track_id = $1`
	var count int
	err := r.db.Pool.QueryRow(ctx, query, trackID).Scan(&count)
	return count, err
}

func (r *commentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}
