package postgres

import (
	"context"
	"fmt"

	"github.com/daedal00/muse/backend/internal/database"
	"github.com/daedal00/muse/backend/internal/models"
	"github.com/daedal00/muse/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type reviewRepository struct {
	db *database.PostgresDB
}

func NewReviewRepository(db *database.PostgresDB) repository.ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(ctx context.Context, review *models.Review) error {
	query := `
		INSERT INTO reviews (id, user_id, album_id, rating, review_text, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	_, err := r.db.Pool.Exec(ctx, query,
		review.ID, review.UserID, review.AlbumID, review.Rating,
		review.ReviewText, review.CreatedAt, review.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}
	
	return nil
}

func (r *reviewRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	query := `
		SELECT id, user_id, album_id, rating, review_text, created_at, updated_at
		FROM reviews 
		WHERE id = $1
	`
	
	review := &models.Review{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&review.ID, &review.UserID, &review.AlbumID, &review.Rating,
		&review.ReviewText, &review.CreatedAt, &review.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("review not found")
		}
		return nil, fmt.Errorf("failed to get review: %w", err)
	}
	
	return review, nil
}

func (r *reviewRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Review, error) {
	query := `
		SELECT id, user_id, album_id, rating, review_text, created_at, updated_at
		FROM reviews 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list reviews by user: %w", err)
	}
	defer rows.Close()
	
	var reviews []*models.Review
	for rows.Next() {
		review := &models.Review{}
		err := rows.Scan(
			&review.ID, &review.UserID, &review.AlbumID, &review.Rating,
			&review.ReviewText, &review.CreatedAt, &review.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}
		reviews = append(reviews, review)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reviews: %w", err)
	}
	
	return reviews, nil
}

func (r *reviewRepository) GetByAlbumID(ctx context.Context, albumID uuid.UUID, limit, offset int) ([]*models.Review, error) {
	query := `
		SELECT id, user_id, album_id, rating, review_text, created_at, updated_at
		FROM reviews 
		WHERE album_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Pool.Query(ctx, query, albumID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list reviews by album: %w", err)
	}
	defer rows.Close()
	
	var reviews []*models.Review
	for rows.Next() {
		review := &models.Review{}
		err := rows.Scan(
			&review.ID, &review.UserID, &review.AlbumID, &review.Rating,
			&review.ReviewText, &review.CreatedAt, &review.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}
		reviews = append(reviews, review)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reviews: %w", err)
	}
	
	return reviews, nil
}

func (r *reviewRepository) GetByUserAndAlbum(ctx context.Context, userID, albumID uuid.UUID) (*models.Review, error) {
	query := `
		SELECT id, user_id, album_id, rating, review_text, created_at, updated_at
		FROM reviews 
		WHERE user_id = $1 AND album_id = $2
	`
	
	review := &models.Review{}
	err := r.db.Pool.QueryRow(ctx, query, userID, albumID).Scan(
		&review.ID, &review.UserID, &review.AlbumID, &review.Rating,
		&review.ReviewText, &review.CreatedAt, &review.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("review not found")
		}
		return nil, fmt.Errorf("failed to get review: %w", err)
	}
	
	return review, nil
}

func (r *reviewRepository) Update(ctx context.Context, review *models.Review) error {
	query := `
		UPDATE reviews 
		SET rating = $2, review_text = $3, updated_at = NOW()
		WHERE id = $1
	`
	
	result, err := r.db.Pool.Exec(ctx, query,
		review.ID, review.Rating, review.ReviewText,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update review: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("review not found")
	}
	
	return nil
}

func (r *reviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM reviews WHERE id = $1`
	
	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("review not found")
	}
	
	return nil
}

func (r *reviewRepository) List(ctx context.Context, limit, offset int) ([]*models.Review, error) {
	query := `
		SELECT id, user_id, album_id, rating, review_text, created_at, updated_at
		FROM reviews 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list reviews: %w", err)
	}
	defer rows.Close()
	
	var reviews []*models.Review
	for rows.Next() {
		review := &models.Review{}
		err := rows.Scan(
			&review.ID, &review.UserID, &review.AlbumID, &review.Rating,
			&review.ReviewText, &review.CreatedAt, &review.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}
		reviews = append(reviews, review)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reviews: %w", err)
	}
	
	return reviews, nil
} 