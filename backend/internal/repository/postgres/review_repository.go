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
		INSERT INTO reviews (id, user_id, spotify_id, spotify_type, rating, review_text, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		review.ID, review.UserID, review.SpotifyID, review.SpotifyType, review.Rating,
		review.ReviewText, review.CreatedAt, review.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	return nil
}

func (r *reviewRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	query := `
		SELECT r.id, r.user_id, r.spotify_id, r.spotify_type, r.rating, r.review_text, r.created_at, r.updated_at,
		       u.id, u.name, u.email, u.bio, u.avatar, u.created_at, u.updated_at
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.id = $1
	`

	review := &models.Review{}
	user := &models.User{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&review.ID, &review.UserID, &review.SpotifyID, &review.SpotifyType, &review.Rating,
		&review.ReviewText, &review.CreatedAt, &review.UpdatedAt,
		&user.ID, &user.Name, &user.Email, &user.Bio, &user.Avatar,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("review not found")
		}
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	review.User = user
	return review, nil
}

func (r *reviewRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Review, error) {
	query := `
		SELECT r.id, r.user_id, r.spotify_id, r.spotify_type, r.rating, r.review_text, r.created_at, r.updated_at,
		       u.id, u.name, u.email, u.bio, u.avatar, u.created_at, u.updated_at
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.user_id = $1
		ORDER BY r.created_at DESC
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
		user := &models.User{}
		err := rows.Scan(
			&review.ID, &review.UserID, &review.SpotifyID, &review.SpotifyType, &review.Rating,
			&review.ReviewText, &review.CreatedAt, &review.UpdatedAt,
			&user.ID, &user.Name, &user.Email, &user.Bio, &user.Avatar,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}
		review.User = user
		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reviews: %w", err)
	}

	return reviews, nil
}

func (r *reviewRepository) GetBySpotifyID(ctx context.Context, spotifyID string, spotifyType string, limit, offset int) ([]*models.Review, error) {
	query := `
		SELECT r.id, r.user_id, r.spotify_id, r.spotify_type, r.rating, r.review_text, r.created_at, r.updated_at,
		       u.id, u.name, u.email, u.bio, u.avatar, u.created_at, u.updated_at
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.spotify_id = $1 AND r.spotify_type = $2
		ORDER BY r.created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Pool.Query(ctx, query, spotifyID, spotifyType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list reviews by spotify item: %w", err)
	}
	defer rows.Close()

	var reviews []*models.Review
	for rows.Next() {
		review := &models.Review{}
		user := &models.User{}
		err := rows.Scan(
			&review.ID, &review.UserID, &review.SpotifyID, &review.SpotifyType, &review.Rating,
			&review.ReviewText, &review.CreatedAt, &review.UpdatedAt,
			&user.ID, &user.Name, &user.Email, &user.Bio, &user.Avatar,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}
		review.User = user
		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reviews: %w", err)
	}

	return reviews, nil
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
		SELECT r.id, r.user_id, r.spotify_id, r.spotify_type, r.rating, r.review_text, r.created_at, r.updated_at,
		       u.id, u.name, u.email, u.bio, u.avatar, u.created_at, u.updated_at
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		ORDER BY r.created_at DESC
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
		user := &models.User{}
		err := rows.Scan(
			&review.ID, &review.UserID, &review.SpotifyID, &review.SpotifyType, &review.Rating,
			&review.ReviewText, &review.CreatedAt, &review.UpdatedAt,
			&user.ID, &user.Name, &user.Email, &user.Bio, &user.Avatar,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}
		review.User = user
		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reviews: %w", err)
	}

	return reviews, nil
}

func (r *reviewRepository) GetAverageRating(ctx context.Context, spotifyID string, spotifyType string) (float64, error) {
	query := `
		SELECT AVG(rating::float) 
		FROM reviews 
		WHERE spotify_id = $1 AND spotify_type = $2
	`

	var avgRating *float64
	err := r.db.Pool.QueryRow(ctx, query, spotifyID, spotifyType).Scan(&avgRating)
	if err != nil {
		return 0, fmt.Errorf("failed to get average rating: %w", err)
	}

	if avgRating == nil {
		return 0, nil // No reviews yet
	}

	return *avgRating, nil
}

func (r *reviewRepository) GetRatingCount(ctx context.Context, spotifyID string, spotifyType string) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM reviews 
		WHERE spotify_id = $1 AND spotify_type = $2
	`

	var count int
	err := r.db.Pool.QueryRow(ctx, query, spotifyID, spotifyType).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get rating count: %w", err)
	}

	return count, nil
}
